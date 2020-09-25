## 超级链接口说明

#### 链注册

```go
Fcn:   "chainRegister",
Args:        
[
       "info",   // 链的名称等信息
	     "ip",     // 链的PAPP的IP地址
	     "serial", // 时间戳序列
	     "csr",    // 证书请求
	     "orgCACert", // 链的各个组织的证书
]
```

#### 获取链的相关信息

```go
Fcn:   "getChainInfo",
Args:        
[
       "chainID",   // 链的ID
]
```

#### 删除链的相关信息

```go
Fcn:   "deleteChain",
Args:        
[
       "chainID",   // 链的ID
]
```

#### 修改链的所有组织证书

```go
Fcn:   "setChainOrgCACert",
Args:        
[
       "chainID",    // 链的ID
	     "orgCACert",  // 修改后的证书信息
]
```

#### 获取链的组织证书

```go
Fcn:   "getChainOrgCACert",
Args:        
[
       "chainID",    // 链的ID
]
```

#### 更新链的某个组织证书

```go
Fcn:   "updateOrgCACert",
Args:        
[
       "chainID",    // 链的ID
       "orgName",    // 组织名称
       "newCert",    // 更新后的证书
]
```

#### getChainInfo后的返回结构

```go
type ReturnToRegister struct {
	ID            string `json:"id"`             //Chain ID
	CERT          string `json:"cert"`           //Chain-PeerAPP-Cert
	ROOTCERT      string `json:"root_cert"`      //Superchain RootCert
}
```

#### 