#include "go_patricia.h"
// Pick up the contains_v6 prototype.
#include "_cgo_export.h"

#include <stdio.h> /* printf */
#include <stdlib.h> /* exit */
#include <arpa/inet.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <fcntl.h>

void _alloc_strary(char ***v, size_t n) {
    int i;
    *v = (char **)calloc(n, sizeof(char *));
}

void _free_strary(char ***v, size_t n) {
    int i;
    free(*v);
}

int load_asn(const char *db_file, patricia_tree_t* tree)
{
    char buf[256];
    FILE *in;
    int count = 0;
    patricia_node_t *node;
    char asn[256];

    int tm = time(NULL);

    srand(tm);

    in = fopen(db_file, "r");
    if (in) {
        while (fgets(buf, 200, in)) {
            char *c = strchr(buf, ',');

            if (c && *c == ',') {
                *c = '\0';

                node = make_and_lookup(tree, c+1);
                if (!node) {
                    printf("couldn't make asn node for asn [%s] prefix [%s]\n", asn, buf+3);
                } else {
                    // Duplicates are allowed.  See discussion starting here:
                    // https://kentik.slack.com/archives/dev-backend/p1473964361000509
                    if (node->data != NULL) {
                      free(node->data);
                    }
                    count++;
                    node->data = strdup(buf);
                }
            }
        }
        fclose(in);
    }

    printf("loaded %d asn records\n", count);
    return count;
}

geo_info_t* new_geo(uint8_t ip[16], uint16_t country, uint16_t region, uint32_t city) {
    geo_info_t* geo = (geo_info_t*)malloc(sizeof(geo_info_t));

    if (geo) {
        memcpy(geo->ip, ip, 16);
        geo->country = country;
        geo->region = region;
        geo->city = city;
    }

    return geo;
}

packed_geo_t* new_packed_geo_from_file(const char* file) {

    packed_geo_t* packed = NULL;
    geo_info_t* geo = NULL;
    int input_fd;
    struct stat st;
    int length;

    input_fd = open (file, O_RDONLY);
    if (input_fd == -1) {
        return NULL;
    }

    fstat(input_fd, &st);
    geo = (geo_info_t*)mmap(0, st.st_size, PROT_READ, MAP_SHARED, input_fd, 0);
    close(input_fd);

    if (geo==(void*) -1){
        return NULL;
    }

    length = st.st_size / sizeof(geo_info_t);
    if (length == 0) {
        return NULL;
    }

    packed = (packed_geo_t*)malloc(sizeof(packed_geo_t));
    if (packed == NULL) {
        return NULL;
    }
    packed->length = length;
    packed->data = geo;
    return packed;
}


char* tmp_template(const char* filename) {
    const char* suffix = ".XXXXXX";
    size_t buflen = strlen(filename) + strlen(suffix) +1;
    char* buf = calloc(1, buflen);
    if (buf == NULL) {
        return NULL;
    }
    if (snprintf(buf, buflen, "%s%s", filename, suffix) >= buflen) {
        free(buf);
        return NULL;
    }
    return buf;
}

int save_packed_geo(geo_info_t** geo, size_t num_rows, const char* file) {
    size_t i;
    size_t count = 0;
    int fd;
    geo_info_t* file_memory = NULL;
    size_t file_len = sizeof(geo_info_t) * num_rows;
    char* tmpfile = NULL;

    if (file == NULL) {
        return -1;
    }

    tmpfile = tmp_template(file);
    if (tmpfile == NULL) {
        return -2;
    }

    fd = mkstemp(tmpfile);
    if (fd < 0) {
        free(tmpfile);
        return fd;
    }

    if (ftruncate(fd, file_len) < 0) {
        free(tmpfile);
        close (fd);
        return -3;
    }

    file_memory = (geo_info_t*)mmap (NULL, file_len, PROT_WRITE, MAP_SHARED, fd, 0);
    if (file_memory == MAP_FAILED) {
        free(tmpfile);
        close (fd);
        return -4;
    }

    // copy data
    for (i=0; i<num_rows; i++) {
        file_memory[i] = **geo++;
        count++;
    }

    munmap (file_memory, file_len);
    close (fd);

    if (rename(tmpfile, file) != 0) {
        free(tmpfile);
        return -5;
    }

    free(tmpfile);
    return count;
}

