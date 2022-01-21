// Code generated by `generate`. DO NOT EDIT.

package oxide

import "net/http"

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.oxide.computer for example.
	server string

	// Client is the *http.Client for performing requests.
	client *http.Client

	// token is the API token used for authentication.
	token         string
	Racks         *RacksService
	Sleds         *SledsService
	Organizations *OrganizationsService
	Disks         *DisksService
	Projects      *ProjectsService
	Users         *UsersService
	Roles         *RolesService
	Instances     *InstancesService
	Sagas         *SagasService
	Metrics       *MetricsService
	Vpcs          *VpcsService
	Subnets       *SubnetsService
	Firewall      *FirewallService
	Routers       *RoutersService
	Routes        *RoutesService
	Hidden        *HiddenService
}

type RacksService service

type SledsService service

type OrganizationsService service

type DisksService service

type ProjectsService service

type UsersService service

type RolesService service

type InstancesService service

type SagasService service

type MetricsService service

type VpcsService service

type SubnetsService service

type FirewallService service

type RoutersService service

type RoutesService service

type HiddenService service
