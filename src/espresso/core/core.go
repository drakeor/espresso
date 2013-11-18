package core

type Core struct {
	DB *DBO
	SESS *SESSOBJ
}

func InitCore(conf *Config) *Core {
	C := &Core{}	

	C.DB = InitDBO(conf.Params.DBHost, conf.Params.DBUsername, conf.Params.DBPassword, conf.Params.DBTable)
	C.DB.Prefix = conf.Params.DBPrefix
	
	C.SESS = InitSessions("cache/sessions")

	return C
}

func (C Core) Close() {
	C.DB.Close()
}