int close_packed_geo(packed_geo_t* packed) {
    size_t file_len = sizeof(geo_info_t) * packed->length;

    munmap (packed->data, file_len);
    free(packed);

    return 0;
}

void close_tree(patricia_tree_t* tree) {
    if (tree != NULL) {
        Destroy_Patricia(tree, (void *)0);
    }
}

patricia_node_t* best_from_ipv4_network_byte_order(patricia_tree_t* tree, uint32_t ip) {

    patricia_node_t*       node = NULL;
    struct in_addr         in;
    prefix_t               prefix;

    prefix.family = AF_INET;
    prefix.bitlen = 32;
    prefix.ref_count = 1000;

    bcopy(&ip, &in, sizeof(in));
    bcopy(&ip, &prefix.add.sin, sizeof(in));
    node = patricia_search_best(tree, &prefix);

    return node;
}

patricia_node_t* best_from_ipv4_host_byte_order(patricia_tree_t* tree, uint32_t ip) {
    uint32_t nip = htonl(ip);
    return best_from_ipv4_network_byte_order(tree, nip);
}

patricia_node_t *best_from_ipv6_network_byte_order(patricia_tree_t *tree, struct in6_addr addr) {
    prefix_t prefix;

    prefix.family = AF_INET6;
    prefix.bitlen = 128;
    prefix.ref_count = 1000;
    memcpy(&prefix.add.sin6, &addr, sizeof(addr));
    return patricia_search_best(tree, &prefix);
}

int load_flow_tag(char* key, flow_tag_t *tag, patricia_tree_t* tree) {

    patricia_node_t*      node = NULL;
    flow_tag_t*           old_tag = NULL;

    node = make_and_lookup(tree, key);
    if (!node) {
        printf("couldn't make node for prefix [%s]\n", key);
        return 0;
    } else {
        if (node->data == NULL) {
            node->data = tag;
        } else {
            old_tag = (flow_tag_t*)node->data;

            while (old_tag->next != NULL) {
                old_tag = old_tag->next;
            }
            old_tag->next = tag;
        }
    }

    return 1;
}

#define FLOW_ALLOC_REGEX(regexbuf, regexcount, str, count)              \
    do {                                                                \
        int i;                                                          \
        regexbuf = calloc(sizeof((regexbuf)[0]), count);                \
        for (i = 0, regexcount = 0; i < count; i++) {                   \
            int regerr = regcomp(&(regexbuf)[regexcount],               \
                    (str)[i], REG_EXTENDED|REG_NOSUB);                  \
            if (regerr) {                                               \
                char regerrstr[512] = {0};                              \
                regerror(regerr, &(regexbuf)[regexcount],               \
                        regerrstr, sizeof(regerrstr));                  \
                fprintf(stderr, "could not compile regex \"%p\":%s\n",  \
                        (str)[i], regerrstr);                           \
                continue;                                               \
            }                                                           \
            regexcount++;                                               \
        }                                                               \
    } while(0)

#define FLOW_FREE_REGEX(regexbuf, regexcount)                           \
    do {                                                                \
        int i;                                                          \
        for (i = 0; i < regexcount; i++) {                              \
            regfree(&(regexbuf)[i]);                                    \
        }                                                               \
        free(regexbuf);                                                 \
    } while(0)

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
        v6_cidr_t* v6cidrs, size_t num_v6_cidrs) {

    flow_tag_t* tag = (flow_tag_t*)calloc(1, sizeof(flow_tag_t));
    if (!tag) {
        return NULL;
    }

    tag->tcp_flags = tcp_flags;
    tag->value = strdup(val);
    tag->next = NULL;
    tag->interfaces = (num_interfaces > 0)? make_int_list(interfaces, num_interfaces): NULL;
    tag->ports = (num_ports > 0)? make_int_list(ports, num_ports): NULL;
    tag->protos = (num_protos > 0)? make_int_list(protos, num_protos): NULL;
    tag->asns = (num_asns > 0)? make_int_list(asns, num_asns): NULL;
    tag->nexthop_asns = (num_nexthop_asns > 0)? make_int_list(nexthop_asns, num_nexthop_asns): NULL;
    tag->macs = (num_macs > 0)? make_int64_list(macs, num_macs): NULL;
    tag->countries = (num_countries > 0)? make_int_list(countries, num_countries): NULL;
    tag->v4cidrs = make_cidrv4_list(v4cidrs, num_v4_cidrs);
    tag->v6cidrs = make_cidrv6_list(v6cidrs, num_v6_cidrs);

    FLOW_ALLOC_REGEX(tag->bgp_aspath, tag->bgp_aspath_count, bgp_aspath, num_bgp_aspaths);
    FLOW_ALLOC_REGEX(tag->bgp_community, tag->bgp_community_count, bgp_community, num_bgp_communities);

    return tag;
}

