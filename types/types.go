// types contains user defined types and constants
package types 

const (
	// Size of Data in KilobBytes
	KB = 1000 * 1

	// Size of Data in KibiBytes, same as 2^10
	KiB = 1 << (10 * 1)

	// Size of Data in MegaBytes
	MB = 1000 * KB

	// Size of Data in MebiBytes, same as 2^20
	MiB = 1 << 20 

	// Size of Data in GigaBytes
	GB = 1000 * MB

	// Size of Data in GibiBytes, same as 2^30
	GiB = 1 << 30

	// Size of Data in TeraBytes
	TB = 1000 * GB

	// Size of Data in TebiBytes, same as 2^40
	TiB = 1 << 40
)