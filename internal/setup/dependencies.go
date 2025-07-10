package setup

import (
	"context"
	"echo-household-budget/config"
	domainmodel "echo-household-budget/internal/domain/model"
	domainRepository "echo-household-budget/internal/domain/repository"
	domainService "echo-household-budget/internal/domain/service"
	"echo-household-budget/internal/domain/service/functioncalling"
	"echo-household-budget/internal/handler"
	"echo-household-budget/internal/infrastructure/llm"
	"echo-household-budget/internal/infrastructure/persistence/repository"
	"echo-household-budget/internal/infrastructure/storage/s3"
	"echo-household-budget/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Dependencies はアプリケーションの依存関係を管理する構造体
type Dependencies struct {
	// Repositories
	KaimemoRepository         repository.KaimemoRepository
	LineRepository            repository.LineRepository
	UserAccountRepository     domainmodel.UserAccountRepository
	CategoryRepository        domainmodel.CategoryRepository
	HouseHoldRepository       domainmodel.HouseHoldRepository
	ShoppingRepository        domainmodel.ShoppingRepository
	ReceiptAnalyzeRepository  domainmodel.ReceiptAnalyzeRepository
	InformationRepository     domainRepository.InformationRepository
	UserInformationRepository domainRepository.UserInformationRepository
	ChatMessageRepository     domainRepository.ChatMessageRepository
	FileStorageRepository     domainRepository.FileStorageRepository

	// Services
	UserAccountService domainService.UserAccountService
	HouseHoldService   domainService.HouseHoldService

	// Function Calling
	ToolRegistry *functioncalling.ToolRegistry
	LLMClient    *llm.LLMClient

	// Use Cases
	SessionManager              usecase.SessionManager
	KaimemoService              usecase.KaimemoService
	ShoppingUsecase             usecase.ShoppingUsecase
	LineAuthService             usecase.LineAuthService
	ReceiptAnalyzeUsecase       usecase.ReceiptAnalyzeUsecase
	CreateInformationUsecase    usecase.CreateInformationUsecase
	FetchInformationUsecase     usecase.FetchInformationUsecase
	PublishInformationUsecase   usecase.PublishInformationUsecase
	FetchUserInformationUsecase usecase.FetchUserInformationUsecase
	RegisterChatMessageUsecase  usecase.RegisterChatMessageUsecase
	FetchChatMessageUsecase     usecase.FetchChatMessageUsecase

	// Handlers
	KaimemoHandler                   handler.KaimemoHandler
	LineAuthHandler                  handler.AuthHandler
	HouseHoldHandler                 handler.HouseHoldHandler
	ReceiptAnalyzeHandler            handler.ReceiptAnalyzeHandler
	CreateInformationHandler         handler.CreateInformationHandler
	FetchInformationHandler          handler.FetchInformationsHandler
	PublishInformationHandler        handler.PublishInformationHandler
	FetchUserInformationHandler      handler.FetchUserInformationHandler
	UpdateReadUserInformationHandler handler.UpdateReadUserInformationHandler
	FetchChatMessagesHandler         handler.FetchChatMessagesHandler
	ChatMessageTelegraphHandler      handler.ChatMessageTelegraphHandler
	DeleteInformationHandler         handler.DeleteInformationHandler
	FetchInformationDetailHandler    handler.FetchInformationDetailHandler
	PutInformationHandler            handler.PutInformationHandler
}

// NewDependencies は依存関係を初期化して返す
func NewDependencies(appConfig *config.AppConfig) *Dependencies {
	// データベース接続の設定
	db, err := config.NewDBConnection(appConfig.DatabaseConfig)
	if err != nil {
		panic(err)
	}

	// AWS S3設定
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(appConfig.S3Config.Region),
		awsconfig.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			appConfig.S3Config.AccessKeyID,
			appConfig.S3Config.SecretAccessKey,
			"",
		))),
	)
	if err != nil {
		panic(err)
	}
	s3Client := awss3.NewFromConfig(cfg)

	// リポジトリの初期化
	deps := &Dependencies{}
	deps.KaimemoRepository = repository.NewNotionRepository(
		appConfig.NotionAPIKey,
		appConfig.NotionKaimemoDatabaseInputID,
		appConfig.NotionKaimemoDatabaseSummaryRecordID,
	)
	deps.LineRepository = repository.NewLineRepository(appConfig.LINEConfig)
	deps.UserAccountRepository = repository.NewUserAccountRepository(db)
	deps.CategoryRepository = repository.NewCategoryRepository(db)
	deps.HouseHoldRepository = repository.NewHouseHoldRepository(db)
	deps.ShoppingRepository = repository.NewShoppingRepository(db)
	deps.ReceiptAnalyzeRepository = repository.NewReceiptRepository(db)
	deps.InformationRepository = repository.NewInformationRepository(db)
	deps.UserInformationRepository = repository.NewUserInformationRepository(db)
	deps.ChatMessageRepository = repository.NewChatMessageRepository(db)
	deps.FileStorageRepository = s3.NewS3FileStorage(s3Client, appConfig.S3Config.BucketName)

	// サービスの初期化
	deps.UserAccountService = domainService.NewUserAccountService(deps.UserAccountRepository, deps.CategoryRepository, deps.HouseHoldRepository)
	deps.HouseHoldService = domainService.NewHouseHoldService(deps.HouseHoldRepository, deps.ShoppingRepository, deps.CategoryRepository)

	// Function Calling の初期化
	deps.ToolRegistry = functioncalling.NewToolRegistry(deps.ShoppingRepository, deps.HouseHoldRepository)
	deps.LLMClient = llm.NewLLMClient(deps.ToolRegistry.GetAllTools())

	// ユースケースの初期化
	deps.SessionManager = usecase.NewSessionManager()
	deps.KaimemoService = usecase.NewKaimemoService(deps.KaimemoRepository)
	deps.ShoppingUsecase = usecase.NewShoppingUsecase(deps.ShoppingRepository)
	deps.LineAuthService = usecase.NewLineAuthService(deps.LineRepository, deps.UserAccountRepository, deps.UserAccountService, deps.SessionManager)
	deps.ReceiptAnalyzeUsecase = usecase.NewReceiptAnalyzeUsecase(deps.ReceiptAnalyzeRepository, deps.FileStorageRepository, deps.HouseHoldService)
	deps.CreateInformationUsecase = usecase.NewCreateInformationUsecase(deps.InformationRepository)
	deps.FetchInformationUsecase = usecase.NewFetchInformationUsecase(deps.InformationRepository)
	deps.PublishInformationUsecase = usecase.NewPublishInformationUsecase(deps.InformationRepository, deps.UserInformationRepository, deps.UserAccountService)
	deps.FetchUserInformationUsecase = usecase.NewFetchUserInformationUsecase(deps.UserInformationRepository)
	deps.RegisterChatMessageUsecase = usecase.NewRegisterChatMessageUsecase(deps.ChatMessageRepository, deps.LLMClient)
	deps.FetchChatMessageUsecase = usecase.NewFetchChatMessageUsecase(deps.ChatMessageRepository)

	// ハンドラーの初期化
	deps.KaimemoHandler = handler.NewKaimemoHandler(deps.KaimemoService, deps.ShoppingUsecase)
	deps.LineAuthHandler = handler.NewLineAuthHandler(deps.LineAuthService, appConfig)
	deps.HouseHoldHandler = handler.NewHouseHoldHandler(deps.HouseHoldService, deps.UserAccountService)
	deps.ReceiptAnalyzeHandler = handler.NewReceiptAnalyzeHandler(deps.ReceiptAnalyzeUsecase)
	deps.CreateInformationHandler = handler.NewCreateInformationHandler(deps.CreateInformationUsecase)
	deps.FetchInformationHandler = handler.NewFetchInformationsHandler(deps.FetchInformationUsecase)
	deps.PublishInformationHandler = handler.NewPublishInformationHandler(deps.PublishInformationUsecase)
	deps.FetchUserInformationHandler = handler.NewFetchUserInformationHandler(deps.FetchUserInformationUsecase)
	deps.UpdateReadUserInformationHandler = handler.NewUpdateReadUserInformationHandler(deps.UserInformationRepository)
	deps.FetchChatMessagesHandler = handler.NewFetchChatMessagesHandler()
	deps.ChatMessageTelegraphHandler = handler.NewChatMessageTelegraphHandler(deps.RegisterChatMessageUsecase, deps.FetchChatMessageUsecase)
	deps.DeleteInformationHandler = handler.NewDeleteInformationHandler()
	deps.FetchInformationDetailHandler = handler.NewFetchInformationDetailHandler()
	deps.PutInformationHandler = handler.NewPutInformationHandler()

	return deps
}