v4_cidr_ary_t* make_cidrv4_list(v4_cidr_t* cidrs, size_t num_cidrs) {
    if (num_cidrs == 0) {
        return NULL;
    }

    v4_cidr_ary_t* list = (v4_cidr_ary_t*)calloc(1, sizeof(v4_cidr_ary_t));

    list->cidrs = (v4_cidr_t*)calloc(num_cidrs, sizeof(v4_cidr_t));
    list->num_cidrs = num_cidrs;

    memmove(list->cidrs, cidrs, sizeof(v4_cidr_t) * num_cidrs);

    return list;
}

v6_cidr_ary_t* make_cidrv6_list(v6_cidr_t* cidrs, size_t num_cidrs) {
    if (num_cidrs == 0) {
        return NULL;
    }

    v6_cidr_ary_t* list = (v6_cidr_ary_t*)calloc(1, sizeof(v6_cidr_ary_t));

    list->cidrs = (v6_cidr_t*)calloc(num_cidrs, sizeof(v6_cidr_t));
    list->num_cidrs = num_cidrs;

    memmove(list->cidrs, cidrs, sizeof(v6_cidr_t) * num_cidrs);

    return list;
}

void free_flow_tag(void* node) {

    flow_tag_t* tag = (flow_tag_t*)node;
    while (tag != NULL) {
        flow_tag_t* next_tag = tag->next;

        if (tag->value != NULL) {
            free(tag->value);
        }

        if (tag->interfaces != NULL) {
            free(tag->interfaces->ints);
            free(tag->interfaces);
        }

        if (tag->ports != NULL) {
            free(tag->ports->ints);
            free(tag->ports);
        }

        if (tag->protos != NULL) {
            free(tag->protos->ints);
            free(tag->protos);
        }

        if (tag->asns != NULL) {
            free(tag->asns->ints);
            free(tag->asns);
        }

        if (tag->nexthop_asns != NULL) {
            free(tag->nexthop_asns->ints);
            free(tag->nexthop_asns);
        }

        if (tag->macs != NULL) {
            free(tag->macs->ints);
            free(tag->macs);
        }

        if (tag->countries != NULL) {
            free(tag->countries->ints);
            free(tag->countries);
        }

        if (tag->v4cidrs != NULL) {
            free(tag->v4cidrs->cidrs);
            free(tag->v4cidrs);
        }
        if (tag->v6cidrs != NULL) {
            free(tag->v6cidrs->cidrs);
            free(tag->v6cidrs);
        }

        FLOW_FREE_REGEX(tag->bgp_aspath, tag->bgp_aspath_count);
        FLOW_FREE_REGEX(tag->bgp_community, tag->bgp_community_count);

        free(tag);
        tag = next_tag;
    }
}

void close_tag_tree(patricia_tree_t* tree) {
    if (tree != NULL) {
        Destroy_Patricia(tree, free_flow_tag);
    }
}

