package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ExerciseCompleteMessage struct {
	UserID              string `json:"user_id"`
	ExerciseID          uint64 `json:"exercise_id"`
	ExerciseLang        string `json:"exercise_lang"`
	SpentTime           uint64 `json:"spent_time"`
	WordsCount          uint16 `json:"words_count"`
	WordsCorrectedCount uint16 `json:"words_corrected_count"`
}

type ExerciseCompleteHandler struct {
	ECUseCase domain.ExerciseCompleteUseCase
	UUseCase  domain.UserUseCase
	validate  *validator.Validate
}

func NewExerciseCompleteHandler(
	v *validator.Validate,
	ec domain.ExerciseCompleteUseCase,
	u domain.UserUseCase,
) *ExerciseCompleteHandler {
	return &ExerciseCompleteHandler{
		ECUseCase: ec,
		UUseCase:  u,
		validate:  v,
	}
}

func (ech *ExerciseCompleteHandler) Handle(msg *message.Message) error {
	ctx := context.Background()

	// Extract tracing context
	propagator := otel.GetTextMapPropagator()
	ctx = propagator.Extract(ctx, propagation.MapCarrier(msg.Metadata))

	ctx, span := otel.Tracer("kafka-consumer").Start(ctx, "Handle Kafka Message")
	defer span.End()

	//ctx, span := otel.Tracer("kafka-consumer").Start(context.Background(), "Handle Kafka Message")
	//defer span.End()

	var ecMsg ExerciseCompleteMessage
	if err := json.Unmarshal(msg.Payload, &ecMsg); err != nil {
		return err
	}

	userID, err := uuid.Parse(ecMsg.UserID)
	if err != nil {
		return fmt.Errorf("UserID not parsed: %w", err)
	}

	user, err := ech.UUseCase.GetByID(ctx, auth.UserID(userID))
	if err != nil {
		return fmt.Errorf("failed to get user from Keycloak: %w", err)
	}

	ec := &domain.ExerciseComplete{
		UserID:              auth.UserID(userID),
		UserName:            user.Username,
		ExerciseID:          domain.ExerciseID(ecMsg.ExerciseID),
		ExerciseLang:        domain.ExerciseLang(ecMsg.ExerciseLang),
		SpentTime:           ecMsg.SpentTime,
		WordsCount:          ecMsg.WordsCount,
		WordsCorrectedCount: ecMsg.WordsCorrectedCount,
		EventTime:           time.Now(),
	}

	// Сохранение данных через ExerciseCompleteUseCase
	if err := ech.ECUseCase.Store(ctx, ec); err != nil {
		return err
	}

	fmt.Println("Message processed successfully:", ec.ExerciseID)

	return nil
}
