using System;
namespace RemoteSigner.Models {
    public static class ErrorCodes {
        public readonly static String InternalServerError = "INTERNAL_SERVER_ERROR";
        public readonly static String NotFound = "NOT_FOUND";
        public readonly static String EmailAlreadyInUse = "EMAIL_ALREADY_IN_USE";
        public readonly static String NoDataAvailable = "NO_DATA_AVAILABLE";
        public readonly static String InvalidLoginInformation = "INVALID_LOGIN_INFORMATION";
        public readonly static String NotLogged = "NOT_LOGGED";
        public readonly static String AlreadyExists = "ALREADY_EXISTS";
        public readonly static String PermissionDenied = "PERMISSION_DENIED";
        public readonly static String InvalidTokenType = "INVALID_TOKEN_TYPE";
        public readonly static String InvalidFieldData = "INVALID_FIELD_DATA";
        public readonly static String AlreadyClient = "ALREADY_CLIENT";
        public readonly static String AlreadyPaid = "ALREADY_PAID";
        public readonly static String PaymentError = "PAYMENT_ERROR";
        public readonly static String InsufficientFunds = "INSUFFICIENT_FUNDS";
        public readonly static String BankingSystemOffline = "BANKING_SYSTEM_OFFLINE";
        public readonly static String OutdatedAPI = "OUTDATED_API";
        public readonly static String BankNotSupported = "BANK_NOT_SUPPORTED";
        public readonly static String VaultSystemOffline = "VAULT_SYSTEM_OFFLINE";
        public readonly static String ServerIsBusy = "SERVER_IS_BUSY";
        public readonly static String Revoked = "REVOKED";
        public readonly static String AlreadySigned = "ALREADY_SIGNED";
        public readonly static String Rejected = "REJECTED";
        public readonly static String OperationNotSupported = "OPERATION_NOT_SUPPORTED";
        public readonly static String GraphQLError = "GRAPHQL_ERROR";
        public readonly static String OperationLimitExceeded = "OPERATION_LIMIT_EXCEEDED";
        public readonly static String InvalidTransactionDate = "INVALID_TRANSACTION_DATE";
        public readonly static String InvalidSignature = "INVALID_SIGNATURE";

    }
}