char* check_tag(void* node, uint16_t tcp_flags, uint16_t port, uint8_t proto,
        uint16_t interface, uint32_t asn, uint32_t nexthop_asn,
        const char *bgp_aspath, const char *bgp_community,
        uint64_t mac, uint32_t country,
        uint32_t nexthopv4, uint8_t *nexthopv6,
        size_t* o_result_len, int* tagOverflow) {

    flow_tag_t*                  tag = (flow_tag_t*)node;
    char*                        results = NULL;

    *tagOverflow = 0;
    *o_result_len = 0;
    if (!tag) {
        return NULL;
    }

    results = (char*)calloc(MAX_DEC, sizeof(char*));
    memset(results, ' ', MAX_DEC * sizeof(char*));
    while (tag != NULL) {

        int good = 1;

        if (good && tag->tcp_flags && !(tag->tcp_flags & tcp_flags)) { good = 0; }
        if (good && tag->protos && !check_int(tag->protos, proto)) { good = 0; }
        if (good && tag->interfaces && !check_int(tag->interfaces, interface)) { good = 0; }
        if (good && tag->ports && !check_int(tag->ports, port)) { good = 0; }
        if (good && tag->asns && !check_int(tag->asns, asn)) { good = 0; }
        if (good && tag->nexthop_asns && !check_int(tag->nexthop_asns, nexthop_asn)) { good = 0; }
        if (good && tag->macs && !check_int64(tag->macs, mac)) { good = 0; }
        if (good && tag->countries && !check_int(tag->countries, country)) { good = 0; }
        if (good && tag->v4cidrs && nexthopv4 != 0 && !check_v4_cidrs(tag->v4cidrs, (uint8_t*) &nexthopv4)) { good = 0; }
        if (good && tag->v6cidrs && nexthopv4 == 0 && nexthopv6 != NULL && !check_v6_cidrs(tag->v6cidrs, nexthopv6)) { good = 0; }
        if (good && !check_regex(tag->bgp_aspath, tag->bgp_aspath_count, bgp_aspath)) {good = 0;}
        if (good && !check_regex(tag->bgp_community, tag->bgp_community_count, bgp_community)) {good = 0;}

        if (good) {
            size_t newlen = strlen(tag->value);
            if ((*o_result_len + newlen) >= MAX_DEC) {
                *tagOverflow = 1;
                return results;
            }

            if (*o_result_len == 0) {
                memcpy(results + *o_result_len, tag->value, newlen);
            } else {
                memcpy(results + 1 + *o_result_len, tag->value, newlen);
                newlen += 1;
            }
            *o_result_len += newlen;
        }

        tag = tag->next;
    }

    return results;
}

// Does nexthop match any of the cidrs in cidr_ary?
int check_v4_cidrs(v4_cidr_ary_t* cidr_ary, uint8_t* nexthop) {
    int i;

    for (i = 0; i < cidr_ary->num_cidrs; i++) {
        if (contains_v4(&cidr_ary->cidrs[i], nexthop)) {
            return 1;
        }
    }
    return 0;
}

// Does nexthop match any of the cidrs in cidr_ary?
int check_v6_cidrs(v6_cidr_ary_t* cidr_ary, uint8_t* nexthop) {
    int i;

    for (i = 0; i < cidr_ary->num_cidrs; i++) {
        if (contains_v6(&cidr_ary->cidrs[i], nexthop)) {
            return 1;
        }
    }
    return 0;
}

// Modeled on Go's net.IPNet.Contains,
// https://golang.org/src/net/ip.go?s=10913:10949#L456
// Does cidr contain (v4) ip?
int contains_v4(v4_cidr_t* cidr, uint8_t* ip) {
    int i;

    uint8_t* nn = cidr->ip;   // network number
    uint8_t* m = cidr->mask;  // mask

    for (i = 0; i < IPv4LEN; i++) {
        if ((nn[i] & m[i]) != (ip[i] & m[i])) {
            return 0;
        }
    }
    return 1;
}

// Does cidr contain (v6) ip?
int contains_v6(v6_cidr_t* cidr, uint8_t* ip) {
    int i;

    uint8_t* nn = cidr->ip;   // network number
    uint8_t* m = cidr->mask;  // mask

    for (i = 0; i < IPv6LEN; i++) {
        if ((nn[i] & m[i]) != (ip[i] & m[i])) {
            return 0;
        }
    }
    return 1;
}

