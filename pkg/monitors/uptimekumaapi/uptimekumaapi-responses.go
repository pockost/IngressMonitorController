package uptimekumaapi

type UptimeKumaApiAuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"bearer"`
}

type UptimeKumaApiErrorResponse struct {
	Detail UptimeKumaApiErrorDetail `json:"detail"`
}

type UptimeKumaApiErrorDetail struct {
	Message string `json:"msg"`
	Type    string `json:"type"`
}

type UptimeKumaApiMonitorListResponse struct {
	Monitors []UptimeKumaApiMonitor `json:"monitors"`
}

type UptimeKumaApiMonitor struct {
	Id                  int    `json:"id,omitempty"`
	Type                string `json:"type"`
	Name                string `json:"name"`
	Interval            int    `json:"interval,omitempty"`
	RetryInterval       int    `json:"retryInterval,omitempty"`
	ResendInterval      int    `json:"resendInterval,omitempty"`
	MaxRetries          int    `json:"maxretries,omitempty"`
	Url                 string `json:"url"`
	Method              string `json:"method"`
	IgnoreTLS 	        bool   `json:"ignoreTls"`
	UpsideDown          bool   `json:"upsideDown"`
	MaxRedirects        int    `json:"maxredirects"`
	AcceptedStatusCodes []string `json:"accepted_statuscodes"`
	SSLExpire           bool   `json:"expiryNotification"`
}
