#ifndef go_patricia_h
#define go_patricia_h

#include <stdlib.h> /* free, atol, calloc */
#include <string.h> /* strlen */
#include <stdbool.h>
#include <time.h>

#include "patricia.h"
#include <regex.h>

#define MAX_CH_KEY 32
#define IPv4LEN 4
#define IPv6LEN 16

static size_t MAX_DEC = 1092;

typedef struct _mask_t {
    uint32_t             lower;
    uint32_t             upper;
} net_mask_t;

typedef struct _int_t {
    uint32_t*            ints;
    size_t               num_ints;
} net_interface_t;

typedef struct _int64_t {
    uint64_t*            ints;
    size_t               num_ints;
} net_interface64_t;

typedef struct strary {
    char    **v;
    size_t  *s;
    int64_t n;
} strary_t;

typedef struct _geo_info_t {
    uint8_t              ip[16];
    uint32_t             city;
    uint16_t             region;
    uint16_t             country;
} geo_info_t;

// Omitted leading _ in these struct names because you refer to them in Go via
// "struct_" + struct-name, which became struct__v4_cidr, which I didn't like.
typedef struct v4_cidr {
    uint8_t              ip[IPv4LEN];
    uint8_t              mask[IPv4LEN];
} v4_cidr_t;
typedef struct v6_cidr {
    uint8_t              ip[IPv6LEN];
    uint8_t              mask[IPv6LEN];
} v6_cidr_t;

typedef struct _v4_cidr_ary {
    v4_cidr_t*           cidrs;
    size_t               num_cidrs;
} v4_cidr_ary_t;

typedef struct _v6_cidr_ary {
    v6_cidr_t*           cidrs;
    size_t               num_cidrs;
} v6_cidr_ary_t;

typedef struct _flow_tag_t {
    uint16_t             tcp_flags;
    net_interface_t*     protos;
    net_interface_t*     ports;
    net_interface_t*     interfaces;
    net_interface_t*     asns;
    net_interface_t*     nexthop_asns;
    char*                value;
    regex_t*             bgp_aspath;
    size_t               bgp_aspath_count;
    regex_t*             bgp_community;
    size_t               bgp_community_count;
    net_interface64_t*   macs;
    net_interface_t*     countries;
    v4_cidr_ary_t*       v4cidrs;
    v6_cidr_ary_t*       v6cidrs;
    struct _flow_tag_t*  next;
} flow_tag_t;

typedef struct _packed_geo_t {
    geo_info_t*         data;
    size_t              length;
} packed_geo_t;

char* strary_item(strary_t *s, int64_t position);
size_t strary_size(strary_t *s, int64_t position);
void free_strary(strary_t *s);

int load_asn(const char *db_file, patricia_tree_t* tree);
int load_geo(const char *ipr, uint16_t cntd, uint16_t regd, uint32_t citd, patricia_tree_t* tree);
int load_flow_tag(char* key, flow_tag_t *tag, patricia_tree_t* tree);
void close_tree(patricia_tree_t* tree);
patricia_node_t* best_from_ipv4_network_byte_order(patricia_tree_t* tree, uint32_t ip);
patricia_node_t* best_from_ipv4_host_byte_order(patricia_tree_t* tree, uint32_t ip);
patricia_node_t* best_from_ipv6_network_byte_order(patricia_tree_t* tree, struct in6_addr);

flow_tag_t* new_flow_tag(uint16_t tcp_flags, const char* val,
        uint32_t* protos, size_t num_protos,
        uint32_t* interfaces, size_t num_interfaces,
        uint32_t* ports, size_t num_ports,
        uint32_t* asns, size_t num_asns,
        uint32_t* nexthop_asns, size_t num_nexthop_asns,
        const void *bgp_aspath[], size_t num_bgp_aspaths,
        const void *bgp_community[], size_t num_bgp_communities,
        uint64_t* macs, size_t num_macs,
        uint32_t* countries, size_t num_countries,
        v4_cidr_t* v4cidrs, size_t num_v4_cidrs,
        v6_cidr_t* v6cidrs, size_t num_v6_cidrs);

v4_cidr_ary_t* make_cidrv4_list(v4_cidr_t* v4cidrs, size_t num_v4_cidrs);
v6_cidr_ary_t* make_cidrv6_list(v6_cidr_t* v6cidrs, size_t num_v6_cidrs);
void free_flow_tag(void* node);
void close_tag_tree(patricia_tree_t* tree);

char* check_tag(void* node, uint16_t tcp_flags, uint16_t port, uint8_t proto,
        uint16_t interface, uint32_t asn, uint32_t nexthop_asn,
        const char *bgp_aspath, const char *bgp_community,
        uint64_t mac, uint32_t country,
        uint32_t nexthopv4, uint8_t* nexthopv6,
        size_t* result_len, int* tagOverflow);

int check_v4_cidrs(v4_cidr_ary_t* v4cidrs, uint8_t* nexthopv4);
int check_v6_cidrs(v6_cidr_ary_t* v6cidrs, uint8_t* nexthopv6);
int contains_v4(v4_cidr_t* v4cidr, uint8_t* nexthopv4);
int contains_v6(v6_cidr_t* v6cidr, uint8_t* nexthopv6);
net_interface_t* make_int_list(uint32_t* interfaces, size_t num_ints);
net_interface64_t* make_int64_list(uint64_t* interfaces, size_t num_ints);
int check_int(net_interface_t* interfaces, uint32_t check);
int check_int64(net_interface64_t* interfaces, uint64_t check);

int check_addr(const net_mask_t* mask, uint32_t addr);
net_mask_t* make_mask(uint32_t net_addr, uint32_t net_mask);

int check_regex(const regex_t *regex, size_t regex_count, const char *str);

// Helper functions for geo
uint16_t get_country(geo_info_t* data);
uint16_t get_region(geo_info_t* data);
uint32_t get_city(geo_info_t* data);

geo_info_t* new_geo(uint8_t ip[16], uint16_t country, uint16_t region, uint32_t city);
int save_packed_geo(geo_info_t** geo, size_t num_rows, const char* file);
geo_info_t* find_in_packed(packed_geo_t* packed, uint8_t search[16]);
packed_geo_t* new_packed_geo_from_file(const char* file);
int close_packed_geo(packed_geo_t* packed);
char* tmp_template(const char* filename);
#endif
