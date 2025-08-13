package ranger

type Access struct {
	Type      string `json:"type"`
	IsAllowed bool   `json:"isAllowed"`
}

type PolicyItem struct {
	Accesses      []Access `json:"accesses"`
	Users         []string `json:"users"`
	Groups        []string `json:"groups"`
	DelegateAdmin bool     `json:"delegateAdmin"`
}

type ResourceType struct {
	Values      []string `json:"values"`
	IsExcludes  bool     `json:"isExcludes"`
	IsRecursive bool     `json:"isRecursive"`
}

type Resources struct {
	Topic           *ResourceType `json:"topic,omitempty"`
	Consumergroup   *ResourceType `json:"consumergroup,omitempty"`
	Cluster         *ResourceType `json:"cluster,omitempty"`
	TransactionalId *ResourceType `json:"transactionalid,omitempty"`
	DelegationToken *ResourceType `json:"delegationtoken,omitempty"`
}

type Policy struct {
	ID             int          `json:"id,omitempty"`
	GUID           string       `json:"guid,omitempty"`
	IsEnabled      bool         `json:"isEnabled"`
	Version        int          `json:"version,omitempty"`
	Service        string       `json:"service"`
	Name           string       `json:"name"`
	PolicyType     int          `json:"policyType"`
	PolicyPriority int          `json:"policyPriority,omitempty"`
	Description    string       `json:"description,omitempty"`
	IsAuditEnabled bool         `json:"isAuditEnabled"`
	Resources      Resources    `json:"resources"`
	PolicyItems    []PolicyItem `json:"policyItems,omitempty"`
	ServiceType    string       `json:"serviceType"`
	IsDenyAllElse  bool         `json:"isDenyAllElse"`
}

type Service struct {
	ID               int               `json:"id,omitempty"`
	GUID             string            `json:"guid,omitempty"`
	IsEnabled        bool              `json:"isEnabled"`
	DisplayName      string            `json:"displayName,omitempty"`
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	TagService       string            `json:"tagService,omitempty"`
	Configs          map[string]string `json:"configs,omitempty"`
	PolicyVersion    int               `json:"policyVersion,omitempty"`
	PolicyUpdateTime int64             `json:"policyUpdateTime,omitempty"`
	TagVersion       int64             `json:"tagVersion,omitempty"`
	TagUpdateTime    int64             `json:"tagUpdateTime,omitempty"`
}
