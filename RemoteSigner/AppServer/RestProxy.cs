using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Newtonsoft.Json;
using RemoteSigner.Log;
using RemoteSigner.Models;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.AppServer {
    class RestProxy {
        static readonly Type[] baseTypes = { typeof(int), typeof(float), typeof(long), typeof(double) };

        Dictionary<string, ProxyMethod> proxyMethods;
        object instance;
        Type classType;

        internal RestProxy(Type restClass, Dictionary<string, Object> injectables) {
            instance = Activator.CreateInstance(restClass);
            classType = restClass;
            Attribute t = restClass.GetCustomAttribute(typeof(REST));
            Logger.Log("RestProxy", $"Creating proxy for {restClass.Name}");
            proxyMethods = new Dictionary<string, ProxyMethod>();
            REST trest = (REST)t;

            // Search Injectables to inject
            FieldInfo[] fields = restClass.GetFields(BindingFlags.NonPublic | BindingFlags.Instance);
            foreach (FieldInfo field in fields) {
                if (field.GetCustomAttribute(typeof(Inject)) != null) {
                    Type ft = field.FieldType;
                    Object injectableInstance;
                    if (injectables.ContainsKey(ft.FullName)) {
                        injectableInstance = injectables[ft.FullName];
                    } else {
                        Logger.Log("RestProxy", $"Creating injectable instance for class {ft.FullName}");
                        injectableInstance = Activator.CreateInstance(ft);
                        injectables.Add(ft.Name, injectableInstance);
                    }
                    field.SetValue(instance, injectableInstance);
                }
            }

            // Search Methods to Map
            MethodInfo[] methods = restClass.GetMethods();
            foreach (var methodInfo in methods) {
                proxyMethods.Add(methodInfo.Name, new ProxyMethod(methodInfo));
                foreach (var paramInfo in methodInfo.GetParameters()) {
                    // Default to body param
                    ProxyParameterRestType restType = ProxyParameterRestType.BODY;
                    string lookName = paramInfo.Name;
                    Attribute p;

                    if ((p = paramInfo.GetCustomAttribute(typeof(QueryParam))) != null) {
                        restType = ProxyParameterRestType.QUERY;
                        lookName = ((QueryParam)p).ParamName ?? paramInfo.Name;
                    }

                    Func<string, object> parser;
                    Type baseType;

                    if ((baseType = GetBaseType(paramInfo.ParameterType)) != null) {
                        parser = x => {
                            object[] dp = { x, Activator.CreateInstance(baseType) };
                            baseType.InvokeMember("TryParse", BindingFlags.InvokeMethod, null, null, dp);
                            return dp[1];
                        };
                    } else if (typeof(string).IsAssignableFrom(paramInfo.ParameterType)) {
                        parser = x => x;
                    } else {
                        parser = x => JsonConvert.DeserializeObject(x, paramInfo.ParameterType);
                    }

                    proxyMethods[methodInfo.Name].ProxyData.Add(new ProxyParameterData(restType, paramInfo.ParameterType, lookName, parser));
                }
            }
        }

        public object CallMethod(string methodName, RestRequest request) {
            return proxyMethods[methodName].Method.Invoke(instance, ParamsBuilder(methodName, request));
        }

        static Type GetBaseType(Type t) {
            try {
                return baseTypes.Where(x => x.IsAssignableFrom(t)).ElementAt(0);
            } catch (ArgumentOutOfRangeException) {
                return null;
            }
        }

        object[] ParamsBuilder(string methodName, RestRequest request) {
            List<ProxyParameterData> proxyData = proxyMethods[methodName].ProxyData;
            object[] callParams = new object[proxyData.Count];
            for (int i = 0; i < proxyData.Count; i++) {

                string parseData = null;

                switch (proxyData[i].RestType) {
                    case ProxyParameterRestType.BODY:
                        parseData = request.BodyData;
                        break;
                    case ProxyParameterRestType.QUERY:
                        parseData = request.QueryString[proxyData[i].LookName];
                        break;
                }

                if (parseData != null) {
                    callParams[i] = proxyData[i].Parse(parseData);
                }
            }

            return callParams;
        }

        class ProxyMethod {
            readonly MethodInfo method;
            readonly List<ProxyParameterData> proxyData;

            public MethodInfo Method {
                get { return method; }
            }

            public List<ProxyParameterData> ProxyData {
                get { return proxyData; }
            }

            public ProxyMethod(MethodInfo method) {
                this.method = method;
                proxyData = new List<ProxyParameterData>();
            }
        }

        class ProxyParameterData {
            ProxyParameterRestType restType;
            readonly Type parameterType;
            readonly string lookName;
            readonly Func<string, object> parse;

            public ProxyParameterRestType RestType {
                get { return restType; }
            }

            public Type ParameterType {
                get { return parameterType; }
            }

            public string LookName {
                get { return lookName; }
            }

            public Func<string, object> Parse {
                get { return parse; }
            }

            public ProxyParameterData(ProxyParameterRestType restType, Type parameterType, string lookName, Func<string, object> parse) {
                this.restType = restType;
                this.parameterType = parameterType;
                this.lookName = lookName;
                this.parse = parse;
            }
        }

        enum ProxyParameterRestType {
            BODY,
            QUERY
        }
    }
}