int check_regex(const regex_t *regex, size_t regex_count, const char *str) {
    int i;

    /* skip null strings and empty regex buf */
    if (regex_count == 0)
        return 1;
    else if (str[0] == '\0')
        return 0;

    for (i = 0; i < regex_count; i++) {
        if (regexec(&regex[i], str, 0, NULL, 0) == 0) {
            return 1;
        }
    }
    return 0;
}

int check_int(net_interface_t* interfaces, uint32_t check) {
    size_t i;

    for (i=0; i<interfaces->num_ints; i++) {
        if (interfaces->ints[i] == check) {
            return 1;
        }
    }

    return 0;
}

int check_int64(net_interface64_t* interfaces, uint64_t check) {
    size_t i;

    for (i=0; i<interfaces->num_ints; i++) {
        if (interfaces->ints[i] == check) {
            return 1;
        }
    }

    return 0;
}

int check_addr(const net_mask_t* mask, uint32_t addr) {
    if (addr >= mask->lower && addr <= mask->upper) {
        return 1;
    }
    return 0;
}

net_interface_t* make_int_list(uint32_t* ints, size_t num_ints) {

    net_interface_t*           interfaces = (net_interface_t*)calloc(1, sizeof(net_interface_t));
    int                        i;

    interfaces->ints = (uint32_t*)calloc(num_ints, sizeof(uint32_t));
    interfaces->num_ints = num_ints;

    for (i=0; i<num_ints; i++) {
        interfaces->ints[i] = *ints++;
    }

    return interfaces;
}

net_interface64_t* make_int64_list(uint64_t* ints, size_t num_ints) {

    net_interface64_t*         interfaces = (net_interface64_t*)calloc(1, sizeof(net_interface64_t));
    int                        i;

    interfaces->ints = (uint64_t*)calloc(num_ints, sizeof(uint64_t));
    interfaces->num_ints = num_ints;

    for (i=0; i<num_ints; i++) {
        interfaces->ints[i] = *ints++;
    }

    return interfaces;
}

net_mask_t* make_mask(uint32_t network_addr, uint32_t mask_bits) {

    uint32_t mask_addr = ~(0xffffffff >> mask_bits);
    net_mask_t* mask = (net_mask_t*)calloc(1, sizeof(net_mask_t));

    if (!mask) {
        return NULL;
    }

    mask->lower = (network_addr & mask_addr);
    mask->upper = (mask->lower | (~mask_addr));

    return mask;
}

char* strary_item(strary_t *s, int64_t position) {
    if (position < s->n) {
        return s->v[position];
    }
    return NULL;
}

size_t strary_size(strary_t *s, int64_t position) {
    if (position < s->n) {
        return s->s[position];
    }
    return 0;
}

void free_strary(strary_t *s) {
    _free_strary(&s->v, s->n);
    free(s->s);
}

geo_info_t* find_in_packed(packed_geo_t* packed, uint8_t search[16]) {

    geo_info_t*        geo;
    geo_info_t*        next;
    int                first = 0;
    int                last;
    int                middle;

    if (packed == NULL) {
        return NULL;
    }

    last = packed->length - 1;
    middle = (first+last)/2;

    while (first <= last) {
        geo = &(packed->data[middle]);
        next = NULL;

        // Now, next is middle+1, since geo is middle.
        if (packed->length > middle+1) {
            next = &(packed->data[middle+1]);
        }

        if (geo != NULL && next != NULL) {
            int c = memcmp(geo->ip, search, 16);
            if (c == 0) {
                break;
            } else if (c < 0 && (middle+1 > last || memcmp(next->ip, search, 16) > 0)) {
                break;
            } else if (c < 0) {
                first = middle + 1;
            } else {
                last = middle - 1;
            }

            middle = (first + last)/2;
        } else {
            return NULL;
        }
    }

    if (first > last) {
        return NULL;
    }

    return geo;
}

// Helper functions for geo
uint16_t get_country(geo_info_t* geo) {
    return geo->country;
}

uint16_t get_region(geo_info_t* geo) {
    return geo->region;
}

uint32_t get_city(geo_info_t* geo) {
    return geo->city;
}
