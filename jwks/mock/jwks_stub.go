package mock

import (
	"github.com/ONSdigital/dp-identity-api/v2/jwks"
)

var KeySetOne = jwks.JSONKey{
	E:   "AQAB",
	Kid: "j+diD4wBP/VZ4+X51XGRdI8Vi0CNV0OpEefKl1ge3A8=",
	Kty: "RSA",
	N:   "vBvi--N-F9MQO81xh71jIbkx81w4_sGhbztTJgIdhycV-lMzG6y3dMBWo9eRsFJuRs3MUFElmRrTVxc7EPWNQGQjUyPFW0_CnPPoGBCwgCyWtpNs5EHAkCHXsfryHb6LbJxH9LEbwOQCHR25_Bnqo_NeXSBJtvUabq3cTUgdOPc61Hskq-m19M1u7u1xu7b5DHD308Qyz3OhaEHx3cLL2za-mKxHe0VDe3sa5UfdaliTdBypFWJgNl6TsxF_G83fksgb3bVchzW45pu4dEhtNLqgXejH2-GwU8YRaAguKGW7dO_v-5uwLgDYQG9wgtAwLIMiXsFU7muig2pJEtlG2w",
}

var KeySetTwo = jwks.JSONKey{
	E:   "AQAB",
	Kid: "Oe/15Omy/K78yrUh2EI6xiQSRyeD5f8D/bcI/UphRR8=",
	Kty: "RSA",
	N:   "6MMhL-GcDj8LspuAes_ZycMTOYUkjURF-3z5vFtn0roie0LlcSgXN9i7VEsU7a-CTdqzBXhm_D4Yu9-RcVYJb8upyzWfrK53l4UoeNrQGhbjZlGKqnuQgU20lRqhKPqmHtAejm81XaW2T-z_bM2oL4U4RjOe5KaWLCpFe8IB92aTFZfXsPcfSodwQar7Po4TsRMg3iqqTk-jxySSYgj72XaCD5c3TojC6rdD_ll1dVub0LYjMESDnfFXDY4iCakk1l5MBwgEXDabJuNajfAotrFUN6svfb9DlXYSR9E_VYKxeDGdWB3QPIoieA_hpNhSM4nhWUApamxaCRC6g4dJjQ",
}

var JWKSResponseData = `{
	"2a8vXmIK67ZZ3hFZ=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApU1DqxJalEmlznkrM+U4aSBMw9u5axcIqNeUq8+ZHo98uKy8Xy5zCOXfWx6KafOPJhbOZInFaSh9UMaluSSw11l/PR4KrGBFzJODQ+RMq6bHW6FlwwHSkMTSfQ0hwzO7y91BiZFmJnaUECf52H3QBApGT4TT060ri5zt1ygpliRwjLLlHW1XX0epzZH3ogrikn4i65e8w6uUcsGBhQvQQqiHvEpcgCQAB",
	"GRBevIroJzPBvaGa=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtvDfudfY9n+8sFJmHGFfgbKqKf8iiEcbvRXNMEi9qd2NGAekhdNJKdeW3sMSwR+sb4Ly6IypowCE2eueYk/GatzYyyolWny/Krdp0EWPT/PnK8Iq1FTIuHxFb08B8iLnH/2nKqgOjVvwEU4eSBh0YHKti2v77a+a4bnx6aOC2YkF2AyIRmbXAHaq4Js9u33X8gGMXZcVsxcSpG8Py/NJ3s+PLKebQFQAB"
}`

var JWKSData = map[string]string{
	"Oe/15Omy/K78yrUh2EI6xiQSRyeD5f8D/bcI/UphRR8=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6MMhL+GcDj8LspuAes/ZycMTOYUkjURF+3z5vFtn0roie0LlcSgXN9i7VEsU7a+CTdqzBXhm/D4Yu9+RcVYJb8upyzWfrK53l4UoeNrQGhbjZlGKqnuQgU20lRqhKPqmHtAejm81XaW2T+z/bM2oL4U4RjOe5KaWLCpFe8IB92aTFZfXsPcfSodwQar7Po4TsRMg3iqqTk+jxySSYgj72XaCD5c3TojC6rdD/ll1dVub0LYjMESDnfFXDY4iCakk1l5MBwgEXDabJuNajfAotrFUN6svfb9DlXYSR9E/VYKxeDGdWB3QPIoieA/hpNhSM4nhWUApamxaCRC6g4dJjQIDAQAB",
	"j+diD4wBP/VZ4+X51XGRdI8Vi0CNV0OpEefKl1ge3A8=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvBvi++N+F9MQO81xh71jIbkx81w4/sGhbztTJgIdhycV+lMzG6y3dMBWo9eRsFJuRs3MUFElmRrTVxc7EPWNQGQjUyPFW0/CnPPoGBCwgCyWtpNs5EHAkCHXsfryHb6LbJxH9LEbwOQCHR25/Bnqo/NeXSBJtvUabq3cTUgdOPc61Hskq+m19M1u7u1xu7b5DHD308Qyz3OhaEHx3cLL2za+mKxHe0VDe3sa5UfdaliTdBypFWJgNl6TsxF/G83fksgb3bVchzW45pu4dEhtNLqgXejH2+GwU8YRaAguKGW7dO/v+5uwLgDYQG9wgtAwLIMiXsFU7muig2pJEtlG2wIDAQAB",
}

var JWKSStubbed = &ManagerMock{
	JWKSGetKeysetFunc: func(_, _ string) (*jwks.JWKS, error) {
		return &jwks.JWKS{
			Keys: []jwks.JSONKey{
				KeySetOne,
				KeySetTwo,
			},
		}, nil
	},
	JWKSToRSAJSONResponseFunc: func(_ *jwks.JWKS) ([]byte, error) {
		return []byte(JWKSResponseData), nil
	},
}
