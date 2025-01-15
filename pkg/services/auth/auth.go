package auth

import (
	"auth/pkg/core/models"
	"auth/pkg/jwt"
	"auth/pkg/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidAppId = errors.New("invalid app id")
var ErrUserAlreadyExists = errors.New("user already exists")

type Auth struct {
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration

	logger *slog.Logger
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, PWHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

func NewAuth(userSaver UserSaver, userProvider UserProvider,
	appProvider AppProvider, tokenTTL time.Duration, logger *slog.Logger) *Auth {
	return &Auth{userSaver: userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
		logger:       logger}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	const op = "auth.Login"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("logging in user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PWHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged succesfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, err

}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	PWHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)

	}

	id, err := a.userSaver.SaveUser(ctx, email, PWHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("user alrady exists", err)

			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)

		}
		log.Error("failed to save user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("registering user")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found")

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("user not found", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is admin", isAdmin))

	return isAdmin, err
}
