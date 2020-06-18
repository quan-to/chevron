package QuantoError

// region Public Use
const InternalServerError = "INTERNAL_SERVER_ERROR"
const NotFound = "NOT_FOUND"
const EmailAlreadyInUse = "EMAIL_ALREADY_IN_USE"
const NoDataAvailable = "NO_DATA_AVAILABLE"
const InvalidLoginInformation = "INVALID_LOGIN_INFORMATION"
const NotLogged = "NOT_LOGGED"
const AlreadyExists = "ALREADY_EXISTS"
const PermissionDenied = "PERMISSION_DENIED"
const InvalidTokenType = "INVALID_TOKEN_TYPE"
const InvalidFieldData = "INVALID_FIELD_DATA"
const AlreadyClient = "ALREADY_CLIENT"
const AlreadyPaid = "ALREADY_PAID"
const PaymentError = "PAYMENT_ERROR"
const InsufficientFunds = "INSUFFICIENT_FUNDS"
const BankingSystemOffline = "BANKING_SYSTEM_OFFLINE"
const OutdatedAPI = "OUTDATED_API"
const BankNotSupported = "BANK_NOT_SUPPORTED"
const VaultSystemOffline = "VAULT_SYSTEM_OFFLINE"
const ServerIsBusy = "SERVER_IS_BUSY"
const Revoked = "REVOKED"
const AlreadySigned = "ALREADY_SIGNED"
const Rejected = "REJECTED"
const OperationNotSupported = "OPERATION_NOT_SUPPORTED"
const GraphQLError = "GRAPHQL_ERROR"
const OperationLimitExceeded = "OPERATION_LIMIT_EXCEEDED"
const InvalidTransactionDate = "INVALID_TRANSACTION_DATE"
const BoletoCreationNotEnabled = "BOLETO_CREATION_NOT_ENABLED"
const BoletoOurNumberExausted = "BOLETO_OUR_NUMBER_EXAUSTED"
const NotImplemented = "NOT_IMPLEMENTED"

// endregion

// region Internal Use - Don"t worry about these if you"re a partner.
const EverythingIsTerrible = "EVERYTHING_IS_TERRIBLE"
const QuantoInternalError = "QUANTO_INTERNAL_ERROR"
const RoutingSystemOffline = "ROUTING_SYSTEM_OFFLINE"
const QITSystemOffline = "QIT_SYSTEM_OFFLINE"
const TargetConnectionError = "CONNECTION_ERROR"
const VaulterIsDead = "VAULTER_IS_DEAD"
const SynchronizationError = "SYNCHRONIZATION_ERROR"
const RoutingError = "ROUTING_ERROR"

// endregion
