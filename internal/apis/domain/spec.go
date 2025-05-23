package domain

type Api interface {
	GetZone(fqdn string) (*Zone, error)
	CreateZone(fqdn string) (*Zone, error)
	DeleteZone(fqdn string) (bool, error)

	GetRecord(zoneFqdn string, id int64) (*Record, error)
	CreateRecord(zoneFqdn, recordType, source, target string, ttl int64) (*Record, error)
	DeleteRecord(zoneFqdn string, id int64) (bool, error)
}
