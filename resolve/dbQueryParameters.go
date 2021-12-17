package resolve

func getQueryParametersUnlocalizedIps() string {
	return "select DISTINCT d.ip from data d left JOIN locations l ON d.ip = l.ip WHERE l.country IS NULL"
}
