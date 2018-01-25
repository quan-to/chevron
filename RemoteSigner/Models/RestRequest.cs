using System;
using System.Collections.Specialized;
using System.Net;
using System.Text;

namespace RemoteSigner.Models {

    [Serializable]
    public class RestRequest {
        readonly string[] acceptTypes;
        readonly Encoding contentEncoding;
        readonly long contentLength;
        readonly string contentType;
        readonly CookieCollection cookies;
        readonly string httpMethod;
        readonly bool isAuthenticated;
        readonly bool isLocal;
        readonly bool isSecureConnection;
        readonly bool isWebSocketRequest;
        readonly bool keepAlive;
        readonly NameValueCollection queryString;
        readonly string rawUrl;
        readonly string userAgent;
        readonly string userHostAddress;
        readonly string userHostName;
        readonly string[] userLanguages;
        readonly string bodyData;
        readonly Uri url;
        readonly Uri urlReferrer;

        //
        // Summary:
        //     Gets the MIME types accepted by the client.
        //
        // Returns:
        //     A System.String array that contains the type names specified in the request's
        //     Accept header or null if the client request did not include an Accept header.
        public string[] AcceptTypes { get { return acceptTypes; } }
        //
        // Summary:
        //     Gets the content encoding that can be used with data sent with the request
        //
        // Returns:
        //     An System.Text.Encoding object suitable for use with the data in the System.Net.HttpListenerRequest.InputStream
        //     property.
        public Encoding ContentEncoding { get { return contentEncoding; } }
        //
        // Summary:
        //     Gets the length of the body data included in the request.
        //
        // Returns:
        //     The value from the request's Content-Length header. This value is -1 if the content
        //     length is not known.
        public long ContentLength { get { return contentLength; } }
        //
        // Summary:
        //     Gets the MIME type of the body data included in the request.
        //
        // Returns:
        //     A System.String that contains the text of the request's Content-Type header.
        public string ContentType { get { return contentType; } }
        //
        // Summary:
        //     Gets the cookies sent with the request.
        //
        // Returns:
        //     A System.Net.CookieCollection that contains cookies that accompany the request.
        //     This property returns an empty collection if the request does not contain cookies.
        public CookieCollection Cookies { get { return cookies; } }
        //
        // Summary:
        //     Gets the HTTP method specified by the client.
        //
        // Returns:
        //     A System.String that contains the method used in the request.
        public string HttpMethod { get { return httpMethod; } }
        //
        // Summary:
        //     Gets a System.Boolean value that indicates whether the client sending this request
        //     is authenticated.
        //
        // Returns:
        //     true if the client was authenticated; otherwise, false.
        public bool IsAuthenticated { get { return isAuthenticated; } }
        //
        // Summary:
        //     Gets a System.Boolean value that indicates whether the request is sent from the
        //     local computer.
        //
        // Returns:
        //     true if the request originated on the same computer as the System.Net.HttpListener
        //     object that provided the request; otherwise, false.
        public bool IsLocal { get { return isLocal; } }
        //
        // Summary:
        //     Gets a System.Boolean value that indicates whether the TCP connection used to
        //     send the request is using the Secure Sockets Layer (SSL) protocol.
        //
        // Returns:
        //     true if the TCP connection is using SSL; otherwise, false.
        public bool IsSecureConnection { get { return isSecureConnection; } }
        //
        // Summary:
        //     Gets a System.Boolean value that indicates whether the TCP connection was a WebSocket
        //     request.
        //
        // Returns:
        //     Returns System.Boolean.true if the TCP connection is a WebSocket request; otherwise,
        //     false.
        public bool IsWebSocketRequest { get { return isWebSocketRequest; } }
        //
        // Summary:
        //     Gets a System.Boolean value that indicates whether the client requests a persistent
        //     connection.
        //
        // Returns:
        //     true if the connection should be kept open; otherwise, false.
        public bool KeepAlive { get { return keepAlive; } }
        //
        // Summary:
        //     Gets the query string included in the request.
        //
        // Returns:
        //     A System.Collections.Specialized.NameValueCollection object that contains the
        //     query data included in the request System.Net.HttpListenerRequest.Url.
        public NameValueCollection QueryString { get { return queryString; } }
        //
        // Summary:
        //     Gets the URL information (without the host and port) requested by the client.
        //
        // Returns:
        //     A System.String that contains the raw URL for this request.
        public string RawUrl { get { return rawUrl; } }
        //
        // Summary:
        //     Gets the user agent presented by the client.
        //
        // Returns:
        //     A System.String object that contains the text of the request's User-Agent header.
        public string UserAgent { get { return userAgent; } }
        //
        // Summary:
        //     Gets the server IP address and port number to which the request is directed.
        //
        // Returns:
        //     A System.String that contains the host address information.
        public string UserHostAddress { get { return userHostAddress; } }
        //
        // Summary:
        //     Gets the DNS name and, if provided, the port number specified by the client.
        //
        // Returns:
        //     A System.String value that contains the text of the request's Host header.
        public string UserHostName { get { return userHostName; } }
        //
        // Summary:
        //     Gets the natural languages that are preferred for the response.
        //
        // Returns:
        //     A System.String array that contains the languages specified in the request's
        //     System.Net.HttpRequestHeader.AcceptLanguage header or null if the client request
        //     did not include an System.Net.HttpRequestHeader.AcceptLanguage header.
        public string[] UserLanguages { get { return userLanguages; } }
        //
        // Summary:
        //     Gets the System.Uri object requested by the client.
        //
        // Returns:
        //     A System.Uri object that identifies the resource requested by the client.
        public Uri Url { get { return url; } }
        //
        // Summary:
        //     Gets the Uniform Resource Identifier (URI) of the resource that referred the
        //     client to the server.
        //
        // Returns:
        //     A System.Uri object that contains the text of the request's System.Net.HttpRequestHeader.Referer
        //     header, or null if the header was not included in the request.
        public Uri UrlReferrer { get { return urlReferrer; } }

        public string BodyData { get { return bodyData; } }

        public RestRequest(HttpListenerRequest request) {
            userAgent = request.UserAgent;
            userHostAddress = request.UserHostAddress;
            userHostName = request.UserHostName;
            userLanguages = request.UserLanguages;
            url = request.Url;
            urlReferrer = request.UrlReferrer;
            acceptTypes = request.AcceptTypes;
            contentEncoding = request.ContentEncoding;
            contentLength = request.ContentLength64;
            contentType = request.ContentType;
            cookies = request.Cookies;
            httpMethod = request.HttpMethod;
            isAuthenticated = request.IsAuthenticated;
            isLocal = request.IsLocal;
            isSecureConnection = request.IsSecureConnection;
            isWebSocketRequest = request.IsWebSocketRequest;
            keepAlive = request.KeepAlive;
            queryString = request.QueryString;
            rawUrl = request.RawUrl;
            if (request.HasEntityBody) {
                using (System.IO.Stream body = request.InputStream) {
                    using (System.IO.StreamReader reader = new System.IO.StreamReader(body, request.ContentEncoding)) {
                        bodyData = reader.ReadToEnd();
                    }
                }
            }
        }
    }
}