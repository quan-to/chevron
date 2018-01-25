namespace RemoteSigner.Models.Attributes {
    interface IHttpMethod {
        string Path { get; }
        string Method { get; }
    }
}
