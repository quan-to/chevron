using System;
namespace RemoteSigner.Database.Attributes {
    [AttributeUsage(AttributeTargets.Class)]
    public class DBTable : Attribute {
        public readonly string TableName;
        public DBTable(string tableName) {
            TableName = tableName;
        }
    }
}
