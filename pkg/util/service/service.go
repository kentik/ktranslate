package service

type Port struct {
	Number   uint32
	Protocol uint8
}

var Services = map[Port]string{
	Port{1, 6}:       "tcpmux",               // TCP Port Service Multiplexer [rfc-1078] | TCP Port Service Multiplexer
	Port{1, 17}:      "tcpmux",               // TCP Port Service Multiplexer
	Port{2, 6}:       "compressnet",          // Management Utility
	Port{2, 17}:      "compressnet",          // Management Utility
	Port{3, 6}:       "compressnet",          // Compression Process
	Port{3, 17}:      "compressnet",          // Compression Process
	Port{5, 6}:       "rje",                  // Remote Job Entry
	Port{5, 17}:      "rje",                  // Remote Job Entry
	Port{7, 132}:     "echo",                 // Missing description for echo
	Port{7, 6}:       "echo",                 // Missing description for echo
	Port{7, 17}:      "echo",                 // Missing description for echo
	Port{9, 132}:     "discard",              // sink null
	Port{9, 6}:       "discard",              // sink null
	Port{9, 17}:      "discard",              // sink null
	Port{11, 6}:      "systat",               // Active Users
	Port{11, 17}:     "systat",               // Active Users
	Port{13, 6}:      "daytime",              // Missing description for daytime
	Port{13, 17}:     "daytime",              // Missing description for daytime
	Port{15, 6}:      "netstat",              // Missing description for netstat
	Port{17, 6}:      "qotd",                 // Quote of the Day
	Port{17, 17}:     "qotd",                 // Quote of the Day
	Port{18, 6}:      "msp",                  // Message Send Protocol | Message Send Protocol (historic)
	Port{18, 17}:     "msp",                  // Message Send Protocol
	Port{19, 6}:      "chargen",              // ttytst source Character Generator | Character Generator
	Port{19, 17}:     "chargen",              // ttytst source Character Generator
	Port{20, 132}:    "ftp-data",             // File Transfer [Default Data] | FTP
	Port{20, 6}:      "ftp-data",             // File Transfer [Default Data]
	Port{20, 17}:     "ftp-data",             // File Transfer [Default Data]
	Port{21, 132}:    "ftp",                  // File Transfer [Control] | File Transfer Protocol [Control]
	Port{21, 6}:      "ftp",                  // File Transfer [Control]
	Port{21, 17}:     "ftp",                  // File Transfer [Control]
	Port{22, 132}:    "ssh",                  // Secure Shell Login | The Secure Shell (SSH) Protocol
	Port{22, 6}:      "ssh",                  // Secure Shell Login
	Port{22, 17}:     "ssh",                  // Secure Shell Login
	Port{23, 6}:      "telnet",               // Missing description for telnet
	Port{23, 17}:     "telnet",               // Missing description for telnet
	Port{24, 6}:      "priv-mail",            // any private mail system
	Port{24, 17}:     "priv-mail",            // any private mail system
	Port{25, 6}:      "smtp",                 // Simple Mail Transfer
	Port{25, 17}:     "smtp",                 // Simple Mail Transfer
	Port{26, 6}:      "rsftp",                // RSFTP
	Port{27, 6}:      "nsw-fe",               // NSW User System FE
	Port{27, 17}:     "nsw-fe",               // NSW User System FE
	Port{29, 6}:      "msg-icp",              // MSG ICP
	Port{29, 17}:     "msg-icp",              // MSG ICP
	Port{31, 6}:      "msg-auth",             // MSG Authentication
	Port{31, 17}:     "msg-auth",             // MSG Authentication
	Port{33, 6}:      "dsp",                  // Display Support Protocol
	Port{33, 17}:     "dsp",                  // Display Support Protocol
	Port{35, 6}:      "priv-print",           // any private printer server
	Port{35, 17}:     "priv-print",           // any private printer server
	Port{37, 6}:      "time",                 // timserver
	Port{37, 17}:     "time",                 // timserver
	Port{38, 6}:      "rap",                  // Route Access Protocol
	Port{38, 17}:     "rap",                  // Route Access Protocol
	Port{39, 6}:      "rlp",                  // Resource Location Protocol
	Port{39, 17}:     "rlp",                  // Resource Location Protocol
	Port{41, 6}:      "graphics",             // Missing description for graphics
	Port{41, 17}:     "graphics",             // Missing description for graphics
	Port{42, 6}:      "nameserver",           // name | Host Name Server
	Port{42, 17}:     "nameserver",           // Host Name Server
	Port{43, 6}:      "whois",                // nicname | Who Is
	Port{43, 17}:     "whois",                // nicname
	Port{44, 6}:      "mpm-flags",            // MPM FLAGS Protocol
	Port{44, 17}:     "mpm-flags",            // MPM FLAGS Protocol
	Port{45, 6}:      "mpm",                  // Message Processing Module [recv]
	Port{45, 17}:     "mpm",                  // Message Processing Module [recv]
	Port{46, 6}:      "mpm-snd",              // MPM [default send]
	Port{46, 17}:     "mpm-snd",              // MPM [default send]
	Port{47, 6}:      "ni-ftp",               // NI FTP
	Port{47, 17}:     "ni-ftp",               // NI FTP
	Port{48, 6}:      "auditd",               // Digital Audit Daemon
	Port{48, 17}:     "auditd",               // Digital Audit Daemon
	Port{49, 6}:      "tacacs",               // Login Host Protocol (TACACS)
	Port{49, 17}:     "tacacs",               // Login Host Protocol (TACACS)
	Port{50, 6}:      "re-mail-ck",           // Remote Mail Checking Protocol
	Port{50, 17}:     "re-mail-ck",           // Remote Mail Checking Protocol
	Port{51, 6}:      "la-maint",             // IMP Logical Address Maintenance
	Port{51, 17}:     "la-maint",             // IMP Logical Address Maintenance
	Port{52, 6}:      "xns-time",             // XNS Time Protocol
	Port{52, 17}:     "xns-time",             // XNS Time Protocol
	Port{53, 6}:      "domain",               // Domain Name Server
	Port{53, 17}:     "domain",               // Domain Name Server
	Port{54, 6}:      "xns-ch",               // XNS Clearinghouse
	Port{54, 17}:     "xns-ch",               // XNS Clearinghouse
	Port{55, 6}:      "isi-gl",               // ISI Graphics Language
	Port{55, 17}:     "isi-gl",               // ISI Graphics Language
	Port{56, 6}:      "xns-auth",             // XNS Authentication
	Port{56, 17}:     "xns-auth",             // XNS Authentication
	Port{57, 6}:      "priv-term",            // any private terminal access
	Port{57, 17}:     "priv-term",            // any private terminal access
	Port{58, 6}:      "xns-mail",             // XNS Mail
	Port{58, 17}:     "xns-mail",             // XNS Mail
	Port{59, 6}:      "priv-file",            // any private file service
	Port{59, 17}:     "priv-file",            // any private file service
	Port{61, 6}:      "ni-mail",              // NI MAIL
	Port{61, 17}:     "ni-mail",              // NI MAIL
	Port{62, 6}:      "acas",                 // ACA Services
	Port{62, 17}:     "acas",                 // ACA Services
	Port{63, 6}:      "via-ftp",              // whoispp | VIA Systems - FTP & whois++ | whois++
	Port{63, 17}:     "via-ftp",              // VIA Systems - FTP & whois++
	Port{64, 6}:      "covia",                // Communications Integrator (CI)
	Port{64, 17}:     "covia",                // Communications Integrator (CI)
	Port{65, 6}:      "tacacs-ds",            // TACACS-Database Service
	Port{65, 17}:     "tacacs-ds",            // TACACS-Database Service
	Port{66, 6}:      "sqlnet",               // sql*net | sql-net | Oracle SQL*NET
	Port{66, 17}:     "sqlnet",               // Oracle SQL*NET
	Port{67, 6}:      "dhcps",                // bootps | DHCP Bootstrap Protocol Server | Bootstrap Protocol Server
	Port{67, 17}:     "dhcps",                // DHCP Bootstrap Protocol Server
	Port{68, 6}:      "dhcpc",                // bootpc | DHCP Bootstrap Protocol Client | Bootstrap Protocol Client
	Port{68, 17}:     "dhcpc",                // DHCP Bootstrap Protocol Client
	Port{69, 6}:      "tftp",                 // Trivial File Transfer
	Port{69, 17}:     "tftp",                 // Trivial File Transfer
	Port{70, 6}:      "gopher",               // Missing description for gopher
	Port{70, 17}:     "gopher",               // Missing description for gopher
	Port{71, 6}:      "netrjs-1",             // Remote Job Service
	Port{71, 17}:     "netrjs-1",             // Remote Job Service
	Port{72, 6}:      "netrjs-2",             // Remote Job Service
	Port{72, 17}:     "netrjs-2",             // Remote Job Service
	Port{73, 6}:      "netrjs-3",             // Remote Job Service
	Port{73, 17}:     "netrjs-3",             // Remote Job Service
	Port{74, 6}:      "netrjs-4",             // Remote Job Service
	Port{74, 17}:     "netrjs-4",             // Remote Job Service
	Port{75, 6}:      "priv-dial",            // any private dial out service
	Port{75, 17}:     "priv-dial",            // any private dial out service
	Port{76, 6}:      "deos",                 // Distributed External Object Store
	Port{76, 17}:     "deos",                 // Distributed External Object Store
	Port{77, 6}:      "priv-rje",             // any private RJE service, netrjs
	Port{77, 17}:     "priv-rje",             // any private RJE service, netjrs
	Port{78, 6}:      "vettcp",               // Missing description for vettcp
	Port{78, 17}:     "vettcp",               // Missing description for vettcp
	Port{79, 6}:      "finger",               // Missing description for finger
	Port{79, 17}:     "finger",               // Missing description for finger
	Port{80, 132}:    "http",                 // www-http | www | World Wide Web HTTP
	Port{80, 6}:      "http",                 // World Wide Web HTTP
	Port{80, 17}:     "http",                 // World Wide Web HTTP
	Port{81, 6}:      "hosts2-ns",            // HOSTS2 Name Server
	Port{81, 17}:     "hosts2-ns",            // HOSTS2 Name Server
	Port{82, 6}:      "xfer",                 // XFER Utility
	Port{82, 17}:     "xfer",                 // XFER Utility
	Port{83, 6}:      "mit-ml-dev",           // MIT ML Device
	Port{83, 17}:     "mit-ml-dev",           // MIT ML Device
	Port{84, 6}:      "ctf",                  // Common Trace Facility
	Port{84, 17}:     "ctf",                  // Common Trace Facility
	Port{85, 6}:      "mit-ml-dev",           // MIT ML Device
	Port{85, 17}:     "mit-ml-dev",           // MIT ML Device
	Port{86, 6}:      "mfcobol",              // Micro Focus Cobol
	Port{86, 17}:     "mfcobol",              // Micro Focus Cobol
	Port{87, 6}:      "priv-term-l",          // any private terminal link, ttylink
	Port{88, 6}:      "kerberos-sec",         // kerberos | Kerberos (v5) | Kerberos
	Port{88, 17}:     "kerberos-sec",         // Kerberos (v5)
	Port{89, 6}:      "su-mit-tg",            // SU MIT Telnet Gateway
	Port{89, 17}:     "su-mit-tg",            // SU MIT Telnet Gateway
	Port{90, 6}:      "dnsix",                // DNSIX Securit Attribute Token Map
	Port{90, 17}:     "dnsix",                // DNSIX Securit Attribute Token Map
	Port{91, 6}:      "mit-dov",              // MIT Dover Spooler
	Port{91, 17}:     "mit-dov",              // MIT Dover Spooler
	Port{92, 6}:      "npp",                  // Network Printing Protocol
	Port{92, 17}:     "npp",                  // Network Printing Protocol
	Port{93, 6}:      "dcp",                  // Device Control Protocol
	Port{93, 17}:     "dcp",                  // Device Control Protocol
	Port{94, 6}:      "objcall",              // Tivoli Object Dispatcher
	Port{94, 17}:     "objcall",              // Tivoli Object Dispatcher
	Port{95, 6}:      "supdup",               // BSD supdupd(8)
	Port{95, 17}:     "supdup",               // Missing description for supdup
	Port{96, 6}:      "dixie",                // DIXIE Protocol Specification
	Port{96, 17}:     "dixie",                // DIXIE Protocol Specification
	Port{97, 6}:      "swift-rvf",            // Swift Remote Virtural File Protocol
	Port{97, 17}:     "swift-rvf",            // Swift Remote Virtural File Protocol
	Port{98, 6}:      "linuxconf",            // tacnews | TAC News
	Port{98, 17}:     "tacnews",              // TAC News
	Port{99, 6}:      "metagram",             // Metagram Relay
	Port{99, 17}:     "metagram",             // Metagram Relay
	Port{100, 6}:     "newacct",              // [unauthorized use]
	Port{101, 6}:     "hostname",             // hostnames NIC Host Name Server | NIC Host Name Server
	Port{101, 17}:    "hostname",             // hostnames NIC Host Name Server
	Port{102, 6}:     "iso-tsap",             // tsap ISO-TSAP Class 0 | ISO-TSAP Class 0
	Port{102, 17}:    "iso-tsap",             // tsap ISO-TSAP Class 0
	Port{103, 6}:     "gppitnp",              // Genesis Point-to-Point Trans Net, or x400 ISO Email | Genesis Point-to-Point Trans Net
	Port{103, 17}:    "gppitnp",              // Genesis Point-to-Point Trans Net
	Port{104, 6}:     "acr-nema",             // ACR-NEMA Digital Imag. & Comm. 300
	Port{104, 17}:    "acr-nema",             // ACR-NEMA Digital Imag. & Comm. 300
	Port{105, 6}:     "csnet-ns",             // cso | Mailbox Name Nameserver | CCSO name server protocol
	Port{105, 17}:    "csnet-ns",             // Mailbox Name Nameserver
	Port{106, 6}:     "pop3pw",               // 3com-tsmux | Eudora compatible PW changer | 3COM-TSMUX
	Port{106, 17}:    "3com-tsmux",           // Missing description for 3com-tsmux
	Port{107, 6}:     "rtelnet",              // Remote Telnet | Remote Telnet Service
	Port{107, 17}:    "rtelnet",              // Remote Telnet Service
	Port{108, 6}:     "snagas",               // SNA Gateway Access Server
	Port{108, 17}:    "snagas",               // SNA Gateway Access Server
	Port{109, 6}:     "pop2",                 // PostOffice V.2 | Post Office Protocol - Version 2
	Port{109, 17}:    "pop2",                 // PostOffice V.2
	Port{110, 6}:     "pop3",                 // PostOffice V.3 | Post Office Protocol - Version 3
	Port{110, 17}:    "pop3",                 // PostOffice V.3
	Port{111, 6}:     "rpcbind",              // sunrpc | portmapper, rpcbind | SUN Remote Procedure Call
	Port{111, 17}:    "rpcbind",              // portmapper, rpcbind
	Port{112, 6}:     "mcidas",               // McIDAS Data Transmission Protocol
	Port{112, 17}:    "mcidas",               // McIDAS Data Transmission Protocol
	Port{113, 6}:     "ident",                // auth | ident, tap, Authentication Service | Authentication Service
	Port{113, 17}:    "auth",                 // ident, tap, Authentication Service
	Port{114, 6}:     "audionews",            // Audio News Multicast
	Port{114, 17}:    "audionews",            // Audio News Multicast
	Port{115, 6}:     "sftp",                 // Simple File Transfer Protocol
	Port{115, 17}:    "sftp",                 // Simple File Transfer Protocol
	Port{116, 6}:     "ansanotify",           // ANSA REX Notify
	Port{116, 17}:    "ansanotify",           // ANSA REX Notify
	Port{117, 6}:     "uucp-path",            // UUCP Path Service
	Port{117, 17}:    "uucp-path",            // UUCP Path Service
	Port{118, 6}:     "sqlserv",              // SQL Services
	Port{118, 17}:    "sqlserv",              // SQL Services
	Port{119, 6}:     "nntp",                 // Network News Transfer Protocol
	Port{119, 17}:    "nntp",                 // Network News Transfer Protocol
	Port{120, 6}:     "cfdptkt",              // Missing description for cfdptkt
	Port{120, 17}:    "cfdptkt",              // Missing description for cfdptkt
	Port{121, 6}:     "erpc",                 // Encore Expedited Remote Pro.Call
	Port{121, 17}:    "erpc",                 // Encore Expedited Remote Pro.Call
	Port{122, 6}:     "smakynet",             // Missing description for smakynet
	Port{122, 17}:    "smakynet",             // Missing description for smakynet
	Port{123, 6}:     "ntp",                  // Network Time Protocol
	Port{123, 17}:    "ntp",                  // Network Time Protocol
	Port{124, 6}:     "ansatrader",           // ANSA REX Trader
	Port{124, 17}:    "ansatrader",           // ANSA REX Trader
	Port{125, 6}:     "locus-map",            // Locus PC-Interface Net Map Ser
	Port{125, 17}:    "locus-map",            // Locus PC-Interface Net Map Ser
	Port{126, 6}:     "unitary",              // nxedit | Unisys Unitary Login | NXEdit
	Port{126, 17}:    "unitary",              // Unisys Unitary Login
	Port{127, 6}:     "locus-con",            // Locus PC-Interface Conn Server
	Port{127, 17}:    "locus-con",            // Locus PC-Interface Conn Server
	Port{128, 6}:     "gss-xlicen",           // GSS X License Verification
	Port{128, 17}:    "gss-xlicen",           // GSS X License Verification
	Port{129, 6}:     "pwdgen",               // Password Generator Protocol
	Port{129, 17}:    "pwdgen",               // Password Generator Protocol
	Port{130, 6}:     "cisco-fna",            // cisco FNATIVE
	Port{130, 17}:    "cisco-fna",            // cisco FNATIVE
	Port{131, 6}:     "cisco-tna",            // cisco TNATIVE
	Port{131, 17}:    "cisco-tna",            // cisco TNATIVE
	Port{132, 6}:     "cisco-sys",            // cisco SYSMAINT
	Port{132, 17}:    "cisco-sys",            // cisco SYSMAINT
	Port{133, 6}:     "statsrv",              // Statistics Service
	Port{133, 17}:    "statsrv",              // Statistics Service
	Port{134, 6}:     "ingres-net",           // INGRES-NET Service
	Port{134, 17}:    "ingres-net",           // INGRES-NET Service
	Port{135, 6}:     "msrpc",                // epmap | Microsoft RPC services | DCE endpoint resolution
	Port{135, 17}:    "msrpc",                // Microsoft RPC services
	Port{136, 6}:     "profile",              // PROFILE Naming System
	Port{136, 17}:    "profile",              // PROFILE Naming System
	Port{137, 6}:     "netbios-ns",           // NETBIOS Name Service
	Port{137, 17}:    "netbios-ns",           // NETBIOS Name Service
	Port{138, 6}:     "netbios-dgm",          // NETBIOS Datagram Service
	Port{138, 17}:    "netbios-dgm",          // NETBIOS Datagram Service
	Port{139, 6}:     "netbios-ssn",          // NETBIOS Session Service
	Port{139, 17}:    "netbios-ssn",          // NETBIOS Session Service
	Port{140, 6}:     "emfis-data",           // EMFIS Data Service
	Port{140, 17}:    "emfis-data",           // EMFIS Data Service
	Port{141, 6}:     "emfis-cntl",           // EMFIS Control Service
	Port{141, 17}:    "emfis-cntl",           // EMFIS Control Service
	Port{142, 6}:     "bl-idm",               // Britton-Lee IDM
	Port{142, 17}:    "bl-idm",               // Britton-Lee IDM
	Port{143, 6}:     "imap",                 // Interim Mail Access Protocol v2 | Internet Message Access Protocol
	Port{143, 17}:    "imap",                 // Interim Mail Access Protocol v2
	Port{144, 6}:     "news",                 // uma | NewS window system | Universal Management Architecture
	Port{144, 17}:    "news",                 // NewS window system
	Port{145, 6}:     "uaac",                 // UAAC Protocol
	Port{145, 17}:    "uaac",                 // UAAC Protocol
	Port{146, 6}:     "iso-tp0",              // ISO-IP0
	Port{146, 17}:    "iso-tp0",              // Missing description for iso-tp0
	Port{147, 6}:     "iso-ip",               // Missing description for iso-ip
	Port{147, 17}:    "iso-ip",               // Missing description for iso-ip
	Port{148, 6}:     "cronus",               // jargon | CRONUS-SUPPORT | Jargon
	Port{148, 17}:    "cronus",               // CRONUS-SUPPORT
	Port{149, 6}:     "aed-512",              // AED 512 Emulation Service
	Port{149, 17}:    "aed-512",              // AED 512 Emulation Service
	Port{150, 6}:     "sql-net",              // Missing description for sql-net
	Port{150, 17}:    "sql-net",              // Missing description for sql-net
	Port{151, 6}:     "hems",                 // Missing description for hems
	Port{151, 17}:    "hems",                 // Missing description for hems
	Port{152, 6}:     "bftp",                 // Background File Transfer Program
	Port{152, 17}:    "bftp",                 // Background File Transfer Program
	Port{153, 6}:     "sgmp",                 // Missing description for sgmp
	Port{153, 17}:    "sgmp",                 // Missing description for sgmp
	Port{154, 6}:     "netsc-prod",           // NETSC
	Port{154, 17}:    "netsc-prod",           // Missing description for netsc-prod
	Port{155, 6}:     "netsc-dev",            // NETSC
	Port{155, 17}:    "netsc-dev",            // Missing description for netsc-dev
	Port{156, 6}:     "sqlsrv",               // SQL Service
	Port{156, 17}:    "sqlsrv",               // SQL Service
	Port{157, 6}:     "knet-cmp",             // KNET VM Command Message Protocol
	Port{157, 17}:    "knet-cmp",             // KNET VM Command Message Protocol
	Port{158, 6}:     "pcmail-srv",           // PCMail Server
	Port{158, 17}:    "pcmail-srv",           // PCMail Server
	Port{159, 6}:     "nss-routing",          // Missing description for nss-routing
	Port{159, 17}:    "nss-routing",          // Missing description for nss-routing
	Port{160, 6}:     "sgmp-traps",           // Missing description for sgmp-traps
	Port{160, 17}:    "sgmp-traps",           // Missing description for sgmp-traps
	Port{161, 6}:     "snmp",                 // Missing description for snmp
	Port{161, 17}:    "snmp",                 // Simple Net Mgmt Proto
	Port{162, 6}:     "snmptrap",             // snmp-trap
	Port{162, 17}:    "snmptrap",             // snmp-trap
	Port{163, 6}:     "cmip-man",             // CMIP TCP Manager
	Port{163, 17}:    "cmip-man",             // CMIP TCP Manager
	Port{164, 6}:     "cmip-agent",           // CMIP TCP Agent
	Port{164, 17}:    "smip-agent",           // CMIP TCP Agent
	Port{165, 6}:     "xns-courier",          // Xerox
	Port{165, 17}:    "xns-courier",          // Xerox
	Port{166, 6}:     "s-net",                // Sirius Systems
	Port{166, 17}:    "s-net",                // Sirius Systems
	Port{167, 6}:     "namp",                 // Missing description for namp
	Port{167, 17}:    "namp",                 // Missing description for namp
	Port{168, 6}:     "rsvd",                 // Missing description for rsvd
	Port{168, 17}:    "rsvd",                 // Missing description for rsvd
	Port{169, 6}:     "send",                 // Missing description for send
	Port{169, 17}:    "send",                 // Missing description for send
	Port{170, 6}:     "print-srv",            // Network PostScript
	Port{170, 17}:    "print-srv",            // Network PostScript
	Port{171, 6}:     "multiplex",            // Network Innovations Multiplex
	Port{171, 17}:    "multiplex",            // Network Innovations Multiplex
	Port{172, 6}:     "cl-1",                 // cl 1 | Network Innovations CL 1
	Port{172, 17}:    "cl-1",                 // Network Innovations CL 1
	Port{173, 6}:     "xyplex-mux",           // Xyplex
	Port{173, 17}:    "xyplex-mux",           // Missing description for xyplex-mux
	Port{174, 6}:     "mailq",                // Missing description for mailq
	Port{174, 17}:    "mailq",                // Missing description for mailq
	Port{175, 6}:     "vmnet",                // Missing description for vmnet
	Port{175, 17}:    "vmnet",                // Missing description for vmnet
	Port{176, 6}:     "genrad-mux",           // Missing description for genrad-mux
	Port{176, 17}:    "genrad-mux",           // Missing description for genrad-mux
	Port{177, 6}:     "xdmcp",                // X Display Mgr. Control Proto | X Display Manager Control Protocol
	Port{177, 17}:    "xdmcp",                // X Display Manager Control Protocol
	Port{178, 6}:     "nextstep",             // NextStep Window Server
	Port{178, 17}:    "nextstep",             // NextStep Window Server
	Port{179, 132}:   "bgp",                  // Border Gateway Protocol
	Port{179, 6}:     "bgp",                  // Border Gateway Protocol
	Port{179, 17}:    "bgp",                  // Border Gateway Protocol
	Port{180, 6}:     "ris",                  // Intergraph
	Port{180, 17}:    "ris",                  // Intergraph
	Port{181, 6}:     "unify",                // Missing description for unify
	Port{181, 17}:    "unify",                // Missing description for unify
	Port{182, 6}:     "audit",                // Unisys Audit SITP
	Port{182, 17}:    "audit",                // Unisys Audit SITP
	Port{183, 6}:     "ocbinder",             // Missing description for ocbinder
	Port{183, 17}:    "ocbinder",             // Missing description for ocbinder
	Port{184, 6}:     "ocserver",             // Missing description for ocserver
	Port{184, 17}:    "ocserver",             // Missing description for ocserver
	Port{185, 6}:     "remote-kis",           // Missing description for remote-kis
	Port{185, 17}:    "remote-kis",           // Missing description for remote-kis
	Port{186, 6}:     "kis",                  // KIS Protocol
	Port{186, 17}:    "kis",                  // KIS Protocol
	Port{187, 6}:     "aci",                  // Application Communication Interface
	Port{187, 17}:    "aci",                  // Application Communication Interface
	Port{188, 6}:     "mumps",                // Plus Five's MUMPS
	Port{188, 17}:    "mumps",                // Plus Five's MUMPS
	Port{189, 6}:     "qft",                  // Queued File Transport
	Port{189, 17}:    "qft",                  // Queued File Transport
	Port{190, 6}:     "gacp",                 // Gateway Access Control Protocol
	Port{190, 17}:    "cacp",                 // Gateway Access Control Protocol
	Port{191, 6}:     "prospero",             // Prospero Directory Service
	Port{191, 17}:    "prospero",             // Prospero Directory Service
	Port{192, 6}:     "osu-nms",              // OSU Network Monitoring System
	Port{192, 17}:    "osu-nms",              // OSU Network Monitoring System
	Port{193, 6}:     "srmp",                 // Spider Remote Monitoring Protocol
	Port{193, 17}:    "srmp",                 // Spider Remote Monitoring Protocol
	Port{194, 6}:     "irc",                  // Internet Relay Chat | Internet Relay Chat Protocol
	Port{194, 17}:    "irc",                  // Internet Relay Chat Protocol
	Port{195, 6}:     "dn6-nlm-aud",          // DNSIX Network Level Module Audit
	Port{195, 17}:    "dn6-nlm-aud",          // DNSIX Network Level Module Audit
	Port{196, 6}:     "dn6-smm-red",          // DNSIX Session Mgt Module Audit Redir
	Port{196, 17}:    "dn6-smm-red",          // DNSIX Session Mgt Module Audit Redir
	Port{197, 6}:     "dls",                  // Directory Location Service
	Port{197, 17}:    "dls",                  // Directory Location Service
	Port{198, 6}:     "dls-mon",              // Directory Location Service Monitor
	Port{198, 17}:    "dls-mon",              // Directory Location Service Monitor
	Port{199, 6}:     "smux",                 // SNMP Unix Multiplexer
	Port{199, 17}:    "smux",                 // Missing description for smux
	Port{200, 6}:     "src",                  // IBM System Resource Controller
	Port{200, 17}:    "src",                  // IBM System Resource Controller
	Port{201, 6}:     "at-rtmp",              // AppleTalk Routing Maintenance
	Port{201, 17}:    "at-rtmp",              // AppleTalk Routing Maintenance
	Port{202, 6}:     "at-nbp",               // AppleTalk Name Binding
	Port{202, 17}:    "at-nbp",               // AppleTalk Name Binding
	Port{203, 6}:     "at-3",                 // AppleTalk Unused
	Port{203, 17}:    "at-3",                 // AppleTalk Unused
	Port{204, 6}:     "at-echo",              // AppleTalk Echo
	Port{204, 17}:    "at-echo",              // AppleTalk Echo
	Port{205, 6}:     "at-5",                 // AppleTalk Unused
	Port{205, 17}:    "at-5",                 // AppleTalk Unused
	Port{206, 6}:     "at-zis",               // AppleTalk Zone Information
	Port{206, 17}:    "at-zis",               // AppleTalk Zone Information
	Port{207, 6}:     "at-7",                 // AppleTalk Unused
	Port{207, 17}:    "at-7",                 // AppleTalk Unused
	Port{208, 6}:     "at-8",                 // AppleTalk Unused
	Port{208, 17}:    "at-8",                 // AppleTalk Unused
	Port{209, 6}:     "tam",                  // qmtp | Trivial Authenticated Mail Protocol | The Quick Mail Transfer Protocol
	Port{209, 17}:    "tam",                  // Trivial Authenticated Mail Protocol
	Port{210, 6}:     "z39.50",               // z39-50 | wais, ANSI Z39.50 | ANSI Z39.50
	Port{210, 17}:    "z39.50",               // wais, ANSI Z39.50
	Port{211, 6}:     "914c-g",               // 914c g | Texas Instruments 914C G Terminal
	Port{211, 17}:    "914c-g",               // Texas Instruments 914C G Terminal
	Port{212, 6}:     "anet",                 // ATEXSSTR
	Port{212, 17}:    "anet",                 // ATEXSSTR
	Port{213, 6}:     "ipx",                  // Missing description for ipx
	Port{213, 17}:    "ipx",                  // Missing description for ipx
	Port{214, 6}:     "vmpwscs",              // VM PWSCS
	Port{214, 17}:    "vmpwscs",              // Missing description for vmpwscs
	Port{215, 6}:     "softpc",               // Insignia Solutions
	Port{215, 17}:    "softpc",               // Insignia Solutions
	Port{216, 6}:     "atls",                 // CAIlic | Access Technology License Server | Computer Associates Int'l License Server
	Port{216, 17}:    "atls",                 // Access Technology License Server
	Port{217, 6}:     "dbase",                // dBASE Unix
	Port{217, 17}:    "dbase",                // dBASE Unix
	Port{218, 6}:     "mpp",                  // Netix Message Posting Protocol
	Port{218, 17}:    "mpp",                  // Netix Message Posting Protocol
	Port{219, 6}:     "uarps",                // Unisys ARPs
	Port{219, 17}:    "uarps",                // Unisys ARPs
	Port{220, 6}:     "imap3",                // Interactive Mail Access Protocol v3
	Port{220, 17}:    "imap3",                // Interactive Mail Access Protocol v3
	Port{221, 6}:     "fln-spx",              // Berkeley rlogind with SPX auth
	Port{221, 17}:    "fln-spx",              // Berkeley rlogind with SPX auth
	Port{222, 6}:     "rsh-spx",              // Berkeley rshd with SPX auth
	Port{222, 17}:    "rsh-spx",              // Berkeley rshd with SPX auth
	Port{223, 6}:     "cdc",                  // Certificate Distribution Center
	Port{223, 17}:    "cdc",                  // Certificate Distribution Center
	Port{224, 6}:     "masqdialer",           // Missing description for masqdialer
	Port{224, 17}:    "masqdialer",           // Missing description for masqdialer
	Port{242, 6}:     "direct",               // Missing description for direct
	Port{242, 17}:    "direct",               // Missing description for direct
	Port{243, 6}:     "sur-meas",             // Survey Measurement
	Port{243, 17}:    "sur-meas",             // Survey Measurement
	Port{244, 6}:     "dayna",                // inbusiness
	Port{244, 17}:    "dayna",                // Missing description for dayna
	Port{245, 6}:     "link",                 // Missing description for link
	Port{245, 17}:    "link",                 // Missing description for link
	Port{246, 6}:     "dsp3270",              // Display Systems Protocol
	Port{246, 17}:    "dsp3270",              // Display Systems Protocol
	Port{247, 6}:     "subntbcst_tftp",       // subntbcst-tftp
	Port{247, 17}:    "subntbcst_tftp",       // Missing description for subntbcst_tftp
	Port{248, 6}:     "bhfhs",                // Missing description for bhfhs
	Port{248, 17}:    "bhfhs",                // Missing description for bhfhs
	Port{256, 6}:     "fw1-secureremote",     // rap | also "rap" | RAP
	Port{256, 17}:    "rap",                  // Missing description for rap
	Port{257, 6}:     "fw1-mc-fwmodule",      // set | FW1 management console for communication w modules and also secure electronic transaction (set) port | Secure Electronic Transaction
	Port{257, 17}:    "set",                  // secure electronic transaction
	Port{258, 6}:     "fw1-mc-gui",           // also yak winsock personal chat
	Port{258, 17}:    "yak-chat",             // yak winsock personal chat
	Port{259, 6}:     "esro-gen",             // efficient short remote operations | Efficient Short Remote Operations
	Port{259, 17}:    "firewall1-rdp",        // Firewall 1 proprietary RDP protocol http:  www.inside-security.de fw1_rdp_poc.html
	Port{260, 6}:     "openport",             // Missing description for openport
	Port{260, 17}:    "openport",             // Missing description for openport
	Port{261, 6}:     "nsiiops",              // iiop name service over tls ssl | IIOP Name Service over TLS SSL
	Port{261, 17}:    "nsiiops",              // iiop name service over tls ssl
	Port{262, 6}:     "arcisdms",             // Missing description for arcisdms
	Port{262, 17}:    "arcisdms",             // Missing description for arcisdms
	Port{263, 6}:     "hdap",                 // Missing description for hdap
	Port{263, 17}:    "hdap",                 // Missing description for hdap
	Port{264, 6}:     "bgmp",                 // Missing description for bgmp
	Port{264, 17}:    "fw1-or-bgmp",          // FW1 secureremote alternate
	Port{265, 6}:     "maybe-fw1",            // x-bone-ctl | X-Bone CTL
	Port{265, 17}:    "x-bone-ctl",           // X-Bone CTL
	Port{266, 6}:     "sst",                  // SCSI on ST
	Port{266, 17}:    "sst",                  // SCSI on ST
	Port{267, 6}:     "td-service",           // Tobit David Service Layer
	Port{267, 17}:    "td-service",           // Tobit David Service Layer
	Port{268, 6}:     "td-replica",           // Tobit David Replica
	Port{268, 17}:    "td-replica",           // Tobit David Replica
	Port{269, 6}:     "manet",                // MANET Protocols
	Port{269, 17}:    "manet",                // MANET Protocols
	Port{270, 6}:     "gist",                 // Q-mode encapsulation for GIST messages
	Port{270, 17}:    "gist",                 // Q-mode encapsulation for GIST messages
	Port{271, 6}:     "pt-tls",               // IETF Network Endpoint Assessment (NEA) Posture Transport Protocol over TLS (PT-TLS)
	Port{280, 6}:     "http-mgmt",            // Missing description for http-mgmt
	Port{280, 17}:    "http-mgmt",            // Missing description for http-mgmt
	Port{281, 6}:     "personal-link",        // Personal Link
	Port{281, 17}:    "personal-link",        // Missing description for personal-link
	Port{282, 6}:     "cableport-ax",         // cable port a x | Cable Port A X
	Port{282, 17}:    "cableport-ax",         // cable port a x
	Port{283, 6}:     "rescap",               // Missing description for rescap
	Port{283, 17}:    "rescap",               // Missing description for rescap
	Port{284, 6}:     "corerjd",              // Missing description for corerjd
	Port{284, 17}:    "corerjd",              // Missing description for corerjd
	Port{286, 6}:     "fxp",                  // FXP Communication
	Port{286, 17}:    "fxp",                  // FXP Communication
	Port{287, 6}:     "k-block",              // Missing description for k-block
	Port{287, 17}:    "k-block",              // K-BLOCK
	Port{308, 6}:     "novastorbakcup",       // novastor backup | Novastor Backup
	Port{308, 17}:    "novastorbakcup",       // novastor backup
	Port{309, 6}:     "entrusttime",          // Missing description for entrusttime
	Port{309, 17}:    "entrusttime",          // Missing description for entrusttime
	Port{310, 6}:     "bhmds",                // Missing description for bhmds
	Port{310, 17}:    "bhmds",                // Missing description for bhmds
	Port{311, 6}:     "asip-webadmin",        // appleshare ip webadmin | AppleShare IP WebAdmin
	Port{311, 17}:    "asip-webadmin",        // appleshare ip webadmin
	Port{312, 6}:     "vslmp",                // Missing description for vslmp
	Port{312, 17}:    "vslmp",                // Missing description for vslmp
	Port{313, 6}:     "magenta-logic",        // Magenta Logic
	Port{313, 17}:    "magenta-logic",        // Missing description for magenta-logic
	Port{314, 6}:     "opalis-robot",         // Opalis Robot
	Port{314, 17}:    "opalis-robot",         // Missing description for opalis-robot
	Port{315, 6}:     "dpsi",                 // Missing description for dpsi
	Port{315, 17}:    "dpsi",                 // Missing description for dpsi
	Port{316, 6}:     "decauth",              // Missing description for decauth
	Port{316, 17}:    "decauth",              // Missing description for decauth
	Port{317, 6}:     "zannet",               // Missing description for zannet
	Port{317, 17}:    "zannet",               // Missing description for zannet
	Port{318, 6}:     "pkix-timestamp",       // PKIX TimeStamp
	Port{318, 17}:    "pkix-timestamp",       // PKIX TimeStamp
	Port{319, 6}:     "ptp-event",            // PTP Event
	Port{319, 17}:    "ptp-event",            // PTP Event
	Port{320, 6}:     "ptp-general",          // PTP General
	Port{320, 17}:    "ptp-general",          // PTP General
	Port{321, 6}:     "pip",                  // Missing description for pip
	Port{321, 17}:    "pip",                  // Missing description for pip
	Port{322, 6}:     "rtsps",                // Missing description for rtsps
	Port{322, 17}:    "rtsps",                // RTSPS
	Port{323, 6}:     "rpki-rtr",             // Resource PKI to Router Protocol
	Port{324, 6}:     "rpki-rtr-tls",         // Resource PKI to Router Protocol over TLS
	Port{333, 6}:     "texar",                // Texar Security Port
	Port{333, 17}:    "texar",                // Texar Security Port
	Port{344, 6}:     "pdap",                 // Prospero Data Access Protocol
	Port{344, 17}:    "pdap",                 // Prospero Data Access Protocol
	Port{345, 6}:     "pawserv",              // Perf Analysis Workbench
	Port{345, 17}:    "pawserv",              // Perf Analysis Workbench
	Port{346, 6}:     "zserv",                // Zebra server
	Port{346, 17}:    "zserv",                // Zebra server
	Port{347, 6}:     "fatserv",              // Fatmen Server
	Port{347, 17}:    "fatserv",              // Fatmen Server
	Port{348, 6}:     "csi-sgwp",             // Cabletron Management Protocol
	Port{348, 17}:    "csi-sgwp",             // Cabletron Management Protocol
	Port{349, 6}:     "mftp",                 // Missing description for mftp
	Port{349, 17}:    "mftp",                 // Missing description for mftp
	Port{350, 6}:     "matip-type-a",         // MATIP Type A
	Port{350, 17}:    "matip-type-a",         // Missing description for matip-type-a
	Port{351, 6}:     "matip-type-b",         // MATIP Type B or bhoetty also safetp | MATIP Type B | bhoetty
	Port{351, 17}:    "matip-type-b",         // MATIP Type B or bhoetty
	Port{352, 6}:     "dtag-ste-sb",          // DTAG, or bhoedap4 | DTAG | bhoedap4
	Port{352, 17}:    "dtag-ste-sb",          // DTAG, or bhoedap4
	Port{353, 6}:     "ndsauth",              // Missing description for ndsauth
	Port{353, 17}:    "ndsauth",              // Missing description for ndsauth
	Port{354, 6}:     "bh611",                // Missing description for bh611
	Port{354, 17}:    "bh611",                // Missing description for bh611
	Port{355, 6}:     "datex-asn",            // Missing description for datex-asn
	Port{355, 17}:    "datex-asn",            // Missing description for datex-asn
	Port{356, 6}:     "cloanto-net-1",        // Cloanto Net 1
	Port{356, 17}:    "cloanto-net-1",        // Missing description for cloanto-net-1
	Port{357, 6}:     "bhevent",              // Missing description for bhevent
	Port{357, 17}:    "bhevent",              // Missing description for bhevent
	Port{358, 6}:     "shrinkwrap",           // Missing description for shrinkwrap
	Port{358, 17}:    "shrinkwrap",           // Missing description for shrinkwrap
	Port{359, 6}:     "tenebris_nts",         // nsrmp | Tenebris Network Trace Service | Network Security Risk Management Protocol
	Port{359, 17}:    "tenebris_nts",         // Tenebris Network Trace Service
	Port{360, 6}:     "scoi2odialog",         // Missing description for scoi2odialog
	Port{360, 17}:    "scoi2odialog",         // Missing description for scoi2odialog
	Port{361, 6}:     "semantix",             // Missing description for semantix
	Port{361, 17}:    "semantix",             // Missing description for semantix
	Port{362, 6}:     "srssend",              // SRS Send
	Port{362, 17}:    "srssend",              // SRS Send
	Port{363, 6}:     "rsvp_tunnel",          // rsvp-tunnel | RSVP Tunnel
	Port{363, 17}:    "rsvp_tunnel",          // Missing description for rsvp_tunnel
	Port{364, 6}:     "aurora-cmgr",          // Aurora CMGR
	Port{364, 17}:    "aurora-cmgr",          // Missing description for aurora-cmgr
	Port{365, 6}:     "dtk",                  // Deception Tool Kit (www.all.net)
	Port{365, 17}:    "dtk",                  // Deception Tool Kit (www.all.net)
	Port{366, 6}:     "odmr",                 // Missing description for odmr
	Port{366, 17}:    "odmr",                 // Missing description for odmr
	Port{367, 6}:     "mortgageware",         // Missing description for mortgageware
	Port{367, 17}:    "mortgageware",         // Missing description for mortgageware
	Port{368, 6}:     "qbikgdp",              // Missing description for qbikgdp
	Port{368, 17}:    "qbikgdp",              // Missing description for qbikgdp
	Port{369, 6}:     "rpc2portmap",          // Missing description for rpc2portmap
	Port{369, 17}:    "rpc2portmap",          // Missing description for rpc2portmap
	Port{370, 6}:     "codaauth2",            // Missing description for codaauth2
	Port{370, 17}:    "codaauth2",            // Missing description for codaauth2
	Port{371, 6}:     "clearcase",            // Missing description for clearcase
	Port{371, 17}:    "clearcase",            // Missing description for clearcase
	Port{372, 6}:     "ulistserv",            // ulistproc | Unix Listserv | ListProcessor
	Port{372, 17}:    "ulistserv",            // Unix Listserv
	Port{373, 6}:     "legent-1",             // Legent Corporation (now Computer Associates Intl.) | Legent Corporation
	Port{373, 17}:    "legent-1",             // Legent Corporation (now Computer Associates Intl.)
	Port{374, 6}:     "legent-2",             // Legent Corporation (now Computer Associates Intl.) | Legent Corporation
	Port{374, 17}:    "legent-2",             // Legent Corporation (now Computer Associates Intl.)
	Port{375, 6}:     "hassle",               // Missing description for hassle
	Port{375, 17}:    "hassle",               // Missing description for hassle
	Port{376, 6}:     "nip",                  // Amiga Envoy Network Inquiry Proto
	Port{376, 17}:    "nip",                  // Amiga Envoy Network Inquiry Proto
	Port{377, 6}:     "tnETOS",               // NEC Corporation
	Port{377, 17}:    "tnETOS",               // NEC Corporation
	Port{378, 6}:     "dsETOS",               // NEC Corporation
	Port{378, 17}:    "dsETOS",               // NEC Corporation
	Port{379, 6}:     "is99c",                // TIA EIA IS-99 modem client
	Port{379, 17}:    "is99c",                // TIA EIA IS-99 modem client
	Port{380, 6}:     "is99s",                // TIA EIA IS-99 modem server
	Port{380, 17}:    "is99s",                // TIA EIA IS-99 modem server
	Port{381, 6}:     "hp-collector",         // hp performance data collector
	Port{381, 17}:    "hp-collector",         // hp performance data collector
	Port{382, 6}:     "hp-managed-node",      // hp performance data managed node
	Port{382, 17}:    "hp-managed-node",      // hp performance data managed node
	Port{383, 6}:     "hp-alarm-mgr",         // hp performance data alarm manager
	Port{383, 17}:    "hp-alarm-mgr",         // hp performance data alarm manager
	Port{384, 6}:     "arns",                 // A Remote Network Server System
	Port{384, 17}:    "arns",                 // A Remote Network Server System
	Port{385, 6}:     "ibm-app",              // IBM Application
	Port{385, 17}:    "ibm-app",              // IBM Application
	Port{386, 6}:     "asa",                  // ASA Message Router Object Def.
	Port{386, 17}:    "asa",                  // ASA Message Router Object Def.
	Port{387, 6}:     "aurp",                 // Appletalk Update-Based Routing Pro.
	Port{387, 17}:    "aurp",                 // Appletalk Update-Based Routing Pro.
	Port{388, 6}:     "unidata-ldm",          // Unidata LDM Version 4 | Unidata LDM
	Port{388, 17}:    "unidata-ldm",          // Unidata LDM Version 4
	Port{389, 6}:     "ldap",                 // Lightweight Directory Access Protocol
	Port{389, 17}:    "ldap",                 // Lightweight Directory Access Protocol
	Port{390, 6}:     "uis",                  // Missing description for uis
	Port{390, 17}:    "uis",                  // Missing description for uis
	Port{391, 6}:     "synotics-relay",       // SynOptics SNMP Relay Port
	Port{391, 17}:    "synotics-relay",       // SynOptics SNMP Relay Port
	Port{392, 6}:     "synotics-broker",      // SynOptics Port Broker Port
	Port{392, 17}:    "synotics-broker",      // SynOptics Port Broker Port
	Port{393, 6}:     "dis",                  // meta5 | Data Interpretation System | Meta5
	Port{393, 17}:    "dis",                  // Data Interpretation System
	Port{394, 6}:     "embl-ndt",             // EMBL Nucleic Data Transfer
	Port{394, 17}:    "embl-ndt",             // EMBL Nucleic Data Transfer
	Port{395, 6}:     "netcp",                // NETscout Control Protocol | NetScout Control Protocol
	Port{395, 17}:    "netcp",                // NETscout Control Protocol
	Port{396, 6}:     "netware-ip",           // Novell Netware over IP
	Port{396, 17}:    "netware-ip",           // Novell Netware over IP
	Port{397, 6}:     "mptn",                 // Multi Protocol Trans. Net.
	Port{397, 17}:    "mptn",                 // Multi Protocol Trans. Net.
	Port{398, 6}:     "kryptolan",            // Missing description for kryptolan
	Port{398, 17}:    "kryptolan",            // Missing description for kryptolan
	Port{399, 6}:     "iso-tsap-c2",          // ISO-TSAP Class 2 | ISO Transport Class 2 Non-Control over TCP | ISO Transport Class 2 Non-Control over UDP
	Port{399, 17}:    "iso-tsap-c2",          // ISO-TSAP Class 2
	Port{400, 6}:     "work-sol",             // osb-sd | Workstation Solutions | Oracle Secure Backup
	Port{400, 17}:    "work-sol",             // Workstation Solutions
	Port{401, 6}:     "ups",                  // Uninterruptible Power Supply
	Port{401, 17}:    "ups",                  // Uninterruptible Power Supply
	Port{402, 6}:     "genie",                // Genie Protocol
	Port{402, 17}:    "genie",                // Genie Protocol
	Port{403, 6}:     "decap",                // Missing description for decap
	Port{403, 17}:    "decap",                // Missing description for decap
	Port{404, 6}:     "nced",                 // Missing description for nced
	Port{404, 17}:    "nced",                 // Missing description for nced
	Port{405, 6}:     "ncld",                 // Missing description for ncld
	Port{405, 17}:    "ncld",                 // Missing description for ncld
	Port{406, 6}:     "imsp",                 // Interactive Mail Support Protocol
	Port{406, 17}:    "imsp",                 // Interactive Mail Support Protocol
	Port{407, 6}:     "timbuktu",             // Missing description for timbuktu
	Port{407, 17}:    "timbuktu",             // Missing description for timbuktu
	Port{408, 6}:     "prm-sm",               // Prospero Resource Manager Sys. Man.
	Port{408, 17}:    "prm-sm",               // Prospero Resource Manager Sys. Man.
	Port{409, 6}:     "prm-nm",               // Prospero Resource Manager Node Man.
	Port{409, 17}:    "prm-nm",               // Prospero Resource Manager Node Man.
	Port{410, 6}:     "decladebug",           // DECLadebug Remote Debug Protocol
	Port{410, 17}:    "decladebug",           // DECLadebug Remote Debug Protocol
	Port{411, 6}:     "rmt",                  // Remote MT Protocol
	Port{411, 17}:    "rmt",                  // Remote MT Protocol
	Port{412, 6}:     "synoptics-trap",       // Trap Convention Port
	Port{412, 17}:    "synoptics-trap",       // Trap Convention Port
	Port{413, 6}:     "smsp",                 // Storage Management Services Protocol
	Port{413, 17}:    "smsp",                 // Missing description for smsp
	Port{414, 6}:     "infoseek",             // Missing description for infoseek
	Port{414, 17}:    "infoseek",             // Missing description for infoseek
	Port{415, 6}:     "bnet",                 // Missing description for bnet
	Port{415, 17}:    "bnet",                 // Missing description for bnet
	Port{416, 6}:     "silverplatter",        // Missing description for silverplatter
	Port{416, 17}:    "silverplatter",        // Missing description for silverplatter
	Port{417, 6}:     "onmux",                // Meeting maker
	Port{417, 17}:    "onmux",                // Meeting maker
	Port{418, 6}:     "hyper-g",              // Missing description for hyper-g
	Port{418, 17}:    "hyper-g",              // Missing description for hyper-g
	Port{419, 6}:     "ariel1",               // Ariel 1
	Port{419, 17}:    "ariel1",               // Missing description for ariel1
	Port{420, 6}:     "smpte",                // Missing description for smpte
	Port{420, 17}:    "smpte",                // Missing description for smpte
	Port{421, 6}:     "ariel2",               // Ariel 2
	Port{421, 17}:    "ariel2",               // Missing description for ariel2
	Port{422, 6}:     "ariel3",               // Ariel 3
	Port{422, 17}:    "ariel3",               // Missing description for ariel3
	Port{423, 6}:     "opc-job-start",        // IBM Operations Planning and Control Start
	Port{423, 17}:    "opc-job-start",        // IBM Operations Planning and Control Start
	Port{424, 6}:     "opc-job-track",        // IBM Operations Planning and Control Track
	Port{424, 17}:    "opc-job-track",        // IBM Operations Planning and Control Track
	Port{425, 6}:     "icad-el",              // ICAD
	Port{425, 17}:    "icad-el",              // Missing description for icad-el
	Port{426, 6}:     "smartsdp",             // Missing description for smartsdp
	Port{426, 17}:    "smartsdp",             // Missing description for smartsdp
	Port{427, 6}:     "svrloc",               // Server Location
	Port{427, 17}:    "svrloc",               // Server Location
	Port{428, 6}:     "ocs_cmu",              // ocs-cmu
	Port{428, 17}:    "ocs_cmu",              // Missing description for ocs_cmu
	Port{429, 6}:     "ocs_amu",              // ocs-amu
	Port{429, 17}:    "ocs_amu",              // Missing description for ocs_amu
	Port{430, 6}:     "utmpsd",               // Missing description for utmpsd
	Port{430, 17}:    "utmpsd",               // Missing description for utmpsd
	Port{431, 6}:     "utmpcd",               // Missing description for utmpcd
	Port{431, 17}:    "utmpcd",               // Missing description for utmpcd
	Port{432, 6}:     "iasd",                 // Missing description for iasd
	Port{432, 17}:    "iasd",                 // Missing description for iasd
	Port{433, 6}:     "nnsp",                 // Usenet, Network News Transfer | NNTP for transit servers (NNSP)
	Port{433, 17}:    "nnsp",                 // Missing description for nnsp
	Port{434, 6}:     "mobileip-agent",       // Missing description for mobileip-agent
	Port{434, 17}:    "mobileip-agent",       // Missing description for mobileip-agent
	Port{435, 6}:     "mobilip-mn",           // Missing description for mobilip-mn
	Port{435, 17}:    "mobilip-mn",           // Missing description for mobilip-mn
	Port{436, 6}:     "dna-cml",              // Missing description for dna-cml
	Port{436, 17}:    "dna-cml",              // Missing description for dna-cml
	Port{437, 6}:     "comscm",               // Missing description for comscm
	Port{437, 17}:    "comscm",               // Missing description for comscm
	Port{438, 6}:     "dsfgw",                // Missing description for dsfgw
	Port{438, 17}:    "dsfgw",                // Missing description for dsfgw
	Port{439, 6}:     "dasp",                 // Missing description for dasp
	Port{439, 17}:    "dasp",                 // Missing description for dasp
	Port{440, 6}:     "sgcp",                 // Missing description for sgcp
	Port{440, 17}:    "sgcp",                 // Missing description for sgcp
	Port{441, 6}:     "decvms-sysmgt",        // Missing description for decvms-sysmgt
	Port{441, 17}:    "decvms-sysmgt",        // Missing description for decvms-sysmgt
	Port{442, 6}:     "cvc_hostd",            // cvc-hostd
	Port{442, 17}:    "cvc_hostd",            // Missing description for cvc_hostd
	Port{443, 132}:   "https",                // http protocol over TLS SSL
	Port{443, 6}:     "https",                // secure http (SSL)
	Port{443, 17}:    "https",                // Missing description for https
	Port{444, 6}:     "snpp",                 // Simple Network Paging Protocol
	Port{444, 17}:    "snpp",                 // Simple Network Paging Protocol
	Port{445, 6}:     "microsoft-ds",         // SMB directly over IP
	Port{445, 17}:    "microsoft-ds",         // Missing description for microsoft-ds
	Port{446, 6}:     "ddm-rdb",              // DDM-Remote Relational Database Access
	Port{446, 17}:    "ddm-rdb",              // Missing description for ddm-rdb
	Port{447, 6}:     "ddm-dfm",              // DDM-Distributed File Management
	Port{447, 17}:    "ddm-dfm",              // Missing description for ddm-dfm
	Port{448, 6}:     "ddm-ssl",              // ddm-byte | DDM-Remote DB Access Using Secure Sockets
	Port{448, 17}:    "ddm-ssl",              // ddm-byte
	Port{449, 6}:     "as-servermap",         // AS Server Mapper
	Port{449, 17}:    "as-servermap",         // AS Server Mapper
	Port{450, 6}:     "tserver",              // Computer Supported Telecomunication Applications
	Port{450, 17}:    "tserver",              // Missing description for tserver
	Port{451, 6}:     "sfs-smp-net",          // Cray Network Semaphore server
	Port{451, 17}:    "sfs-smp-net",          // Cray Network Semaphore server
	Port{452, 6}:     "sfs-config",           // Cray SFS config server
	Port{452, 17}:    "sfs-config",           // Cray SFS config server
	Port{453, 6}:     "creativeserver",       // Missing description for creativeserver
	Port{453, 17}:    "creativeserver",       // Missing description for creativeserver
	Port{454, 6}:     "contentserver",        // Missing description for contentserver
	Port{454, 17}:    "contentserver",        // Missing description for contentserver
	Port{455, 6}:     "creativepartnr",       // Missing description for creativepartnr
	Port{455, 17}:    "creativepartnr",       // Missing description for creativepartnr
	Port{456, 6}:     "macon",                // macon-tcp | macon-udp
	Port{456, 17}:    "macon",                // Missing description for macon
	Port{457, 6}:     "scohelp",              // Missing description for scohelp
	Port{457, 17}:    "scohelp",              // Missing description for scohelp
	Port{458, 6}:     "appleqtc",             // apple quick time
	Port{458, 17}:    "appleqtc",             // apple quick time
	Port{459, 6}:     "ampr-rcmd",            // Missing description for ampr-rcmd
	Port{459, 17}:    "ampr-rcmd",            // Missing description for ampr-rcmd
	Port{460, 6}:     "skronk",               // Missing description for skronk
	Port{460, 17}:    "skronk",               // Missing description for skronk
	Port{461, 6}:     "datasurfsrv",          // DataRampSrv
	Port{461, 17}:    "datasurfsrv",          // Missing description for datasurfsrv
	Port{462, 6}:     "datasurfsrvsec",       // DataRampSrvSec
	Port{462, 17}:    "datasurfsrvsec",       // Missing description for datasurfsrvsec
	Port{463, 6}:     "alpes",                // Missing description for alpes
	Port{463, 17}:    "alpes",                // Missing description for alpes
	Port{464, 6}:     "kpasswd5",             // Kerberos (v5) | kpasswd
	Port{464, 17}:    "kpasswd5",             // Kerberos (v5)
	Port{465, 6}:     "smtps",                // submissions | igmpv3lite | urd | smtp protocol over TLS SSL (was ssmtp) | URL Rendesvous Directory for SSM | IGMP over UDP for SSM | URL Rendezvous Directory for SSM | Message Submission over TLS protocol
	Port{465, 17}:    "smtps",                // smtp protocol over TLS SSL (was ssmtp)
	Port{466, 6}:     "digital-vrc",          // Missing description for digital-vrc
	Port{466, 17}:    "digital-vrc",          // Missing description for digital-vrc
	Port{467, 6}:     "mylex-mapd",           // Missing description for mylex-mapd
	Port{467, 17}:    "mylex-mapd",           // Missing description for mylex-mapd
	Port{468, 6}:     "photuris",             // Photuris Key Management | proturis
	Port{468, 17}:    "photuris",             // Missing description for photuris
	Port{469, 6}:     "rcp",                  // Radio Control Protocol
	Port{469, 17}:    "rcp",                  // Radio Control Protocol
	Port{470, 6}:     "scx-proxy",            // Missing description for scx-proxy
	Port{470, 17}:    "scx-proxy",            // Missing description for scx-proxy
	Port{471, 6}:     "mondex",               // Missing description for mondex
	Port{471, 17}:    "mondex",               // Missing description for mondex
	Port{472, 6}:     "ljk-login",            // Missing description for ljk-login
	Port{472, 17}:    "ljk-login",            // Missing description for ljk-login
	Port{473, 6}:     "hybrid-pop",           // Missing description for hybrid-pop
	Port{473, 17}:    "hybrid-pop",           // Missing description for hybrid-pop
	Port{474, 6}:     "tn-tl-w1",             // tn-tl-w2
	Port{474, 17}:    "tn-tl-w2",             // Missing description for tn-tl-w2
	Port{475, 6}:     "tcpnethaspsrv",        // Missing description for tcpnethaspsrv
	Port{475, 17}:    "tcpnethaspsrv",        // Missing description for tcpnethaspsrv
	Port{476, 6}:     "tn-tl-fd1",            // Missing description for tn-tl-fd1
	Port{476, 17}:    "tn-tl-fd1",            // Missing description for tn-tl-fd1
	Port{477, 6}:     "ss7ns",                // Missing description for ss7ns
	Port{477, 17}:    "ss7ns",                // Missing description for ss7ns
	Port{478, 6}:     "spsc",                 // Missing description for spsc
	Port{478, 17}:    "spsc",                 // Missing description for spsc
	Port{479, 6}:     "iafserver",            // Missing description for iafserver
	Port{479, 17}:    "iafserver",            // Missing description for iafserver
	Port{480, 6}:     "loadsrv",              // iafdbase
	Port{480, 17}:    "iafdbase",             // Missing description for iafdbase
	Port{481, 6}:     "dvs",                  // ph | Ph service
	Port{481, 17}:    "ph",                   // Missing description for ph
	Port{482, 6}:     "bgs-nsi",              // Missing description for bgs-nsi
	Port{482, 17}:    "xlog",                 // Missing description for xlog
	Port{483, 6}:     "ulpnet",               // Missing description for ulpnet
	Port{483, 17}:    "ulpnet",               // Missing description for ulpnet
	Port{484, 6}:     "integra-sme",          // Integra Software Management Environment
	Port{484, 17}:    "integra-sme",          // Integra Software Management Environment
	Port{485, 6}:     "powerburst",           // Air Soft Power Burst
	Port{485, 17}:    "powerburst",           // Air Soft Power Burst
	Port{486, 6}:     "sstats",               // avian
	Port{486, 17}:    "avian",                // Missing description for avian
	Port{487, 6}:     "saft",                 // saft Simple Asynchronous File Transfer
	Port{487, 17}:    "saft",                 // saft Simple Asynchronous File Transfer
	Port{488, 6}:     "gss-http",             // Missing description for gss-http
	Port{488, 17}:    "gss-http",             // Missing description for gss-http
	Port{489, 6}:     "nest-protocol",        // Missing description for nest-protocol
	Port{489, 17}:    "nest-protocol",        // Missing description for nest-protocol
	Port{490, 6}:     "micom-pfs",            // Missing description for micom-pfs
	Port{490, 17}:    "micom-pfs",            // Missing description for micom-pfs
	Port{491, 6}:     "go-login",             // Missing description for go-login
	Port{491, 17}:    "go-login",             // Missing description for go-login
	Port{492, 6}:     "ticf-1",               // Transport Independent Convergence for FNA
	Port{492, 17}:    "ticf-1",               // Transport Independent Convergence for FNA
	Port{493, 6}:     "ticf-2",               // Transport Independent Convergence for FNA
	Port{493, 17}:    "ticf-2",               // Transport Independent Convergence for FNA
	Port{494, 6}:     "pov-ray",              // Missing description for pov-ray
	Port{494, 17}:    "pov-ray",              // Missing description for pov-ray
	Port{495, 6}:     "intecourier",          // Missing description for intecourier
	Port{495, 17}:    "intecourier",          // Missing description for intecourier
	Port{496, 6}:     "pim-rp-disc",          // Missing description for pim-rp-disc
	Port{496, 17}:    "pim-rp-disc",          // Missing description for pim-rp-disc
	Port{497, 6}:     "retrospect",           // Retrospect backup and restore service
	Port{497, 17}:    "retrospect",           // Missing description for retrospect
	Port{498, 6}:     "siam",                 // Missing description for siam
	Port{498, 17}:    "siam",                 // Missing description for siam
	Port{499, 6}:     "iso-ill",              // ISO ILL Protocol
	Port{499, 17}:    "iso-ill",              // ISO ILL Protocol
	Port{500, 6}:     "isakmp",               // Missing description for isakmp
	Port{500, 17}:    "isakmp",               // Missing description for isakmp
	Port{501, 6}:     "stmf",                 // Missing description for stmf
	Port{501, 17}:    "stmf",                 // Missing description for stmf
	Port{502, 6}:     "mbap",                 // Modbus Application Protocol
	Port{502, 17}:    "mbap",                 // Modbus Application Protocol
	Port{503, 6}:     "intrinsa",             // Missing description for intrinsa
	Port{503, 17}:    "intrinsa",             // Missing description for intrinsa
	Port{504, 6}:     "citadel",              // Missing description for citadel
	Port{504, 17}:    "citadel",              // Missing description for citadel
	Port{505, 6}:     "mailbox-lm",           // Missing description for mailbox-lm
	Port{505, 17}:    "mailbox-lm",           // Missing description for mailbox-lm
	Port{506, 6}:     "ohimsrv",              // Missing description for ohimsrv
	Port{506, 17}:    "ohimsrv",              // Missing description for ohimsrv
	Port{507, 6}:     "crs",                  // Missing description for crs
	Port{507, 17}:    "crs",                  // Missing description for crs
	Port{508, 6}:     "xvttp",                // Missing description for xvttp
	Port{508, 17}:    "xvttp",                // Missing description for xvttp
	Port{509, 6}:     "snare",                // Missing description for snare
	Port{509, 17}:    "snare",                // Missing description for snare
	Port{510, 6}:     "fcp",                  // FirstClass Protocol
	Port{510, 17}:    "fcp",                  // FirstClass Protocol
	Port{511, 6}:     "passgo",               // Missing description for passgo
	Port{511, 17}:    "passgo",               // Missing description for passgo
	Port{512, 6}:     "exec",                 // biff | comsat | BSD rexecd(8) | remote process execution; authentication performed using passwords and UNIX login names | used by mail system to notify users of new mail received; currently receives messages only from processes on the same machine
	Port{512, 17}:    "biff",                 // comsat
	Port{513, 6}:     "login",                // who | BSD rlogind(8) | remote login a la telnet; automatic authentication performed based on priviledged port numbers and distributed data bases which identify "authentication domains" | maintains data bases showing who's logged in to machines on a local net and the load average of the machine
	Port{513, 17}:    "who",                  // BSD rwhod(8)
	Port{514, 6}:     "shell",                // syslog | BSD rshd(8) | cmd like exec, but automatic authentication is performed as for login server
	Port{514, 17}:    "syslog",               // BSD syslogd(8)
	Port{515, 6}:     "printer",              // spooler (lpd) | spooler
	Port{515, 17}:    "printer",              // spooler (lpd)
	Port{516, 6}:     "videotex",             // Missing description for videotex
	Port{516, 17}:    "videotex",             // Missing description for videotex
	Port{517, 6}:     "talk",                 // like tenex link, but across | like tenex link, but across machine - unfortunately, doesn't use link protocol (this is actually just a rendezvous port from which a tcp connection is established)
	Port{517, 17}:    "talk",                 // BSD talkd(8)
	Port{518, 6}:     "ntalk",                // (talkd)
	Port{518, 17}:    "ntalk",                // (talkd)
	Port{519, 6}:     "utime",                // unixtime
	Port{519, 17}:    "utime",                // unixtime
	Port{520, 6}:     "efs",                  // router | extended file name server | local routing process (on site); uses variant of Xerox NS routing information protocol - RIP
	Port{520, 17}:    "route",                // router routed -- RIP
	Port{521, 6}:     "ripng",                // Missing description for ripng
	Port{521, 17}:    "ripng",                // Missing description for ripng
	Port{522, 6}:     "ulp",                  // Missing description for ulp
	Port{522, 17}:    "ulp",                  // Missing description for ulp
	Port{523, 6}:     "ibm-db2",              // Missing description for ibm-db2
	Port{523, 17}:    "ibm-db2",              // Missing description for ibm-db2
	Port{524, 6}:     "ncp",                  // Missing description for ncp
	Port{524, 17}:    "ncp",                  // Missing description for ncp
	Port{525, 6}:     "timed",                // timeserver
	Port{525, 17}:    "timed",                // timeserver
	Port{526, 6}:     "tempo",                // newdate
	Port{526, 17}:    "tempo",                // newdate
	Port{527, 6}:     "stx",                  // Stock IXChange
	Port{527, 17}:    "stx",                  // Stock IXChange
	Port{528, 6}:     "custix",               // Customer IXChange
	Port{528, 17}:    "custix",               // Customer IXChange
	Port{529, 6}:     "irc",                  // irc-serv | IRC-SERV
	Port{529, 17}:    "irc",                  // Missing description for irc
	Port{530, 6}:     "courier",              // rpc
	Port{530, 17}:    "courier",              // rpc
	Port{531, 6}:     "conference",           // chat
	Port{531, 17}:    "conference",           // chat
	Port{532, 6}:     "netnews",              // readnews
	Port{532, 17}:    "netnews",              // readnews
	Port{533, 6}:     "netwall",              // for emergency broadcasts
	Port{533, 17}:    "netwall",              // for emergency broadcasts
	Port{534, 6}:     "mm-admin",             // windream | MegaMedia Admin | windream Admin
	Port{534, 17}:    "mm-admin",             // MegaMedia Admin
	Port{535, 6}:     "iiop",                 // Missing description for iiop
	Port{535, 17}:    "iiop",                 // Missing description for iiop
	Port{536, 6}:     "opalis-rdv",           // Missing description for opalis-rdv
	Port{536, 17}:    "opalis-rdv",           // Missing description for opalis-rdv
	Port{537, 6}:     "nmsp",                 // Networked Media Streaming Protocol
	Port{537, 17}:    "nmsp",                 // Networked Media Streaming Protocol
	Port{538, 6}:     "gdomap",               // Missing description for gdomap
	Port{538, 17}:    "gdomap",               // Missing description for gdomap
	Port{539, 6}:     "apertus-ldp",          // Apertus Technologies Load Determination
	Port{539, 17}:    "apertus-ldp",          // Apertus Technologies Load Determination
	Port{540, 6}:     "uucp",                 // uucpd
	Port{540, 17}:    "uucp",                 // uucpd
	Port{541, 6}:     "uucp-rlogin",          // Missing description for uucp-rlogin
	Port{541, 17}:    "uucp-rlogin",          // Missing description for uucp-rlogin
	Port{542, 6}:     "commerce",             // Missing description for commerce
	Port{542, 17}:    "commerce",             // Missing description for commerce
	Port{543, 6}:     "klogin",               // Kerberos (v4 v5)
	Port{543, 17}:    "klogin",               // Kerberos (v4 v5)
	Port{544, 6}:     "kshell",               // krcmd Kerberos (v4 v5) | krcmd
	Port{544, 17}:    "kshell",               // krcmd Kerberos (v4 v5)
	Port{545, 6}:     "ekshell",              // Kerberos encrypted remote shell -kfall | appleqtcsrvr
	Port{545, 17}:    "appleqtcsrvr",         // Missing description for appleqtcsrvr
	Port{546, 6}:     "dhcpv6-client",        // DHCPv6 Client
	Port{546, 17}:    "dhcpv6-client",        // DHCPv6 Client
	Port{547, 6}:     "dhcpv6-server",        // DHCPv6 Server
	Port{547, 17}:    "dhcpv6-server",        // DHCPv6 Server
	Port{548, 6}:     "afp",                  // afpovertcp | AFP over TCP
	Port{548, 17}:    "afp",                  // AFP over UDP
	Port{549, 6}:     "idfp",                 // Missing description for idfp
	Port{549, 17}:    "idfp",                 // Missing description for idfp
	Port{550, 6}:     "new-rwho",             // new-who
	Port{550, 17}:    "new-rwho",             // new-who
	Port{551, 6}:     "cybercash",            // Missing description for cybercash
	Port{551, 17}:    "cybercash",            // Missing description for cybercash
	Port{552, 6}:     "deviceshare",          // devshr-nts
	Port{552, 17}:    "deviceshare",          // Missing description for deviceshare
	Port{553, 6}:     "pirp",                 // Missing description for pirp
	Port{553, 17}:    "pirp",                 // Missing description for pirp
	Port{554, 6}:     "rtsp",                 // Real Time Stream Control Protocol | Real Time Streaming Protocol (RTSP)
	Port{554, 17}:    "rtsp",                 // Real Time Stream Control Protocol
	Port{555, 6}:     "dsf",                  // Missing description for dsf
	Port{555, 17}:    "dsf",                  // Missing description for dsf
	Port{556, 6}:     "remotefs",             // rfs, rfs_server, Brunhoff remote filesystem | rfs server
	Port{556, 17}:    "remotefs",             // rfs, rfs_server, Brunhoff remote filesystem
	Port{557, 6}:     "openvms-sysipc",       // Missing description for openvms-sysipc
	Port{557, 17}:    "openvms-sysipc",       // Missing description for openvms-sysipc
	Port{558, 6}:     "sdnskmp",              // Missing description for sdnskmp
	Port{558, 17}:    "sdnskmp",              // Missing description for sdnskmp
	Port{559, 6}:     "teedtap",              // Missing description for teedtap
	Port{559, 17}:    "teedtap",              // Missing description for teedtap
	Port{560, 6}:     "rmonitor",             // rmonitord
	Port{560, 17}:    "rmonitor",             // rmonitord
	Port{561, 6}:     "monitor",              // Missing description for monitor
	Port{561, 17}:    "monitor",              // Missing description for monitor
	Port{562, 6}:     "chshell",              // chcmd
	Port{562, 17}:    "chshell",              // chcmd
	Port{563, 6}:     "snews",                // nntps | nntp protocol over TLS SSL (was snntp)
	Port{563, 17}:    "snews",                // Missing description for snews
	Port{564, 6}:     "9pfs",                 // plan 9 file service
	Port{564, 17}:    "9pfs",                 // plan 9 file service
	Port{565, 6}:     "whoami",               // Missing description for whoami
	Port{565, 17}:    "whoami",               // Missing description for whoami
	Port{566, 6}:     "streettalk",           // Missing description for streettalk
	Port{566, 17}:    "streettalk",           // Missing description for streettalk
	Port{567, 6}:     "banyan-rpc",           // Missing description for banyan-rpc
	Port{567, 17}:    "banyan-rpc",           // Missing description for banyan-rpc
	Port{568, 6}:     "ms-shuttle",           // Microsoft shuttle | microsoft shuttle
	Port{568, 17}:    "ms-shuttle",           // Microsoft shuttle
	Port{569, 6}:     "ms-rome",              // Microsoft rome | microsoft rome
	Port{569, 17}:    "ms-rome",              // Microsoft rome
	Port{570, 6}:     "meter",                // demon
	Port{570, 17}:    "meter",                // demon
	Port{571, 6}:     "umeter",               // meter | udemon
	Port{571, 17}:    "umeter",               // udemon
	Port{572, 6}:     "sonar",                // Missing description for sonar
	Port{572, 17}:    "sonar",                // Missing description for sonar
	Port{573, 6}:     "banyan-vip",           // Missing description for banyan-vip
	Port{573, 17}:    "banyan-vip",           // Missing description for banyan-vip
	Port{574, 6}:     "ftp-agent",            // FTP Software Agent System
	Port{574, 17}:    "ftp-agent",            // FTP Software Agent System
	Port{575, 6}:     "vemmi",                // Missing description for vemmi
	Port{575, 17}:    "vemmi",                // Missing description for vemmi
	Port{576, 6}:     "ipcd",                 // Missing description for ipcd
	Port{576, 17}:    "ipcd",                 // Missing description for ipcd
	Port{577, 6}:     "vnas",                 // Missing description for vnas
	Port{577, 17}:    "vnas",                 // Missing description for vnas
	Port{578, 6}:     "ipdd",                 // Missing description for ipdd
	Port{578, 17}:    "ipdd",                 // Missing description for ipdd
	Port{579, 6}:     "decbsrv",              // Missing description for decbsrv
	Port{579, 17}:    "decbsrv",              // Missing description for decbsrv
	Port{580, 6}:     "sntp-heartbeat",       // SNTP HEARTBEAT
	Port{580, 17}:    "sntp-heartbeat",       // Missing description for sntp-heartbeat
	Port{581, 6}:     "bdp",                  // Bundle Discovery Protocol
	Port{581, 17}:    "bdp",                  // Bundle Discovery Protocol
	Port{582, 6}:     "scc-security",         // SCC Security
	Port{582, 17}:    "scc-security",         // Missing description for scc-security
	Port{583, 6}:     "philips-vc",           // Philips Video-Conferencing
	Port{583, 17}:    "philips-vc",           // Philips Video-Conferencing
	Port{584, 6}:     "keyserver",            // Key Server
	Port{584, 17}:    "keyserver",            // Missing description for keyserver
	Port{585, 6}:     "imap4-ssl",            // IMAP4+SSL (use of 585 is not recommended,
	Port{585, 17}:    "imap4-ssl",            // use 993 instead)
	Port{586, 6}:     "password-chg",         // Password Change
	Port{586, 17}:    "password-chg",         // Missing description for password-chg
	Port{587, 6}:     "submission",           // Message Submission
	Port{587, 17}:    "submission",           // Missing description for submission
	Port{588, 6}:     "cal",                  // Missing description for cal
	Port{588, 17}:    "cal",                  // Missing description for cal
	Port{589, 6}:     "eyelink",              // Missing description for eyelink
	Port{589, 17}:    "eyelink",              // Missing description for eyelink
	Port{590, 6}:     "tns-cml",              // TNS CML
	Port{590, 17}:    "tns-cml",              // Missing description for tns-cml
	Port{591, 6}:     "http-alt",             // FileMaker, Inc. - HTTP Alternate | FileMaker, Inc. - HTTP Alternate (see Port 80)
	Port{591, 17}:    "http-alt",             // FileMaker, Inc. - HTTP Alternate
	Port{592, 6}:     "eudora-set",           // Eudora Set
	Port{592, 17}:    "eudora-set",           // Missing description for eudora-set
	Port{593, 6}:     "http-rpc-epmap",       // HTTP RPC Ep Map
	Port{593, 17}:    "http-rpc-epmap",       // HTTP RPC Ep Map
	Port{594, 6}:     "tpip",                 // Missing description for tpip
	Port{594, 17}:    "tpip",                 // Missing description for tpip
	Port{595, 6}:     "cab-protocol",         // CAB Protocol
	Port{595, 17}:    "cab-protocol",         // Missing description for cab-protocol
	Port{596, 6}:     "smsd",                 // Missing description for smsd
	Port{596, 17}:    "smsd",                 // Missing description for smsd
	Port{597, 6}:     "ptcnameservice",       // PTC Name Service
	Port{597, 17}:    "ptcnameservice",       // PTC Name Service
	Port{598, 6}:     "sco-websrvrmg3",       // SCO Web Server Manager 3
	Port{598, 17}:    "sco-websrvrmg3",       // SCO Web Server Manager 3
	Port{599, 6}:     "acp",                  // Aeolon Core Protocol
	Port{599, 17}:    "acp",                  // Aeolon Core Protocol
	Port{600, 6}:     "ipcserver",            // Sun IPC server
	Port{600, 17}:    "ipcserver",            // Sun IPC server
	Port{601, 6}:     "syslog-conn",          // Reliable Syslog Service
	Port{601, 17}:    "syslog-conn",          // Reliable Syslog Service
	Port{602, 6}:     "xmlrpc-beep",          // XML-RPC over BEEP
	Port{602, 17}:    "xmlrpc-beep",          // XML-RPC over BEEP
	Port{603, 6}:     "mnotes",               // idxp | CommonTime Mnotes PDA Synchronization | IDXP
	Port{603, 17}:    "idxp",                 // IDXP
	Port{604, 6}:     "tunnel",               // Missing description for tunnel
	Port{604, 17}:    "tunnel",               // TUNNEL
	Port{605, 6}:     "soap-beep",            // SOAP over BEEP
	Port{605, 17}:    "soap-beep",            // SOAP over BEEP
	Port{606, 6}:     "urm",                  // Cray Unified Resource Manager
	Port{606, 17}:    "urm",                  // Cray Unified Resource Manager
	Port{607, 6}:     "nqs",                  // Missing description for nqs
	Port{607, 17}:    "nqs",                  // Missing description for nqs
	Port{608, 6}:     "sift-uft",             // Sender-Initiated Unsolicited File Transfer
	Port{608, 17}:    "sift-uft",             // Sender-Initiated Unsolicited File Transfer
	Port{609, 6}:     "npmp-trap",            // Missing description for npmp-trap
	Port{609, 17}:    "npmp-trap",            // Missing description for npmp-trap
	Port{610, 6}:     "npmp-local",           // Missing description for npmp-local
	Port{610, 17}:    "npmp-local",           // Missing description for npmp-local
	Port{611, 6}:     "npmp-gui",             // Missing description for npmp-gui
	Port{611, 17}:    "npmp-gui",             // Missing description for npmp-gui
	Port{612, 6}:     "hmmp-ind",             // HMMP Indication
	Port{612, 17}:    "hmmp-ind",             // HMMP Indication
	Port{613, 6}:     "hmmp-op",              // HMMP Operation
	Port{613, 17}:    "hmmp-op",              // HMMP Operation
	Port{614, 6}:     "sshell",               // SSLshell
	Port{614, 17}:    "sshell",               // SSLshell
	Port{615, 6}:     "sco-inetmgr",          // Internet Configuration Manager
	Port{615, 17}:    "sco-inetmgr",          // Internet Configuration Manager
	Port{616, 6}:     "sco-sysmgr",           // SCO System Administration Server
	Port{616, 17}:    "sco-sysmgr",           // SCO System Administration Server
	Port{617, 6}:     "sco-dtmgr",            // SCO Desktop Administration Server or Arkeia (www.arkeia.com) backup software | SCO Desktop Administration Server
	Port{617, 17}:    "sco-dtmgr",            // SCO Desktop Administration Server
	Port{618, 6}:     "dei-icda",             // Missing description for dei-icda
	Port{618, 17}:    "dei-icda",             // DEI-ICDA
	Port{619, 6}:     "compaq-evm",           // Compaq EVM
	Port{619, 17}:    "compaq-evm",           // Compaq EVM
	Port{620, 6}:     "sco-websrvrmgr",       // SCO WebServer Manager
	Port{620, 17}:    "sco-websrvrmgr",       // SCO WebServer Manager
	Port{621, 6}:     "escp-ip",              // ESCP
	Port{621, 17}:    "escp-ip",              // ESCP
	Port{622, 6}:     "collaborator",         // Missing description for collaborator
	Port{622, 17}:    "collaborator",         // Collaborator
	Port{623, 6}:     "oob-ws-http",          // asf-rmcp | DMTF out-of-band web services management protocol | ASF Remote Management and Control Protocol
	Port{623, 17}:    "asf-rmcp",             // ASF Remote Management and Control
	Port{624, 6}:     "cryptoadmin",          // Crypto Admin
	Port{624, 17}:    "cryptoadmin",          // Crypto Admin
	Port{625, 6}:     "apple-xsrvr-admin",    // dec_dlm | dec-dlm | Apple Mac Xserver admin | DEC DLM
	Port{625, 17}:    "dec_dlm",              // DEC DLM
	Port{626, 6}:     "apple-imap-admin",     // asia | Apple IMAP mail admin | ASIA
	Port{626, 17}:    "serialnumberd",        // Mac OS X Server serial number (licensing) daemon
	Port{627, 6}:     "passgo-tivoli",        // PassGo Tivoli
	Port{627, 17}:    "passgo-tivoli",        // PassGo Tivoli
	Port{628, 6}:     "qmqp",                 // Qmail Quick Mail Queueing
	Port{628, 17}:    "qmqp",                 // QMQP
	Port{629, 6}:     "3com-amp3",            // 3Com AMP3
	Port{629, 17}:    "3com-amp3",            // 3Com AMP3
	Port{630, 6}:     "rda",                  // Missing description for rda
	Port{630, 17}:    "rda",                  // RDA
	Port{631, 6}:     "ipp",                  // Internet Printing Protocol -- for one implementation see http:  www.cups.org (Common UNIX Printing System) | IPP (Internet Printing Protocol)
	Port{631, 17}:    "ipp",                  // Internet Printing Protocol
	Port{632, 6}:     "bmpp",                 // Missing description for bmpp
	Port{632, 17}:    "bmpp",                 // Missing description for bmpp
	Port{633, 6}:     "servstat",             // Service Status update (Sterling Software)
	Port{633, 17}:    "servstat",             // Service Status update (Sterling Software)
	Port{634, 6}:     "ginad",                // Missing description for ginad
	Port{634, 17}:    "ginad",                // Missing description for ginad
	Port{635, 6}:     "rlzdbase",             // RLZ DBase
	Port{635, 17}:    "mount",                // NFS Mount Service
	Port{636, 6}:     "ldapssl",              // ldaps | LDAP over SSL | ldap protocol over TLS SSL (was sldap)
	Port{636, 17}:    "ldaps",                // ldap protocol over TLS SSL (was sldap)
	Port{637, 6}:     "lanserver",            // Missing description for lanserver
	Port{637, 17}:    "lanserver",            // Missing description for lanserver
	Port{638, 6}:     "mcns-sec",             // Missing description for mcns-sec
	Port{638, 17}:    "mcns-sec",             // Missing description for mcns-sec
	Port{639, 6}:     "msdp",                 // Missing description for msdp
	Port{639, 17}:    "msdp",                 // MSDP
	Port{640, 6}:     "entrust-sps",          // Missing description for entrust-sps
	Port{640, 17}:    "pcnfs",                // PC-NFS DOS Authentication
	Port{641, 6}:     "repcmd",               // Missing description for repcmd
	Port{641, 17}:    "repcmd",               // Missing description for repcmd
	Port{642, 6}:     "esro-emsdp",           // ESRO-EMSDP V1.3
	Port{642, 17}:    "esro-emsdp",           // ESRO-EMSDP V1.3
	Port{643, 6}:     "sanity",               // Missing description for sanity
	Port{643, 17}:    "sanity",               // SANity
	Port{644, 6}:     "dwr",                  // Missing description for dwr
	Port{644, 17}:    "dwr",                  // Missing description for dwr
	Port{645, 6}:     "pssc",                 // Missing description for pssc
	Port{645, 17}:    "pssc",                 // PSSC
	Port{646, 6}:     "ldp",                  // Label Distribution
	Port{646, 17}:    "ldp",                  // Label Distribution
	Port{647, 6}:     "dhcp-failover",        // DHCP Failover
	Port{647, 17}:    "dhcp-failover",        // DHCP Failover
	Port{648, 6}:     "rrp",                  // Registry Registrar Protocol (RRP)
	Port{648, 17}:    "rrp",                  // Registry Registrar Protocol (RRP)
	Port{649, 6}:     "cadview-3d",           // Cadview-3d - streaming 3d models over the internet
	Port{649, 17}:    "cadview-3d",           // Cadview-3d - streaming 3d models over the internet
	Port{650, 6}:     "obex",                 // Missing description for obex
	Port{650, 17}:    "bwnfs",                // BW-NFS DOS Authentication
	Port{651, 6}:     "ieee-mms",             // IEEE MMS
	Port{651, 17}:    "ieee-mms",             // IEEE MMS
	Port{652, 6}:     "hello-port",           // HELLO_PORT
	Port{652, 17}:    "hello-port",           // HELLO_PORT
	Port{653, 6}:     "repscmd",              // RepCmd
	Port{653, 17}:    "repscmd",              // RepCmd
	Port{654, 6}:     "aodv",                 // Missing description for aodv
	Port{654, 17}:    "aodv",                 // AODV
	Port{655, 6}:     "tinc",                 // Missing description for tinc
	Port{655, 17}:    "tinc",                 // TINC
	Port{656, 6}:     "spmp",                 // Missing description for spmp
	Port{656, 17}:    "spmp",                 // SPMP
	Port{657, 6}:     "rmc",                  // Missing description for rmc
	Port{657, 17}:    "rmc",                  // RMC
	Port{658, 6}:     "tenfold",              // Missing description for tenfold
	Port{658, 17}:    "tenfold",              // TenFold
	Port{660, 6}:     "mac-srvr-admin",       // MacOS Server Admin
	Port{660, 17}:    "mac-srvr-admin",       // MacOS Server Admin
	Port{661, 6}:     "hap",                  // Missing description for hap
	Port{661, 17}:    "hap",                  // HAP
	Port{662, 6}:     "pftp",                 // Missing description for pftp
	Port{662, 17}:    "pftp",                 // PFTP
	Port{663, 6}:     "purenoise",            // Missing description for purenoise
	Port{663, 17}:    "purenoise",            // PureNoise
	Port{664, 6}:     "secure-aux-bus",       // asf-secure-rmcp | oob-ws-https | DMTF out-of-band secure web services management protocol | ASF Secure Remote Management and Control Protocol
	Port{664, 17}:    "secure-aux-bus",       // Missing description for secure-aux-bus
	Port{665, 6}:     "sun-dr",               // Sun DR
	Port{665, 17}:    "sun-dr",               // Sun DR
	Port{666, 6}:     "doom",                 // mdqs | Id Software Doom | doom Id Software
	Port{666, 17}:    "doom",                 // doom Id Software
	Port{667, 6}:     "disclose",             // campaign contribution disclosures - SDR Technologies
	Port{667, 17}:    "disclose",             // campaign contribution disclosures - SDR Technologies
	Port{668, 6}:     "mecomm",               // Missing description for mecomm
	Port{668, 17}:    "mecomm",               // MeComm
	Port{669, 6}:     "meregister",           // Missing description for meregister
	Port{669, 17}:    "meregister",           // MeRegister
	Port{670, 6}:     "vacdsm-sws",           // Missing description for vacdsm-sws
	Port{670, 17}:    "vacdsm-sws",           // VACDSM-SWS
	Port{671, 6}:     "vacdsm-app",           // Missing description for vacdsm-app
	Port{671, 17}:    "vacdsm-app",           // VACDSM-APP
	Port{672, 6}:     "vpps-qua",             // Missing description for vpps-qua
	Port{672, 17}:    "vpps-qua",             // VPPS-QUA
	Port{673, 6}:     "cimplex",              // Missing description for cimplex
	Port{673, 17}:    "cimplex",              // CIMPLEX
	Port{674, 6}:     "acap",                 // ACAP server of Communigate (www.stalker.com)
	Port{674, 17}:    "acap",                 // ACAP
	Port{675, 6}:     "dctp",                 // Missing description for dctp
	Port{675, 17}:    "dctp",                 // DCTP
	Port{676, 6}:     "vpps-via",             // VPPS Via
	Port{676, 17}:    "vpps-via",             // VPPS Via
	Port{677, 6}:     "vpp",                  // Virtual Presence Protocol
	Port{677, 17}:    "vpp",                  // Virtual Presence Protocol
	Port{678, 6}:     "ggf-ncp",              // GNU Generation Foundation NCP
	Port{678, 17}:    "ggf-ncp",              // GNU Generation Foundation NCP
	Port{679, 6}:     "mrm",                  // Missing description for mrm
	Port{679, 17}:    "mrm",                  // MRM
	Port{680, 6}:     "entrust-aaas",         // Missing description for entrust-aaas
	Port{680, 17}:    "entrust-aaas",         // Missing description for entrust-aaas
	Port{681, 6}:     "entrust-aams",         // Missing description for entrust-aams
	Port{681, 17}:    "entrust-aams",         // Missing description for entrust-aams
	Port{682, 6}:     "xfr",                  // Missing description for xfr
	Port{682, 17}:    "xfr",                  // XFR
	Port{683, 6}:     "corba-iiop",           // CORBA IIOP
	Port{683, 17}:    "corba-iiop",           // Missing description for corba-iiop
	Port{684, 6}:     "corba-iiop-ssl",       // CORBA IIOP SSL
	Port{684, 17}:    "corba-iiop-ssl",       // CORBA IIOP SSL
	Port{685, 6}:     "mdc-portmapper",       // MDC Port Mapper
	Port{685, 17}:    "mdc-portmapper",       // MDC Port Mapper
	Port{686, 6}:     "hcp-wismar",           // Hardware Control Protocol Wismar
	Port{686, 17}:    "hcp-wismar",           // Hardware Control Protocol Wismar
	Port{687, 6}:     "asipregistry",         // Missing description for asipregistry
	Port{687, 17}:    "asipregistry",         // Missing description for asipregistry
	Port{688, 6}:     "realm-rusd",           // ApplianceWare managment protocol
	Port{688, 17}:    "realm-rusd",           // ApplianceWare managment protocol
	Port{689, 6}:     "nmap",                 // Missing description for nmap
	Port{689, 17}:    "nmap",                 // NMAP
	Port{690, 6}:     "vatp",                 // Velazquez Application Transfer Protocol | Velneo Application Transfer Protocol
	Port{690, 17}:    "vatp",                 // Velazquez Application Transfer Protocol
	Port{691, 6}:     "resvc",                // msexch-routing | The Microsoft Exchange 2000 Server Routing Service | MS Exchange Routing
	Port{691, 17}:    "msexch-routing",       // MS Exchange Routing
	Port{692, 6}:     "hyperwave-isp",        // Missing description for hyperwave-isp
	Port{692, 17}:    "hyperwave-isp",        // Hyperwave-ISP
	Port{693, 6}:     "connendp",             // almanid Connection Endpoint
	Port{693, 17}:    "connendp",             // almanid Connection Endpoint
	Port{694, 6}:     "ha-cluster",           // Missing description for ha-cluster
	Port{694, 17}:    "ha-cluster",           // Missing description for ha-cluster
	Port{695, 6}:     "ieee-mms-ssl",         // Missing description for ieee-mms-ssl
	Port{695, 17}:    "ieee-mms-ssl",         // IEEE-MMS-SSL
	Port{696, 6}:     "rushd",                // Missing description for rushd
	Port{696, 17}:    "rushd",                // RUSHD
	Port{697, 6}:     "uuidgen",              // Missing description for uuidgen
	Port{697, 17}:    "uuidgen",              // UUIDGEN
	Port{698, 6}:     "olsr",                 // Missing description for olsr
	Port{698, 17}:    "olsr",                 // OLSR
	Port{699, 6}:     "accessnetwork",        // Access Network
	Port{699, 17}:    "accessnetwork",        // Access Network
	Port{700, 6}:     "epp",                  // Extensible Provisioning Protocol
	Port{700, 17}:    "epp",                  // Extensible Provisioning Protocol
	Port{701, 6}:     "lmp",                  // Link Management Protocol (LMP)
	Port{701, 17}:    "lmp",                  // Link Management Protocol (LMP)
	Port{702, 6}:     "iris-beep",            // IRIS over BEEP
	Port{702, 17}:    "iris-beep",            // IRIS over BEEP
	Port{704, 6}:     "elcsd",                // errlog copy server daemon
	Port{704, 17}:    "elcsd",                // errlog copy server daemon
	Port{705, 6}:     "agentx",               // Missing description for agentx
	Port{705, 17}:    "agentx",               // AgentX
	Port{706, 6}:     "silc",                 // Secure Internet Live Conferencing -- http:  silcnet.org
	Port{706, 17}:    "silc",                 // SILC
	Port{707, 6}:     "borland-dsj",          // Borland DSJ
	Port{707, 17}:    "borland-dsj",          // Borland DSJ
	Port{709, 6}:     "entrustmanager",       // entrust-kmsh | EntrustManager - NorTel DES auth network see 389 tcp | Entrust Key Management Service Handler
	Port{709, 17}:    "entrustmanager",       // EntrustManager - NorTel DES auth network see 389 tcp
	Port{710, 6}:     "entrust-ash",          // Entrust Administration Service Handler
	Port{710, 17}:    "entrust-ash",          // Entrust Administration Service Handler
	Port{711, 6}:     "cisco-tdp",            // Cisco TDP
	Port{711, 17}:    "cisco-tdp",            // Cisco TDP
	Port{712, 6}:     "tbrpf",                // Missing description for tbrpf
	Port{712, 17}:    "tbrpf",                // TBRPF
	Port{713, 6}:     "iris-xpc",             // IRIS over XPC
	Port{713, 17}:    "iris-xpc",             // IRIS over XPC
	Port{714, 6}:     "iris-xpcs",            // IRIS over XPCS
	Port{714, 17}:    "iris-xpcs",            // IRIS over XPCS
	Port{715, 6}:     "iris-lwz",             // Missing description for iris-lwz
	Port{715, 17}:    "iris-lwz",             // IRIS-LWZ
	Port{716, 6}:     "pana",                 // PANA Messages
	Port{716, 17}:    "pana",                 // PANA Messages
	Port{723, 6}:     "omfs",                 // OpenMosix File System
	Port{729, 6}:     "netviewdm1",           // IBM NetView DM 6000 Server Client
	Port{729, 17}:    "netviewdm1",           // IBM NetView DM 6000 Server Client
	Port{730, 6}:     "netviewdm2",           // IBM NetView DM 6000 send tcp
	Port{730, 17}:    "netviewdm2",           // IBM NetView DM 6000 send tcp
	Port{731, 6}:     "netviewdm3",           // IBM NetView DM 6000 receive tcp
	Port{731, 17}:    "netviewdm3",           // IBM NetView DM 6000 receive tcp
	Port{737, 17}:    "sometimes-rpc2",       // Rusersd on my OpenBSD box
	Port{740, 6}:     "netcp",                // NETscout Control Protocol
	Port{740, 17}:    "netcp",                // NETscout Control Protocol
	Port{741, 6}:     "netgw",                // Missing description for netgw
	Port{741, 17}:    "netgw",                // Missing description for netgw
	Port{742, 6}:     "netrcs",               // Network based Rev. Cont. Sys.
	Port{742, 17}:    "netrcs",               // Network based Rev. Cont. Sys.
	Port{744, 6}:     "flexlm",               // Flexible License Manager
	Port{744, 17}:    "flexlm",               // Flexible License Manager
	Port{747, 6}:     "fujitsu-dev",          // Fujitsu Device Control
	Port{747, 17}:    "fujitsu-dev",          // Fujitsu Device Control
	Port{748, 6}:     "ris-cm",               // Russell Info Sci Calendar Manager
	Port{748, 17}:    "ris-cm",               // Russell Info Sci Calendar Manager
	Port{749, 6}:     "kerberos-adm",         // Kerberos 5 admin changepw | kerberos administration
	Port{749, 17}:    "kerberos-adm",         // Kerberos 5 admin changepw
	Port{750, 6}:     "kerberos",             // kerberos-iv | loadav | rfile | kdc Kerberos (v4) | kerberos version iv
	Port{750, 17}:    "kerberos",             // kdc Kerberos (v4)
	Port{751, 6}:     "kerberos_master",      // pump | Kerberos `kadmin' (v4)
	Port{751, 17}:    "kerberos_master",      // Kerberos `kadmin' (v4)
	Port{752, 6}:     "qrh",                  // Missing description for qrh
	Port{752, 17}:    "qrh",                  // Missing description for qrh
	Port{753, 6}:     "rrh",                  // Missing description for rrh
	Port{753, 17}:    "rrh",                  // Missing description for rrh
	Port{754, 6}:     "krb_prop",             // tell | kerberos v5 server propagation | send
	Port{754, 17}:    "tell",                 // send
	Port{758, 6}:     "nlogin",               // Missing description for nlogin
	Port{758, 17}:    "nlogin",               // Missing description for nlogin
	Port{759, 6}:     "con",                  // Missing description for con
	Port{759, 17}:    "con",                  // Missing description for con
	Port{760, 6}:     "krbupdate",            // ns | kreg Kerberos (v4) registration
	Port{760, 17}:    "ns",                   // Missing description for ns
	Port{761, 6}:     "kpasswd",              // rxe | kpwd Kerberos (v4) "passwd"
	Port{761, 17}:    "rxe",                  // Missing description for rxe
	Port{762, 6}:     "quotad",               // Missing description for quotad
	Port{762, 17}:    "quotad",               // Missing description for quotad
	Port{763, 6}:     "cycleserv",            // Missing description for cycleserv
	Port{763, 17}:    "cycleserv",            // Missing description for cycleserv
	Port{764, 6}:     "omserv",               // Missing description for omserv
	Port{764, 17}:    "omserv",               // Missing description for omserv
	Port{765, 6}:     "webster",              // Missing description for webster
	Port{765, 17}:    "webster",              // Missing description for webster
	Port{767, 6}:     "phonebook",            // phone
	Port{767, 17}:    "phonebook",            // phone
	Port{769, 6}:     "vid",                  // Missing description for vid
	Port{769, 17}:    "vid",                  // Missing description for vid
	Port{770, 6}:     "cadlock",              // Missing description for cadlock
	Port{770, 17}:    "cadlock",              // Missing description for cadlock
	Port{771, 6}:     "rtip",                 // Missing description for rtip
	Port{771, 17}:    "rtip",                 // Missing description for rtip
	Port{772, 6}:     "cycleserv2",           // Missing description for cycleserv2
	Port{772, 17}:    "cycleserv2",           // Missing description for cycleserv2
	Port{773, 6}:     "submit",               // notify
	Port{773, 17}:    "notify",               // Missing description for notify
	Port{774, 6}:     "rpasswd",              // acmaint_dbd | acmaint-dbd
	Port{774, 17}:    "acmaint_dbd",          // Missing description for acmaint_dbd
	Port{775, 6}:     "entomb",               // acmaint_transd | acmaint-transd
	Port{775, 17}:    "acmaint_transd",       // Missing description for acmaint_transd
	Port{776, 6}:     "wpages",               // Missing description for wpages
	Port{776, 17}:    "wpages",               // Missing description for wpages
	Port{777, 6}:     "multiling-http",       // Multiling HTTP
	Port{777, 17}:    "multiling-http",       // Multiling HTTP
	Port{780, 6}:     "wpgs",                 // Missing description for wpgs
	Port{780, 17}:    "wpgs",                 // Missing description for wpgs
	Port{781, 6}:     "hp-collector",         // hp performance data collector
	Port{781, 17}:    "hp-collector",         // hp performance data collector
	Port{782, 6}:     "hp-managed-node",      // hp performance data managed node
	Port{782, 17}:    "hp-managed-node",      // hp performance data managed node
	Port{783, 6}:     "spamassassin",         // Apache SpamAssassin spamd
	Port{786, 6}:     "concert",              // Missing description for concert
	Port{786, 17}:    "concert",              // Missing description for concert
	Port{787, 6}:     "qsc",                  // Missing description for qsc
	Port{799, 6}:     "controlit",            // Remotely possible
	Port{800, 6}:     "mdbs_daemon",          // mdbs-daemon
	Port{800, 17}:    "mdbs_daemon",          // Missing description for mdbs_daemon
	Port{801, 6}:     "device",               // Missing description for device
	Port{801, 17}:    "device",               // Missing description for device
	Port{802, 6}:     "mbap-s",               // Modbus Application Protocol Secure
	Port{808, 6}:     "ccproxy-http",         // CCProxy HTTP Gopher FTP (over HTTP) proxy
	Port{810, 6}:     "fcp-udp",              // FCP | FCP Datagram
	Port{810, 17}:    "fcp-udp",              // FCP Datagram
	Port{828, 6}:     "itm-mcell-s",          // Missing description for itm-mcell-s
	Port{828, 17}:    "itm-mcell-s",          // Missing description for itm-mcell-s
	Port{829, 6}:     "pkix-3-ca-ra",         // PKIX-3 CA RA
	Port{829, 17}:    "pkix-3-ca-ra",         // PKIX-3 CA RA
	Port{830, 6}:     "netconf-ssh",          // NETCONF over SSH
	Port{830, 17}:    "netconf-ssh",          // NETCONF over SSH
	Port{831, 6}:     "netconf-beep",         // NETCONF over BEEP
	Port{831, 17}:    "netconf-beep",         // NETCONF over BEEP
	Port{832, 6}:     "netconfsoaphttp",      // NETCONF for SOAP over HTTPS
	Port{832, 17}:    "netconfsoaphttp",      // NETCONF for SOAP over HTTPS
	Port{833, 6}:     "netconfsoapbeep",      // NETCONF for SOAP over BEEP
	Port{833, 17}:    "netconfsoapbeep",      // NETCONF for SOAP over BEEP
	Port{847, 6}:     "dhcp-failover2",       // dhcp-failover 2
	Port{847, 17}:    "dhcp-failover2",       // dhcp-failover 2
	Port{848, 6}:     "gdoi",                 // Missing description for gdoi
	Port{848, 17}:    "gdoi",                 // GDOI
	Port{853, 6}:     "domain-s",             // DNS query-response protocol run over TLS DTLS
	Port{854, 6}:     "dlep",                 // Dynamic Link Exchange Protocol (DLEP)
	Port{860, 6}:     "iscsi",                // Missing description for iscsi
	Port{860, 17}:    "iscsi",                // iSCSI
	Port{861, 6}:     "owamp-control",        // Missing description for owamp-control
	Port{861, 17}:    "owamp-control",        // OWAMP-Control
	Port{862, 6}:     "twamp-control",        // Two-way Active Measurement Protocol (TWAMP) Control
	Port{862, 17}:    "twamp-control",        // Two-way Active Measurement Protocol (TWAMP) Control
	Port{871, 6}:     "supfilesrv",           // SUP server
	Port{873, 6}:     "rsync",                // Rsync server ( http:  rsync.samba.org )
	Port{873, 17}:    "rsync",                // Missing description for rsync
	Port{886, 6}:     "iclcnet-locate",       // ICL coNETion locate server
	Port{886, 17}:    "iclcnet-locate",       // ICL coNETion locate server
	Port{887, 6}:     "iclcnet_svinfo",       // iclcnet-svinfo | ICL coNETion server info
	Port{887, 17}:    "iclcnet_svinfo",       // ICL coNETion server info
	Port{888, 6}:     "accessbuilder",        // cddbp | or Audio CD Database | CD Database Protocol
	Port{888, 17}:    "accessbuilder",        // Missing description for accessbuilder
	Port{898, 6}:     "sun-manageconsole",    // Solaris Management Console Java listener (Solaris 8 & 9)
	Port{900, 6}:     "omginitialrefs",       // OMG Initial Refs
	Port{900, 17}:    "omginitialrefs",       // OMG Initial Refs
	Port{901, 6}:     "samba-swat",           // smpnameres | Samba SWAT tool.  Also used by ISS RealSecure. | SMPNAMERES
	Port{901, 17}:    "smpnameres",           // SMPNAMERES
	Port{902, 6}:     "iss-realsecure",       // ideafarm-door | ISS RealSecure Sensor | self documenting Telnet Door | self documenting Door: send 0x00 for info
	Port{902, 17}:    "ideafarm-door",        // self documenting Door: send 0x00 for info
	Port{903, 6}:     "iss-console-mgr",      // ideafarm-panic | ISS Console Manager | self documenting Telnet Panic Door | self documenting Panic Door: send 0x00 for info
	Port{903, 17}:    "ideafarm-panic",       // self documenting Panic Door: send 0x00 for info
	Port{910, 6}:     "kink",                 // Kerberized Internet Negotiation of Keys (KINK)
	Port{910, 17}:    "kink",                 // Kerberized Internet Negotiation of Keys (KINK)
	Port{911, 6}:     "xact-backup",          // Missing description for xact-backup
	Port{911, 17}:    "xact-backup",          // Missing description for xact-backup
	Port{912, 6}:     "apex-mesh",            // APEX relay-relay service
	Port{912, 17}:    "apex-mesh",            // APEX relay-relay service
	Port{913, 6}:     "apex-edge",            // APEX endpoint-relay service
	Port{913, 17}:    "apex-edge",            // APEX endpoint-relay service
	Port{950, 6}:     "oftep-rpc",            // Often RPC.statd (on Redhat Linux)
	Port{953, 6}:     "rndc",                 // RNDC is used by BIND 9 (& probably other NS) | BIND9 remote name daemon controller
	Port{975, 6}:     "securenetpro-sensor",  // Missing description for securenetpro-sensor
	Port{989, 6}:     "ftps-data",            // ftp protocol, data, over TLS SSL
	Port{989, 17}:    "ftps-data",            // ftp protocol, data, over TLS SSL
	Port{990, 6}:     "ftps",                 // ftp protocol, control, over TLS SSL
	Port{990, 17}:    "ftps",                 // ftp protocol, control, over TLS SSL
	Port{991, 6}:     "nas",                  // Netnews Administration System
	Port{991, 17}:    "nas",                  // Netnews Administration System
	Port{992, 6}:     "telnets",              // telnet protocol over TLS SSL
	Port{992, 17}:    "telnets",              // telnet protocol over TLS SSL
	Port{993, 6}:     "imaps",                // imap4 protocol over TLS SSL | IMAP over TLS protocol
	Port{993, 17}:    "imaps",                // imap4 protocol over TLS SSL
	Port{994, 6}:     "ircs",                 // irc protocol over TLS SSL
	Port{994, 17}:    "ircs",                 // irc protocol over TLS SSL
	Port{995, 6}:     "pop3s",                // POP3 protocol over TLS SSL | pop3 protocol over TLS SSL (was spop3) | POP3 over TLS protocol
	Port{995, 17}:    "pop3s",                // pop3 protocol over TLS SSL (was spop3)
	Port{996, 6}:     "xtreelic",             // XTREE License Server | vsinet
	Port{996, 17}:    "vsinet",               // Missing description for vsinet
	Port{997, 6}:     "maitrd",               // Missing description for maitrd
	Port{997, 17}:    "maitrd",               // Missing description for maitrd
	Port{998, 6}:     "busboy",               // puparp
	Port{998, 17}:    "puparp",               // Missing description for puparp
	Port{999, 6}:     "garcon",               // puprouter | applix | Applix ac
	Port{999, 17}:    "applix",               // Applix ac
	Port{1000, 6}:    "cadlock",              // cadlock2
	Port{1000, 17}:   "ock",                  // Missing description for ock
	Port{1001, 6}:    "webpush",              // HTTP Web Push
	Port{1002, 6}:    "windows-icfw",         // Windows Internet Connection Firewall or Internet Locator Server for NetMeeting.
	Port{1008, 6}:    "ufsd",                 // ufsd  UFS-aware server
	Port{1008, 17}:   "ufsd",                 // Missing description for ufsd
	Port{1010, 6}:    "surf",                 // Missing description for surf
	Port{1010, 17}:   "surf",                 // Missing description for surf
	Port{1012, 17}:   "sometimes-rpc1",       // This is rstatd on my openBSD box
	Port{1021, 6}:    "exp1",                 // RFC3692-style Experiment 1 (*)    [RFC4727] | RFC3692-style Experiment 1
	Port{1021, 17}:   "exp1",                 // RFC3692-style Experiment 1 (*)    [RFC4727]
	Port{1022, 6}:    "exp2",                 // RFC3692-style Experiment 2 (*)    [RFC4727] | RFC3692-style Experiment 2
	Port{1022, 17}:   "exp2",                 // RFC3692-style Experiment 2 (*)    [RFC4727]
	Port{1023, 6}:    "netvenuechat",         // Nortel NetVenue Notification, Chat, Intercom
	Port{1024, 6}:    "kdm",                  // K Display Manager (KDE version of xdm)
	Port{1025, 6}:    "NFS-or-IIS",           // blackjack | IIS, NFS, or listener RFS remote_file_sharing | network blackjack
	Port{1025, 17}:   "blackjack",            // network blackjack
	Port{1026, 6}:    "LSA-or-nterm",         // cap | nterm remote_login network_terminal | Calendar Access Protocol
	Port{1026, 17}:   "win-rpc",              // Commonly used to send MS Messenger spam
	Port{1027, 6}:    "IIS",                  // 6a44 | IPv6 Behind NAT44 CPEs
	Port{1028, 17}:   "ms-lsa",               // Missing description for ms-lsa
	Port{1029, 6}:    "ms-lsa",               // solid-mux | Solid Mux Server
	Port{1029, 17}:   "solid-mux",            // Solid Mux Server
	Port{1030, 6}:    "iad1",                 // BBN IAD
	Port{1030, 17}:   "iad1",                 // BBN IAD
	Port{1031, 6}:    "iad2",                 // BBN IAD
	Port{1031, 17}:   "iad2",                 // BBN IAD
	Port{1032, 6}:    "iad3",                 // BBN IAD
	Port{1032, 17}:   "iad3",                 // BBN IAD
	Port{1033, 6}:    "netinfo",              // netinfo-local | Netinfo is apparently on many OS X boxes. | local netinfo port
	Port{1033, 17}:   "netinfo-local",        // local netinfo port
	Port{1034, 6}:    "zincite-a",            // activesync | Zincite.A backdoor | ActiveSync Notifications
	Port{1034, 17}:   "activesync-notify",    // Windows Mobile device ActiveSync Notifications
	Port{1035, 6}:    "multidropper",         // mxxrlogin | A Multidropper Adware, or PhoneFree | MX-XR RPC
	Port{1035, 17}:   "mxxrlogin",            // MX-XR RPC
	Port{1036, 6}:    "nsstp",                // Nebula Secure Segment Transfer Protocol
	Port{1036, 17}:   "nsstp",                // Nebula Secure Segment Transfer Protocol
	Port{1037, 6}:    "ams",                  // Missing description for ams
	Port{1037, 17}:   "ams",                  // AMS
	Port{1038, 6}:    "mtqp",                 // Message Tracking Query Protocol
	Port{1038, 17}:   "mtqp",                 // Message Tracking Query Protocol
	Port{1039, 6}:    "sbl",                  // Streamlined Blackhole
	Port{1039, 17}:   "sbl",                  // Streamlined Blackhole
	Port{1040, 6}:    "netsaint",             // netarx | Netsaint status daemon | Netarx Netcare
	Port{1040, 17}:   "netarx",               // Netarx Netcare
	Port{1041, 6}:    "danf-ak2",             // AK2 Product
	Port{1041, 17}:   "danf-ak2",             // AK2 Product
	Port{1042, 6}:    "afrog",                // Subnet Roaming
	Port{1042, 17}:   "afrog",                // Subnet Roaming
	Port{1043, 6}:    "boinc",                // boinc-client | BOINC Client Control or Microsoft IIS | BOINC Client Control
	Port{1043, 17}:   "boinc",                // BOINC Client Control
	Port{1044, 6}:    "dcutility",            // Dev Consortium Utility
	Port{1044, 17}:   "dcutility",            // Dev Consortium Utility
	Port{1045, 6}:    "fpitp",                // Fingerprint Image Transfer Protocol
	Port{1045, 17}:   "fpitp",                // Fingerprint Image Transfer Protocol
	Port{1046, 6}:    "wfremotertm",          // WebFilter Remote Monitor
	Port{1046, 17}:   "wfremotertm",          // WebFilter Remote Monitor
	Port{1047, 6}:    "neod1",                // Sun's NEO Object Request Broker
	Port{1047, 17}:   "neod1",                // Sun's NEO Object Request Broker
	Port{1048, 6}:    "neod2",                // Sun's NEO Object Request Broker
	Port{1048, 17}:   "neod2",                // Sun's NEO Object Request Broker
	Port{1049, 6}:    "td-postman",           // Tobit David Postman VPMN
	Port{1049, 17}:   "td-postman",           // Tobit David Postman VPMN
	Port{1050, 6}:    "java-or-OTGfileshare", // cma | J2EE nameserver, also OTG, also called Disk Application extender. Could also be MiniCommand backdoor OTGlicenseserv | CORBA Management Agent
	Port{1050, 17}:   "cma",                  // CORBA Management Agent
	Port{1051, 6}:    "optima-vnet",          // Optima VNET
	Port{1051, 17}:   "optima-vnet",          // Missing description for optima-vnet
	Port{1052, 6}:    "ddt",                  // Dynamic DNS tools | Dynamic DNS Tools
	Port{1052, 17}:   "ddt",                  // Dynamic DNS tools
	Port{1053, 6}:    "remote-as",            // Remote Assistant (RA)
	Port{1053, 17}:   "remote-as",            // Remote Assistant (RA)
	Port{1054, 6}:    "brvread",              // Missing description for brvread
	Port{1054, 17}:   "brvread",              // BRVREAD
	Port{1055, 6}:    "ansyslmd",             // ANSYS - License Manager
	Port{1055, 17}:   "ansyslmd",             // Missing description for ansyslmd
	Port{1056, 6}:    "vfo",                  // Missing description for vfo
	Port{1056, 17}:   "vfo",                  // VFO
	Port{1057, 6}:    "startron",             // Missing description for startron
	Port{1057, 17}:   "startron",             // STARTRON
	Port{1058, 6}:    "nim",                  // Missing description for nim
	Port{1058, 17}:   "nim",                  // Missing description for nim
	Port{1059, 6}:    "nimreg",               // Missing description for nimreg
	Port{1059, 17}:   "nimreg",               // Missing description for nimreg
	Port{1060, 6}:    "polestar",             // Missing description for polestar
	Port{1060, 17}:   "polestar",             // Missing description for polestar
	Port{1061, 6}:    "kiosk",                // Missing description for kiosk
	Port{1061, 17}:   "kiosk",                // KIOSK
	Port{1062, 6}:    "veracity",             // Missing description for veracity
	Port{1062, 17}:   "veracity",             // Missing description for veracity
	Port{1063, 6}:    "kyoceranetdev",        // Missing description for kyoceranetdev
	Port{1063, 17}:   "kyoceranetdev",        // KyoceraNetDev
	Port{1064, 6}:    "jstel",                // Missing description for jstel
	Port{1064, 17}:   "jstel",                // JSTEL
	Port{1065, 6}:    "syscomlan",            // Missing description for syscomlan
	Port{1065, 17}:   "syscomlan",            // SYSCOMLAN
	Port{1066, 6}:    "fpo-fns",              // Missing description for fpo-fns
	Port{1066, 17}:   "fpo-fns",              // Missing description for fpo-fns
	Port{1067, 6}:    "instl_boots",          // instl-boots | Installation Bootstrap Proto. Serv.
	Port{1067, 17}:   "instl_boots",          // Installation Bootstrap Proto. Serv.
	Port{1068, 6}:    "instl_bootc",          // instl-bootc | Installation Bootstrap Proto. Cli.
	Port{1068, 17}:   "instl_bootc",          // Installation Bootstrap Proto. Cli.
	Port{1069, 6}:    "cognex-insight",       // Missing description for cognex-insight
	Port{1069, 17}:   "cognex-insight",       // Missing description for cognex-insight
	Port{1070, 6}:    "gmrupdateserv",        // Missing description for gmrupdateserv
	Port{1070, 17}:   "gmrupdateserv",        // GMRUpdateSERV
	Port{1071, 6}:    "bsquare-voip",         // Missing description for bsquare-voip
	Port{1071, 17}:   "bsquare-voip",         // BSQUARE-VOIP
	Port{1072, 6}:    "cardax",               // Missing description for cardax
	Port{1072, 17}:   "cardax",               // CARDAX
	Port{1073, 6}:    "bridgecontrol",        // Bridge Control
	Port{1073, 17}:   "bridgecontrol",        // Bridge Control
	Port{1074, 6}:    "warmspotMgmt",         // Warmspot Management Protocol
	Port{1074, 17}:   "warmspotMgmt",         // Warmspot Management Protocol
	Port{1075, 6}:    "rdrmshc",              // Missing description for rdrmshc
	Port{1075, 17}:   "rdrmshc",              // RDRMSHC
	Port{1076, 6}:    "sns_credit",           // dab-sti-c | Shared Network Services (SNS) for Canadian credit card authorizations | DAB STI-C
	Port{1076, 17}:   "dab-sti-c",            // DAB STI-C
	Port{1077, 6}:    "imgames",              // Missing description for imgames
	Port{1077, 17}:   "imgames",              // IMGames
	Port{1078, 6}:    "avocent-proxy",        // Avocent Proxy Protocol
	Port{1078, 17}:   "avocent-proxy",        // Avocent Proxy Protocol
	Port{1079, 6}:    "asprovatalk",          // Missing description for asprovatalk
	Port{1079, 17}:   "asprovatalk",          // ASPROVATalk
	Port{1080, 6}:    "socks",                // Missing description for socks
	Port{1080, 17}:   "socks",                // Missing description for socks
	Port{1081, 6}:    "pvuniwien",            // Missing description for pvuniwien
	Port{1081, 17}:   "pvuniwien",            // PVUNIWIEN
	Port{1082, 6}:    "amt-esd-prot",         // Missing description for amt-esd-prot
	Port{1082, 17}:   "amt-esd-prot",         // AMT-ESD-PROT
	Port{1083, 6}:    "ansoft-lm-1",          // Anasoft License Manager
	Port{1083, 17}:   "ansoft-lm-1",          // Anasoft License Manager
	Port{1084, 6}:    "ansoft-lm-2",          // Anasoft License Manager
	Port{1084, 17}:   "ansoft-lm-2",          // Anasoft License Manager
	Port{1085, 6}:    "webobjects",           // Web Objects
	Port{1085, 17}:   "webobjects",           // Web Objects
	Port{1086, 6}:    "cplscrambler-lg",      // CPL Scrambler Logging
	Port{1086, 17}:   "cplscrambler-lg",      // CPL Scrambler Logging
	Port{1087, 6}:    "cplscrambler-in",      // CPL Scrambler Internal
	Port{1087, 17}:   "cplscrambler-in",      // CPL Scrambler Internal
	Port{1088, 6}:    "cplscrambler-al",      // CPL Scrambler Alarm Log
	Port{1088, 17}:   "cplscrambler-al",      // CPL Scrambler Alarm Log
	Port{1089, 6}:    "ff-annunc",            // FF Annunciation
	Port{1089, 17}:   "ff-annunc",            // FF Annunciation
	Port{1090, 6}:    "ff-fms",               // FF Fieldbus Message Specification
	Port{1090, 17}:   "ff-fms",               // FF Fieldbus Message Specification
	Port{1091, 6}:    "ff-sm",                // FF System Management
	Port{1091, 17}:   "ff-sm",                // FF System Management
	Port{1092, 6}:    "obrpd",                // Open Business Reporting Protocol
	Port{1092, 17}:   "obrpd",                // Open Business Reporting Protocol
	Port{1093, 6}:    "proofd",               // Missing description for proofd
	Port{1093, 17}:   "proofd",               // PROOFD
	Port{1094, 6}:    "rootd",                // Missing description for rootd
	Port{1094, 17}:   "rootd",                // ROOTD
	Port{1095, 6}:    "nicelink",             // Missing description for nicelink
	Port{1095, 17}:   "nicelink",             // NICELink
	Port{1096, 6}:    "cnrprotocol",          // Common Name Resolution Protocol
	Port{1096, 17}:   "cnrprotocol",          // Common Name Resolution Protocol
	Port{1097, 6}:    "sunclustermgr",        // Sun Cluster Manager
	Port{1097, 17}:   "sunclustermgr",        // Sun Cluster Manager
	Port{1098, 6}:    "rmiactivation",        // RMI Activation
	Port{1098, 17}:   "rmiactivation",        // RMI Activation
	Port{1099, 6}:    "rmiregistry",          // RMI Registry
	Port{1099, 17}:   "rmiregistry",          // RMI Registry
	Port{1100, 6}:    "mctp",                 // Missing description for mctp
	Port{1100, 17}:   "mctp",                 // MCTP
	Port{1101, 6}:    "pt2-discover",         // Missing description for pt2-discover
	Port{1101, 17}:   "pt2-discover",         // PT2-DISCOVER
	Port{1102, 6}:    "adobeserver-1",        // ADOBE SERVER 1
	Port{1102, 17}:   "adobeserver-1",        // ADOBE SERVER 1
	Port{1103, 6}:    "xaudio",               // adobeserver-2 | Xaserver X Audio Server | ADOBE SERVER 2
	Port{1103, 17}:   "adobeserver-2",        // ADOBE SERVER 2
	Port{1104, 6}:    "xrl",                  // Missing description for xrl
	Port{1104, 17}:   "xrl",                  // XRL
	Port{1105, 6}:    "ftranhc",              // Missing description for ftranhc
	Port{1105, 17}:   "ftranhc",              // FTRANHC
	Port{1106, 6}:    "isoipsigport-1",       // Missing description for isoipsigport-1
	Port{1106, 17}:   "isoipsigport-1",       // ISOIPSIGPORT-1
	Port{1107, 6}:    "isoipsigport-2",       // Missing description for isoipsigport-2
	Port{1107, 17}:   "isoipsigport-2",       // ISOIPSIGPORT-2
	Port{1108, 6}:    "ratio-adp",            // Missing description for ratio-adp
	Port{1108, 17}:   "ratio-adp",            // Missing description for ratio-adp
	Port{1109, 6}:    "kpop",                 // Pop with Kerberos
	Port{1110, 6}:    "nfsd-status",          // nfsd-keepalive | webadmstart | Cluster status info | Start web admin server | Client status info
	Port{1110, 17}:   "nfsd-keepalive",       // Client status info
	Port{1111, 6}:    "lmsocialserver",       // LM Social Server
	Port{1111, 17}:   "lmsocialserver",       // LM Social Server
	Port{1112, 6}:    "msql",                 // icp | mini-sql server | Intelligent Communication Protocol
	Port{1112, 17}:   "icp",                  // Intelligent Communication Protocol
	Port{1113, 6}:    "ltp-deepspace",        // Licklider Transmission Protocol
	Port{1113, 17}:   "ltp-deepspace",        // Licklider Transmission Protocol
	Port{1114, 6}:    "mini-sql",             // Mini SQL
	Port{1114, 17}:   "mini-sql",             // Mini SQL
	Port{1115, 6}:    "ardus-trns",           // ARDUS Transfer
	Port{1115, 17}:   "ardus-trns",           // ARDUS Transfer
	Port{1116, 6}:    "ardus-cntl",           // ARDUS Control
	Port{1116, 17}:   "ardus-cntl",           // ARDUS Control
	Port{1117, 6}:    "ardus-mtrns",          // ARDUS Multicast Transfer
	Port{1117, 17}:   "ardus-mtrns",          // ARDUS Multicast Transfer
	Port{1118, 6}:    "sacred",               // Missing description for sacred
	Port{1118, 17}:   "sacred",               // SACRED
	Port{1119, 6}:    "bnetgame",             // Battle.net Chat Game Protocol
	Port{1119, 17}:   "bnetgame",             // Battle.net Chat Game Protocol
	Port{1120, 6}:    "bnetfile",             // Battle.net File Transfer Protocol
	Port{1120, 17}:   "bnetfile",             // Battle.net File Transfer Protocol
	Port{1121, 6}:    "rmpp",                 // Datalode RMPP
	Port{1121, 17}:   "rmpp",                 // Datalode RMPP
	Port{1122, 6}:    "availant-mgr",         // Missing description for availant-mgr
	Port{1122, 17}:   "availant-mgr",         // Missing description for availant-mgr
	Port{1123, 6}:    "murray",               // Missing description for murray
	Port{1123, 17}:   "murray",               // Murray
	Port{1124, 6}:    "hpvmmcontrol",         // HP VMM Control
	Port{1124, 17}:   "hpvmmcontrol",         // HP VMM Control
	Port{1125, 6}:    "hpvmmagent",           // HP VMM Agent
	Port{1125, 17}:   "hpvmmagent",           // HP VMM Agent
	Port{1126, 6}:    "hpvmmdata",            // HP VMM Agent
	Port{1126, 17}:   "hpvmmdata",            // HP VMM Agent
	Port{1127, 6}:    "supfiledbg",           // kwdb-commn | SUP debugging | KWDB Remote Communication
	Port{1127, 17}:   "kwdb-commn",           // KWDB Remote Communication
	Port{1128, 6}:    "saphostctrl",          // SAPHostControl over SOAP HTTP
	Port{1128, 17}:   "saphostctrl",          // SAPHostControl over SOAP HTTP
	Port{1129, 6}:    "saphostctrls",         // SAPHostControl over SOAP HTTPS
	Port{1129, 17}:   "saphostctrls",         // SAPHostControl over SOAP HTTPS
	Port{1130, 6}:    "casp",                 // CAC App Service Protocol
	Port{1130, 17}:   "casp",                 // CAC App Service Protocol
	Port{1131, 6}:    "caspssl",              // CAC App Service Protocol Encripted
	Port{1131, 17}:   "caspssl",              // CAC App Service Protocol Encripted
	Port{1132, 6}:    "kvm-via-ip",           // KVM-via-IP Management Service
	Port{1132, 17}:   "kvm-via-ip",           // KVM-via-IP Management Service
	Port{1133, 6}:    "dfn",                  // Data Flow Network
	Port{1133, 17}:   "dfn",                  // Data Flow Network
	Port{1134, 6}:    "aplx",                 // MicroAPL APLX
	Port{1134, 17}:   "aplx",                 // MicroAPL APLX
	Port{1135, 6}:    "omnivision",           // OmniVision Communication Service
	Port{1135, 17}:   "omnivision",           // OmniVision Communication Service
	Port{1136, 6}:    "hhb-gateway",          // HHB Gateway Control
	Port{1136, 17}:   "hhb-gateway",          // HHB Gateway Control
	Port{1137, 6}:    "trim",                 // TRIM Workgroup Service
	Port{1137, 17}:   "trim",                 // TRIM Workgroup Service
	Port{1138, 6}:    "encrypted_admin",      // encrypted-admin | encrypted admin requests
	Port{1138, 17}:   "encrypted_admin",      // encrypted admin requests
	Port{1139, 6}:    "cce3x",                // evm | ClearCommerce Engine 3.x ( www.clearcommerce.com) | Enterprise Virtual Manager
	Port{1139, 17}:   "evm",                  // Enterprise Virtual Manager
	Port{1140, 6}:    "autonoc",              // AutoNOC Network Operations Protocol
	Port{1140, 17}:   "autonoc",              // AutoNOC Network Operations Protocol
	Port{1141, 6}:    "mxomss",               // User Message Service
	Port{1141, 17}:   "mxomss",               // User Message Service
	Port{1142, 6}:    "edtools",              // User Discovery Service
	Port{1142, 17}:   "edtools",              // User Discovery Service
	Port{1143, 6}:    "imyx",                 // Infomatryx Exchange
	Port{1143, 17}:   "imyx",                 // Infomatryx Exchange
	Port{1144, 6}:    "fuscript",             // Fusion Script
	Port{1144, 17}:   "fuscript",             // Fusion Script
	Port{1145, 6}:    "x9-icue",              // X9 iCue Show Control
	Port{1145, 17}:   "x9-icue",              // X9 iCue Show Control
	Port{1146, 6}:    "audit-transfer",       // audit transfer
	Port{1146, 17}:   "audit-transfer",       // audit transfer
	Port{1147, 6}:    "capioverlan",          // Missing description for capioverlan
	Port{1147, 17}:   "capioverlan",          // CAPIoverLAN
	Port{1148, 6}:    "elfiq-repl",           // Elfiq Replication Service
	Port{1148, 17}:   "elfiq-repl",           // Elfiq Replication Service
	Port{1149, 6}:    "bvtsonar",             // BVT Sonar Service | BlueView Sonar Service
	Port{1149, 17}:   "bvtsonar",             // BVT Sonar Service
	Port{1150, 6}:    "blaze",                // Blaze File Server
	Port{1150, 17}:   "blaze",                // Blaze File Server
	Port{1151, 6}:    "unizensus",            // Unizensus Login Server
	Port{1151, 17}:   "unizensus",            // Unizensus Login Server
	Port{1152, 6}:    "winpoplanmess",        // Winpopup LAN Messenger
	Port{1152, 17}:   "winpoplanmess",        // Winpopup LAN Messenger
	Port{1153, 6}:    "c1222-acse",           // ANSI C12.22 Port
	Port{1153, 17}:   "c1222-acse",           // ANSI C12.22 Port
	Port{1154, 6}:    "resacommunity",        // Community Service
	Port{1154, 17}:   "resacommunity",        // Community Service
	Port{1155, 6}:    "nfa",                  // Network File Access
	Port{1155, 17}:   "nfa",                  // Network File Access
	Port{1156, 6}:    "iascontrol-oms",       // iasControl OMS
	Port{1156, 17}:   "iascontrol-oms",       // iasControl OMS
	Port{1157, 6}:    "iascontrol",           // Oracle iASControl
	Port{1157, 17}:   "iascontrol",           // Oracle iASControl
	Port{1158, 6}:    "lsnr",                 // dbcontrol-oms | Oracle DB listener | dbControl OMS
	Port{1158, 17}:   "dbcontrol-oms",        // dbControl OMS
	Port{1159, 6}:    "oracle-oms",           // Oracle OMS
	Port{1159, 17}:   "oracle-oms",           // Oracle OMS
	Port{1160, 6}:    "olsv",                 // DB Lite Mult-User Server
	Port{1160, 17}:   "olsv",                 // DB Lite Mult-User Server
	Port{1161, 6}:    "health-polling",       // Health Polling
	Port{1161, 17}:   "health-polling",       // Health Polling
	Port{1162, 6}:    "health-trap",          // Health Trap
	Port{1162, 17}:   "health-trap",          // Health Trap
	Port{1163, 6}:    "sddp",                 // SmartDialer Data Protocol
	Port{1163, 17}:   "sddp",                 // SmartDialer Data Protocol
	Port{1164, 6}:    "qsm-proxy",            // QSM Proxy Service
	Port{1164, 17}:   "qsm-proxy",            // QSM Proxy Service
	Port{1165, 6}:    "qsm-gui",              // QSM GUI Service
	Port{1165, 17}:   "qsm-gui",              // QSM GUI Service
	Port{1166, 6}:    "qsm-remote",           // QSM RemoteExec
	Port{1166, 17}:   "qsm-remote",           // QSM RemoteExec
	Port{1167, 132}:  "cisco-ipsla",          // Cisco IP SLAs Control Protocol
	Port{1167, 6}:    "cisco-ipsla",          // Cisco IP SLAs Control Protocol
	Port{1167, 17}:   "cisco-ipsla",          // Cisco IP SLAs Control Protocol
	Port{1168, 6}:    "vchat",                // VChat Conference Service
	Port{1168, 17}:   "vchat",                // VChat Conference Service
	Port{1169, 6}:    "tripwire",             // Missing description for tripwire
	Port{1169, 17}:   "tripwire",             // TRIPWIRE
	Port{1170, 6}:    "atc-lm",               // AT+C License Manager
	Port{1170, 17}:   "atc-lm",               // AT+C License Manager
	Port{1171, 6}:    "atc-appserver",        // AT+C FmiApplicationServer
	Port{1171, 17}:   "atc-appserver",        // AT+C FmiApplicationServer
	Port{1172, 6}:    "dnap",                 // DNA Protocol
	Port{1172, 17}:   "dnap",                 // DNA Protocol
	Port{1173, 6}:    "d-cinema-rrp",         // D-Cinema Request-Response
	Port{1173, 17}:   "d-cinema-rrp",         // D-Cinema Request-Response
	Port{1174, 6}:    "fnet-remote-ui",       // FlashNet Remote Admin
	Port{1174, 17}:   "fnet-remote-ui",       // FlashNet Remote Admin
	Port{1175, 6}:    "dossier",              // Dossier Server
	Port{1175, 17}:   "dossier",              // Dossier Server
	Port{1176, 6}:    "indigo-server",        // Indigo Home Server
	Port{1176, 17}:   "indigo-server",        // Indigo Home Server
	Port{1177, 6}:    "dkmessenger",          // DKMessenger Protocol
	Port{1177, 17}:   "dkmessenger",          // DKMessenger Protocol
	Port{1178, 6}:    "skkserv",              // sgi-storman | SKK (kanji input) | SGI Storage Manager
	Port{1178, 17}:   "sgi-storman",          // SGI Storage Manager
	Port{1179, 6}:    "b2n",                  // Backup To Neighbor
	Port{1179, 17}:   "b2n",                  // Backup To Neighbor
	Port{1180, 6}:    "mc-client",            // Millicent Client Proxy
	Port{1180, 17}:   "mc-client",            // Millicent Client Proxy
	Port{1181, 6}:    "3comnetman",           // 3Com Net Management
	Port{1181, 17}:   "3comnetman",           // 3Com Net Management
	Port{1182, 6}:    "accelenet",            // accelenet-data | AcceleNet Control | AcceleNet Data
	Port{1182, 17}:   "accelenet-data",       // AcceleNet Data
	Port{1183, 6}:    "llsurfup-http",        // LL Surfup HTTP
	Port{1183, 17}:   "llsurfup-http",        // LL Surfup HTTP
	Port{1184, 6}:    "llsurfup-https",       // LL Surfup HTTPS
	Port{1184, 17}:   "llsurfup-https",       // LL Surfup HTTPS
	Port{1185, 6}:    "catchpole",            // Catchpole port
	Port{1185, 17}:   "catchpole",            // Catchpole port
	Port{1186, 6}:    "mysql-cluster",        // MySQL Cluster Manager
	Port{1186, 17}:   "mysql-cluster",        // MySQL Cluster Manager
	Port{1187, 6}:    "alias",                // Alias Service
	Port{1187, 17}:   "alias",                // Alias Service
	Port{1188, 6}:    "hp-webadmin",          // HP Web Admin
	Port{1188, 17}:   "hp-webadmin",          // HP Web Admin
	Port{1189, 6}:    "unet",                 // Unet Connection
	Port{1189, 17}:   "unet",                 // Unet Connection
	Port{1190, 6}:    "commlinx-avl",         // CommLinx GPS   AVL System
	Port{1190, 17}:   "commlinx-avl",         // CommLinx GPS   AVL System
	Port{1191, 6}:    "gpfs",                 // General Parallel File System
	Port{1191, 17}:   "gpfs",                 // General Parallel File System
	Port{1192, 6}:    "caids-sensor",         // caids sensors channel
	Port{1192, 17}:   "caids-sensor",         // caids sensors channel
	Port{1193, 6}:    "fiveacross",           // Five Across Server
	Port{1193, 17}:   "fiveacross",           // Five Across Server
	Port{1194, 6}:    "openvpn",              // Missing description for openvpn
	Port{1194, 17}:   "openvpn",              // OpenVPN
	Port{1195, 6}:    "rsf-1",                // RSF-1 clustering
	Port{1195, 17}:   "rsf-1",                // RSF-1 clustering
	Port{1196, 6}:    "netmagic",             // Network Magic
	Port{1196, 17}:   "netmagic",             // Network Magic
	Port{1197, 6}:    "carrius-rshell",       // Carrius Remote Access
	Port{1197, 17}:   "carrius-rshell",       // Carrius Remote Access
	Port{1198, 6}:    "cajo-discovery",       // cajo reference discovery
	Port{1198, 17}:   "cajo-discovery",       // cajo reference discovery
	Port{1199, 6}:    "dmidi",                // Missing description for dmidi
	Port{1199, 17}:   "dmidi",                // DMIDI
	Port{1200, 6}:    "scol",                 // Missing description for scol
	Port{1200, 17}:   "scol",                 // SCOL
	Port{1201, 6}:    "nucleus-sand",         // Nucleus Sand Database Server
	Port{1201, 17}:   "nucleus-sand",         // Nucleus Sand Database Server
	Port{1202, 6}:    "caiccipc",             // Missing description for caiccipc
	Port{1202, 17}:   "caiccipc",             // Missing description for caiccipc
	Port{1203, 6}:    "ssslic-mgr",           // License Validation
	Port{1203, 17}:   "ssslic-mgr",           // License Validation
	Port{1204, 6}:    "ssslog-mgr",           // Log Request Listener
	Port{1204, 17}:   "ssslog-mgr",           // Log Request Listener
	Port{1205, 6}:    "accord-mgc",           // Missing description for accord-mgc
	Port{1205, 17}:   "accord-mgc",           // Accord-MGC
	Port{1206, 6}:    "anthony-data",         // Anthony Data
	Port{1206, 17}:   "anthony-data",         // Anthony Data
	Port{1207, 6}:    "metasage",             // Missing description for metasage
	Port{1207, 17}:   "metasage",             // MetaSage
	Port{1208, 6}:    "seagull-ais",          // SEAGULL AIS
	Port{1208, 17}:   "seagull-ais",          // SEAGULL AIS
	Port{1209, 6}:    "ipcd3",                // Missing description for ipcd3
	Port{1209, 17}:   "ipcd3",                // IPCD3
	Port{1210, 6}:    "eoss",                 // Missing description for eoss
	Port{1210, 17}:   "eoss",                 // EOSS
	Port{1211, 6}:    "groove-dpp",           // Groove DPP
	Port{1211, 17}:   "groove-dpp",           // Groove DPP
	Port{1212, 6}:    "lupa",                 // Missing description for lupa
	Port{1212, 17}:   "lupa",                 // Missing description for lupa
	Port{1213, 6}:    "mpc-lifenet",          // MPC LIFENET | Medtronic Physio-Control LIFENET
	Port{1213, 17}:   "mpc-lifenet",          // MPC LIFENET
	Port{1214, 6}:    "fasttrack",            // kazaa | Kazaa File Sharing | KAZAA
	Port{1214, 17}:   "fasttrack",            // Kazaa File Sharing
	Port{1215, 6}:    "scanstat-1",           // scanSTAT 1.0
	Port{1215, 17}:   "scanstat-1",           // scanSTAT 1.0
	Port{1216, 6}:    "etebac5",              // ETEBAC 5
	Port{1216, 17}:   "etebac5",              // ETEBAC 5
	Port{1217, 6}:    "hpss-ndapi",           // HPSS NonDCE Gateway
	Port{1217, 17}:   "hpss-ndapi",           // HPSS NonDCE Gateway
	Port{1218, 6}:    "aeroflight-ads",       // AeroFlight ADs
	Port{1218, 17}:   "aeroflight-ads",       // AeroFlight ADs
	Port{1219, 6}:    "aeroflight-ret",       // Missing description for aeroflight-ret
	Port{1219, 17}:   "aeroflight-ret",       // AeroFlight-Ret
	Port{1220, 6}:    "quicktime",            // qt-serveradmin | Apple Darwin and QuickTime Streaming Administration Servers | QT SERVER ADMIN
	Port{1220, 17}:   "qt-serveradmin",       // QT SERVER ADMIN
	Port{1221, 6}:    "sweetware-apps",       // SweetWARE Apps
	Port{1221, 17}:   "sweetware-apps",       // SweetWARE Apps
	Port{1222, 6}:    "nerv",                 // SNI R&D network
	Port{1222, 17}:   "nerv",                 // SNI R&D network
	Port{1223, 6}:    "tgp",                  // TrulyGlobal Protocol
	Port{1223, 17}:   "tgp",                  // TrulyGlobal Protocol
	Port{1224, 6}:    "vpnz",                 // Missing description for vpnz
	Port{1224, 17}:   "vpnz",                 // VPNz
	Port{1225, 6}:    "slinkysearch",         // Missing description for slinkysearch
	Port{1225, 17}:   "slinkysearch",         // SLINKYSEARCH
	Port{1226, 6}:    "stgxfws",              // Missing description for stgxfws
	Port{1226, 17}:   "stgxfws",              // STGXFWS
	Port{1227, 6}:    "dns2go",               // Missing description for dns2go
	Port{1227, 17}:   "dns2go",               // DNS2Go
	Port{1228, 6}:    "florence",             // Missing description for florence
	Port{1228, 17}:   "florence",             // FLORENCE
	Port{1229, 6}:    "zented",               // ZENworks Tiered Electronic Distribution
	Port{1229, 17}:   "zented",               // ZENworks Tiered Electronic Distribution
	Port{1230, 6}:    "periscope",            // Missing description for periscope
	Port{1230, 17}:   "periscope",            // Periscope
	Port{1231, 6}:    "menandmice-lpm",       // Missing description for menandmice-lpm
	Port{1231, 17}:   "menandmice-lpm",       // Missing description for menandmice-lpm
	Port{1232, 6}:    "first-defense",        // Remote systems monitoring
	Port{1233, 6}:    "univ-appserver",       // Universal App Server
	Port{1233, 17}:   "univ-appserver",       // Universal App Server
	Port{1234, 6}:    "hotline",              // search-agent | Infoseek Search Agent
	Port{1234, 17}:   "search-agent",         // Infoseek Search Agent
	Port{1235, 6}:    "mosaicsyssvc1",        // Missing description for mosaicsyssvc1
	Port{1235, 17}:   "mosaicsyssvc1",        // Missing description for mosaicsyssvc1
	Port{1236, 6}:    "bvcontrol",            // Missing description for bvcontrol
	Port{1236, 17}:   "bvcontrol",            // Missing description for bvcontrol
	Port{1237, 6}:    "tsdos390",             // Missing description for tsdos390
	Port{1237, 17}:   "tsdos390",             // Missing description for tsdos390
	Port{1238, 6}:    "hacl-qs",              // Missing description for hacl-qs
	Port{1238, 17}:   "hacl-qs",              // Missing description for hacl-qs
	Port{1239, 6}:    "nmsd",                 // Missing description for nmsd
	Port{1239, 17}:   "nmsd",                 // NMSD
	Port{1240, 6}:    "instantia",            // Missing description for instantia
	Port{1240, 17}:   "instantia",            // Instantia
	Port{1241, 6}:    "nessus",               // Nessus or remote message server
	Port{1241, 17}:   "nessus",               // Missing description for nessus
	Port{1242, 6}:    "nmasoverip",           // NMAS over IP
	Port{1242, 17}:   "nmasoverip",           // NMAS over IP
	Port{1243, 6}:    "serialgateway",        // Missing description for serialgateway
	Port{1243, 17}:   "serialgateway",        // SerialGateway
	Port{1244, 6}:    "isbconference1",       // Missing description for isbconference1
	Port{1244, 17}:   "isbconference1",       // Missing description for isbconference1
	Port{1245, 6}:    "isbconference2",       // Missing description for isbconference2
	Port{1245, 17}:   "isbconference2",       // Missing description for isbconference2
	Port{1246, 6}:    "payrouter",            // Missing description for payrouter
	Port{1246, 17}:   "payrouter",            // Missing description for payrouter
	Port{1247, 6}:    "visionpyramid",        // Missing description for visionpyramid
	Port{1247, 17}:   "visionpyramid",        // VisionPyramid
	Port{1248, 6}:    "hermes",               // Missing description for hermes
	Port{1248, 17}:   "hermes",               // Missing description for hermes
	Port{1249, 6}:    "mesavistaco",          // Mesa Vista Co
	Port{1249, 17}:   "mesavistaco",          // Mesa Vista Co
	Port{1250, 6}:    "swldy-sias",           // Missing description for swldy-sias
	Port{1250, 17}:   "swldy-sias",           // Missing description for swldy-sias
	Port{1251, 6}:    "servergraph",          // Missing description for servergraph
	Port{1251, 17}:   "servergraph",          // Missing description for servergraph
	Port{1252, 6}:    "bspne-pcc",            // Missing description for bspne-pcc
	Port{1252, 17}:   "bspne-pcc",            // Missing description for bspne-pcc
	Port{1253, 6}:    "q55-pcc",              // Missing description for q55-pcc
	Port{1253, 17}:   "q55-pcc",              // Missing description for q55-pcc
	Port{1254, 6}:    "de-noc",               // Missing description for de-noc
	Port{1254, 17}:   "de-noc",               // Missing description for de-noc
	Port{1255, 6}:    "de-cache-query",       // Missing description for de-cache-query
	Port{1255, 17}:   "de-cache-query",       // Missing description for de-cache-query
	Port{1256, 6}:    "de-server",            // Missing description for de-server
	Port{1256, 17}:   "de-server",            // Missing description for de-server
	Port{1257, 6}:    "shockwave2",           // Shockwave 2
	Port{1257, 17}:   "shockwave2",           // Shockwave 2
	Port{1258, 6}:    "opennl",               // Open Network Library
	Port{1258, 17}:   "opennl",               // Open Network Library
	Port{1259, 6}:    "opennl-voice",         // Open Network Library Voice
	Port{1259, 17}:   "opennl-voice",         // Open Network Library Voice
	Port{1260, 6}:    "ibm-ssd",              // Missing description for ibm-ssd
	Port{1260, 17}:   "ibm-ssd",              // Missing description for ibm-ssd
	Port{1261, 6}:    "mpshrsv",              // Missing description for mpshrsv
	Port{1261, 17}:   "mpshrsv",              // Missing description for mpshrsv
	Port{1262, 6}:    "qnts-orb",             // Missing description for qnts-orb
	Port{1262, 17}:   "qnts-orb",             // QNTS-ORB
	Port{1263, 6}:    "dka",                  // Missing description for dka
	Port{1263, 17}:   "dka",                  // Missing description for dka
	Port{1264, 6}:    "prat",                 // Missing description for prat
	Port{1264, 17}:   "prat",                 // PRAT
	Port{1265, 6}:    "dssiapi",              // Missing description for dssiapi
	Port{1265, 17}:   "dssiapi",              // DSSIAPI
	Port{1266, 6}:    "dellpwrappks",         // Missing description for dellpwrappks
	Port{1266, 17}:   "dellpwrappks",         // DELLPWRAPPKS
	Port{1267, 6}:    "epc",                  // eTrust Policy Compliance
	Port{1267, 17}:   "epc",                  // eTrust Policy Compliance
	Port{1268, 6}:    "propel-msgsys",        // Missing description for propel-msgsys
	Port{1268, 17}:   "propel-msgsys",        // PROPEL-MSGSYS
	Port{1269, 6}:    "watilapp",             // Missing description for watilapp
	Port{1269, 17}:   "watilapp",             // WATiLaPP
	Port{1270, 6}:    "ssserver",             // opsmgr | Sun StorEdge Configuration Service | Microsoft Operations Manager
	Port{1270, 17}:   "opsmgr",               // Microsoft Operations Manager
	Port{1271, 6}:    "excw",                 // Missing description for excw
	Port{1271, 17}:   "excw",                 // eXcW
	Port{1272, 6}:    "cspmlockmgr",          // Missing description for cspmlockmgr
	Port{1272, 17}:   "cspmlockmgr",          // CSPMLockMgr
	Port{1273, 6}:    "emc-gateway",          // Missing description for emc-gateway
	Port{1273, 17}:   "emc-gateway",          // EMC-Gateway
	Port{1274, 6}:    "t1distproc",           // Missing description for t1distproc
	Port{1274, 17}:   "t1distproc",           // Missing description for t1distproc
	Port{1275, 6}:    "ivcollector",          // Missing description for ivcollector
	Port{1275, 17}:   "ivcollector",          // Missing description for ivcollector
	Port{1276, 6}:    "ivmanager",            // Missing description for ivmanager
	Port{1276, 17}:   "ivmanager",            // Missing description for ivmanager
	Port{1277, 6}:    "miva-mqs",             // mqs
	Port{1277, 17}:   "miva-mqs",             // mqs
	Port{1278, 6}:    "dellwebadmin-1",       // Dell Web Admin 1
	Port{1278, 17}:   "dellwebadmin-1",       // Dell Web Admin 1
	Port{1279, 6}:    "dellwebadmin-2",       // Dell Web Admin 2
	Port{1279, 17}:   "dellwebadmin-2",       // Dell Web Admin 2
	Port{1280, 6}:    "pictrography",         // Missing description for pictrography
	Port{1280, 17}:   "pictrography",         // Pictrography
	Port{1281, 6}:    "healthd",              // Missing description for healthd
	Port{1281, 17}:   "healthd",              // Missing description for healthd
	Port{1282, 6}:    "emperion",             // Missing description for emperion
	Port{1282, 17}:   "emperion",             // Emperion
	Port{1283, 6}:    "productinfo",          // Product Information
	Port{1283, 17}:   "productinfo",          // Product Information
	Port{1284, 6}:    "iee-qfx",              // Missing description for iee-qfx
	Port{1284, 17}:   "iee-qfx",              // IEE-QFX
	Port{1285, 6}:    "neoiface",             // Missing description for neoiface
	Port{1285, 17}:   "neoiface",             // Missing description for neoiface
	Port{1286, 6}:    "netuitive",            // Missing description for netuitive
	Port{1286, 17}:   "netuitive",            // Missing description for netuitive
	Port{1287, 6}:    "routematch",           // RouteMatch Com
	Port{1287, 17}:   "routematch",           // RouteMatch Com
	Port{1288, 6}:    "navbuddy",             // Missing description for navbuddy
	Port{1288, 17}:   "navbuddy",             // NavBuddy
	Port{1289, 6}:    "jwalkserver",          // Missing description for jwalkserver
	Port{1289, 17}:   "jwalkserver",          // JWalkServer
	Port{1290, 6}:    "winjaserver",          // Missing description for winjaserver
	Port{1290, 17}:   "winjaserver",          // WinJaServer
	Port{1291, 6}:    "seagulllms",           // Missing description for seagulllms
	Port{1291, 17}:   "seagulllms",           // SEAGULLLMS
	Port{1292, 6}:    "dsdn",                 // Missing description for dsdn
	Port{1292, 17}:   "dsdn",                 // Missing description for dsdn
	Port{1293, 6}:    "pkt-krb-ipsec",        // Missing description for pkt-krb-ipsec
	Port{1293, 17}:   "pkt-krb-ipsec",        // PKT-KRB-IPSec
	Port{1294, 6}:    "cmmdriver",            // Missing description for cmmdriver
	Port{1294, 17}:   "cmmdriver",            // CMMdriver
	Port{1295, 6}:    "ehtp",                 // End-by-Hop Transmission Protocol
	Port{1295, 17}:   "ehtp",                 // End-by-Hop Transmission Protocol
	Port{1296, 6}:    "dproxy",               // Missing description for dproxy
	Port{1296, 17}:   "dproxy",               // Missing description for dproxy
	Port{1297, 6}:    "sdproxy",              // Missing description for sdproxy
	Port{1297, 17}:   "sdproxy",              // Missing description for sdproxy
	Port{1298, 6}:    "lpcp",                 // Missing description for lpcp
	Port{1298, 17}:   "lpcp",                 // Missing description for lpcp
	Port{1299, 6}:    "hp-sci",               // Missing description for hp-sci
	Port{1299, 17}:   "hp-sci",               // Missing description for hp-sci
	Port{1300, 6}:    "h323hostcallsc",       // H323 Host Call Secure | H.323 Secure Call Control Signalling
	Port{1300, 17}:   "h323hostcallsc",       // H323 Host Call Secure
	Port{1301, 6}:    "ci3-software-1",       // Missing description for ci3-software-1
	Port{1301, 17}:   "ci3-software-1",       // CI3-Software-1
	Port{1302, 6}:    "ci3-software-2",       // Missing description for ci3-software-2
	Port{1302, 17}:   "ci3-software-2",       // CI3-Software-2
	Port{1303, 6}:    "sftsrv",               // Missing description for sftsrv
	Port{1303, 17}:   "sftsrv",               // Missing description for sftsrv
	Port{1304, 6}:    "boomerang",            // Missing description for boomerang
	Port{1304, 17}:   "boomerang",            // Boomerang
	Port{1305, 6}:    "pe-mike",              // Missing description for pe-mike
	Port{1305, 17}:   "pe-mike",              // Missing description for pe-mike
	Port{1306, 6}:    "re-conn-proto",        // Missing description for re-conn-proto
	Port{1306, 17}:   "re-conn-proto",        // RE-Conn-Proto
	Port{1307, 6}:    "pacmand",              // Missing description for pacmand
	Port{1307, 17}:   "pacmand",              // Pacmand
	Port{1308, 6}:    "odsi",                 // Optical Domain Service Interconnect (ODSI)
	Port{1308, 17}:   "odsi",                 // Optical Domain Service Interconnect (ODSI)
	Port{1309, 6}:    "jtag-server",          // JTAG server
	Port{1309, 17}:   "jtag-server",          // JTAG server
	Port{1310, 6}:    "husky",                // Missing description for husky
	Port{1310, 17}:   "husky",                // Husky
	Port{1311, 6}:    "rxmon",                // Missing description for rxmon
	Port{1311, 17}:   "rxmon",                // Missing description for rxmon
	Port{1312, 6}:    "sti-envision",         // STI Envision
	Port{1312, 17}:   "sti-envision",         // STI Envision
	Port{1313, 6}:    "bmc_patroldb",         // bmc-patroldb
	Port{1313, 17}:   "bmc_patroldb",         // BMC_PATROLDB
	Port{1314, 6}:    "pdps",                 // Photoscript Distributed Printing System
	Port{1314, 17}:   "pdps",                 // Photoscript Distributed Printing System
	Port{1315, 6}:    "els",                  // E.L.S., Event Listener Service
	Port{1315, 17}:   "els",                  // E.L.S., Event Listener Service
	Port{1316, 6}:    "exbit-escp",           // Missing description for exbit-escp
	Port{1316, 17}:   "exbit-escp",           // Exbit-ESCP
	Port{1317, 6}:    "vrts-ipcserver",       // Missing description for vrts-ipcserver
	Port{1317, 17}:   "vrts-ipcserver",       // Missing description for vrts-ipcserver
	Port{1318, 6}:    "krb5gatekeeper",       // Missing description for krb5gatekeeper
	Port{1318, 17}:   "krb5gatekeeper",       // Missing description for krb5gatekeeper
	Port{1319, 6}:    "amx-icsp",             // Missing description for amx-icsp
	Port{1319, 17}:   "amx-icsp",             // AMX-ICSP
	Port{1320, 6}:    "amx-axbnet",           // Missing description for amx-axbnet
	Port{1320, 17}:   "amx-axbnet",           // AMX-AXBNET
	Port{1321, 6}:    "pip",                  // Missing description for pip
	Port{1321, 17}:   "pip",                  // PIP
	Port{1322, 6}:    "novation",             // Missing description for novation
	Port{1322, 17}:   "novation",             // Novation
	Port{1323, 6}:    "brcd",                 // Missing description for brcd
	Port{1323, 17}:   "brcd",                 // Missing description for brcd
	Port{1324, 6}:    "delta-mcp",            // Missing description for delta-mcp
	Port{1324, 17}:   "delta-mcp",            // Missing description for delta-mcp
	Port{1325, 6}:    "dx-instrument",        // Missing description for dx-instrument
	Port{1325, 17}:   "dx-instrument",        // DX-Instrument
	Port{1326, 6}:    "wimsic",               // Missing description for wimsic
	Port{1326, 17}:   "wimsic",               // WIMSIC
	Port{1327, 6}:    "ultrex",               // Missing description for ultrex
	Port{1327, 17}:   "ultrex",               // Ultrex
	Port{1328, 6}:    "ewall",                // Missing description for ewall
	Port{1328, 17}:   "ewall",                // EWALL
	Port{1329, 6}:    "netdb-export",         // Missing description for netdb-export
	Port{1329, 17}:   "netdb-export",         // Missing description for netdb-export
	Port{1330, 6}:    "streetperfect",        // Missing description for streetperfect
	Port{1330, 17}:   "streetperfect",        // StreetPerfect
	Port{1331, 6}:    "intersan",             // Missing description for intersan
	Port{1331, 17}:   "intersan",             // Missing description for intersan
	Port{1332, 6}:    "pcia-rxp-b",           // PCIA RXP-B
	Port{1332, 17}:   "pcia-rxp-b",           // PCIA RXP-B
	Port{1333, 6}:    "passwrd-policy",       // Password Policy
	Port{1333, 17}:   "passwrd-policy",       // Password Policy
	Port{1334, 6}:    "writesrv",             // Missing description for writesrv
	Port{1334, 17}:   "writesrv",             // Missing description for writesrv
	Port{1335, 6}:    "digital-notary",       // Digital Notary Protocol
	Port{1335, 17}:   "digital-notary",       // Digital Notary Protocol
	Port{1336, 6}:    "ischat",               // Instant Service Chat
	Port{1336, 17}:   "ischat",               // Instant Service Chat
	Port{1337, 6}:    "waste",                // menandmice-dns | Nullsoft WASTE encrypted P2P app | menandmice DNS
	Port{1337, 17}:   "menandmice-dns",       // menandmice DNS
	Port{1338, 6}:    "wmc-log-svc",          // WMC-log-svr
	Port{1338, 17}:   "wmc-log-svc",          // WMC-log-svr
	Port{1339, 6}:    "kjtsiteserver",        // Missing description for kjtsiteserver
	Port{1339, 17}:   "kjtsiteserver",        // Missing description for kjtsiteserver
	Port{1340, 6}:    "naap",                 // Missing description for naap
	Port{1340, 17}:   "naap",                 // NAAP
	Port{1341, 6}:    "qubes",                // Missing description for qubes
	Port{1341, 17}:   "qubes",                // QuBES
	Port{1342, 6}:    "esbroker",             // Missing description for esbroker
	Port{1342, 17}:   "esbroker",             // ESBroker
	Port{1343, 6}:    "re101",                // Missing description for re101
	Port{1343, 17}:   "re101",                // Missing description for re101
	Port{1344, 6}:    "icap",                 // Missing description for icap
	Port{1344, 17}:   "icap",                 // ICAP
	Port{1345, 6}:    "vpjp",                 // Missing description for vpjp
	Port{1345, 17}:   "vpjp",                 // VPJP
	Port{1346, 6}:    "alta-ana-lm",          // Alta Analytics License Manager
	Port{1346, 17}:   "alta-ana-lm",          // Alta Analytics License Manager
	Port{1347, 6}:    "bbn-mmc",              // multi media conferencing
	Port{1347, 17}:   "bbn-mmc",              // multi media conferencing
	Port{1348, 6}:    "bbn-mmx",              // multi media conferencing
	Port{1348, 17}:   "bbn-mmx",              // multi media conferencing
	Port{1349, 6}:    "sbook",                // Registration Network Protocol
	Port{1349, 17}:   "sbook",                // Registration Network Protocol
	Port{1350, 6}:    "editbench",            // Registration Network Protocol
	Port{1350, 17}:   "editbench",            // Registration Network Protocol
	Port{1351, 6}:    "equationbuilder",      // Digital Tool Works (MIT)
	Port{1351, 17}:   "equationbuilder",      // Digital Tool Works (MIT)
	Port{1352, 6}:    "lotusnotes",           // lotusnote | Lotus Note
	Port{1352, 17}:   "lotusnotes",           // Lotus Note
	Port{1353, 6}:    "relief",               // Relief Consulting
	Port{1353, 17}:   "relief",               // Relief Consulting
	Port{1354, 6}:    "rightbrain",           // XSIP-network | RightBrain Software | Five Across XSIP Network
	Port{1354, 17}:   "rightbrain",           // RightBrain Software
	Port{1355, 6}:    "intuitive-edge",       // Intuitive Edge
	Port{1355, 17}:   "intuitive-edge",       // Intuitive Edge
	Port{1356, 6}:    "cuillamartin",         // CuillaMartin Company
	Port{1356, 17}:   "cuillamartin",         // CuillaMartin Company
	Port{1357, 6}:    "pegboard",             // Electronic PegBoard
	Port{1357, 17}:   "pegboard",             // Electronic PegBoard
	Port{1358, 6}:    "connlcli",             // Missing description for connlcli
	Port{1358, 17}:   "connlcli",             // Missing description for connlcli
	Port{1359, 6}:    "ftsrv",                // Missing description for ftsrv
	Port{1359, 17}:   "ftsrv",                // Missing description for ftsrv
	Port{1360, 6}:    "mimer",                // Missing description for mimer
	Port{1360, 17}:   "mimer",                // Missing description for mimer
	Port{1361, 6}:    "linx",                 // Missing description for linx
	Port{1361, 17}:   "linx",                 // Missing description for linx
	Port{1362, 6}:    "timeflies",            // Missing description for timeflies
	Port{1362, 17}:   "timeflies",            // Missing description for timeflies
	Port{1363, 6}:    "ndm-requester",        // Network DataMover Requester
	Port{1363, 17}:   "ndm-requester",        // Network DataMover Requester
	Port{1364, 6}:    "ndm-server",           // Network DataMover Server
	Port{1364, 17}:   "ndm-server",           // Network DataMover Server
	Port{1365, 6}:    "adapt-sna",            // Network Software Associates
	Port{1365, 17}:   "adapt-sna",            // Network Software Associates
	Port{1366, 6}:    "netware-csp",          // Novell NetWare Comm Service Platform
	Port{1366, 17}:   "netware-csp",          // Novell NetWare Comm Service Platform
	Port{1367, 6}:    "dcs",                  // Missing description for dcs
	Port{1367, 17}:   "dcs",                  // Missing description for dcs
	Port{1368, 6}:    "screencast",           // Missing description for screencast
	Port{1368, 17}:   "screencast",           // Missing description for screencast
	Port{1369, 6}:    "gv-us",                // GlobalView to Unix Shell
	Port{1369, 17}:   "gv-us",                // GlobalView to Unix Shell
	Port{1370, 6}:    "us-gv",                // Unix Shell to GlobalView
	Port{1370, 17}:   "us-gv",                // Unix Shell to GlobalView
	Port{1371, 6}:    "fc-cli",               // Fujitsu Config Protocol
	Port{1371, 17}:   "fc-cli",               // Fujitsu Config Protocol
	Port{1372, 6}:    "fc-ser",               // Fujitsu Config Protocol
	Port{1372, 17}:   "fc-ser",               // Fujitsu Config Protocol
	Port{1373, 6}:    "chromagrafx",          // Missing description for chromagrafx
	Port{1373, 17}:   "chromagrafx",          // Missing description for chromagrafx
	Port{1374, 6}:    "molly",                // EPI Software Systems
	Port{1374, 17}:   "molly",                // EPI Software Systems
	Port{1375, 6}:    "bytex",                // Missing description for bytex
	Port{1375, 17}:   "bytex",                // Missing description for bytex
	Port{1376, 6}:    "ibm-pps",              // IBM Person to Person Software
	Port{1376, 17}:   "ibm-pps",              // IBM Person to Person Software
	Port{1377, 6}:    "cichlid",              // Cichlid License Manager
	Port{1377, 17}:   "cichlid",              // Cichlid License Manager
	Port{1378, 6}:    "elan",                 // Elan License Manager
	Port{1378, 17}:   "elan",                 // Elan License Manager
	Port{1379, 6}:    "dbreporter",           // Integrity Solutions
	Port{1379, 17}:   "dbreporter",           // Integrity Solutions
	Port{1380, 6}:    "telesis-licman",       // Telesis Network License Manager
	Port{1380, 17}:   "telesis-licman",       // Telesis Network License Manager
	Port{1381, 6}:    "apple-licman",         // Apple Network License Manager
	Port{1381, 17}:   "apple-licman",         // Apple Network License Manager
	Port{1382, 6}:    "udt_os",               // udt-os
	Port{1382, 17}:   "udt_os",               // Missing description for udt_os
	Port{1383, 6}:    "gwha",                 // GW Hannaway Network License Manager
	Port{1383, 17}:   "gwha",                 // GW Hannaway Network License Manager
	Port{1384, 6}:    "os-licman",            // Objective Solutions License Manager
	Port{1384, 17}:   "os-licman",            // Objective Solutions License Manager
	Port{1385, 6}:    "atex_elmd",            // atex-elmd | Atex Publishing License Manager
	Port{1385, 17}:   "atex_elmd",            // Atex Publishing License Manager
	Port{1386, 6}:    "checksum",             // CheckSum License Manager
	Port{1386, 17}:   "checksum",             // CheckSum License Manager
	Port{1387, 6}:    "cadsi-lm",             // Computer Aided Design Software Inc LM
	Port{1387, 17}:   "cadsi-lm",             // Computer Aided Design Software Inc LM
	Port{1388, 6}:    "objective-dbc",        // Objective Solutions DataBase Cache
	Port{1388, 17}:   "objective-dbc",        // Objective Solutions DataBase Cache
	Port{1389, 6}:    "iclpv-dm",             // Document Manager
	Port{1389, 17}:   "iclpv-dm",             // Document Manager
	Port{1390, 6}:    "iclpv-sc",             // Storage Controller
	Port{1390, 17}:   "iclpv-sc",             // Storage Controller
	Port{1391, 6}:    "iclpv-sas",            // Storage Access Server
	Port{1391, 17}:   "iclpv-sas",            // Storage Access Server
	Port{1392, 6}:    "iclpv-pm",             // Print Manager
	Port{1392, 17}:   "iclpv-pm",             // Print Manager
	Port{1393, 6}:    "iclpv-nls",            // Network Log Server
	Port{1393, 17}:   "iclpv-nls",            // Network Log Server
	Port{1394, 6}:    "iclpv-nlc",            // Network Log Client
	Port{1394, 17}:   "iclpv-nlc",            // Network Log Client
	Port{1395, 6}:    "iclpv-wsm",            // PC Workstation Manager software
	Port{1395, 17}:   "iclpv-wsm",            // PC Workstation Manager software
	Port{1396, 6}:    "dvl-activemail",       // DVL Active Mail
	Port{1396, 17}:   "dvl-activemail",       // DVL Active Mail
	Port{1397, 6}:    "audio-activmail",      // Audio Active Mail
	Port{1397, 17}:   "audio-activmail",      // Audio Active Mail
	Port{1398, 6}:    "video-activmail",      // Video Active Mail
	Port{1398, 17}:   "video-activmail",      // Video Active Mail
	Port{1399, 6}:    "cadkey-licman",        // Cadkey License Manager
	Port{1399, 17}:   "cadkey-licman",        // Cadkey License Manager
	Port{1400, 6}:    "cadkey-tablet",        // Cadkey Tablet Daemon
	Port{1400, 17}:   "cadkey-tablet",        // Cadkey Tablet Daemon
	Port{1401, 6}:    "goldleaf-licman",      // Goldleaf License Manager
	Port{1401, 17}:   "goldleaf-licman",      // Goldleaf License Manager
	Port{1402, 6}:    "prm-sm-np",            // Prospero Resource Manager
	Port{1402, 17}:   "prm-sm-np",            // Prospero Resource Manager
	Port{1403, 6}:    "prm-nm-np",            // Prospero Resource Manager
	Port{1403, 17}:   "prm-nm-np",            // Prospero Resource Manager
	Port{1404, 6}:    "igi-lm",               // Infinite Graphics License Manager
	Port{1404, 17}:   "igi-lm",               // Infinite Graphics License Manager
	Port{1405, 6}:    "ibm-res",              // IBM Remote Execution Starter
	Port{1405, 17}:   "ibm-res",              // IBM Remote Execution Starter
	Port{1406, 6}:    "netlabs-lm",           // NetLabs License Manager
	Port{1406, 17}:   "netlabs-lm",           // NetLabs License Manager
	Port{1407, 6}:    "dbsa-lm",              // tibet-server | DBSA License Manager | TIBET Data Server
	Port{1407, 17}:   "dbsa-lm",              // DBSA License Manager
	Port{1408, 6}:    "sophia-lm",            // Sophia License Manager
	Port{1408, 17}:   "sophia-lm",            // Sophia License Manager
	Port{1409, 6}:    "here-lm",              // Here License Manager
	Port{1409, 17}:   "here-lm",              // Here License Manager
	Port{1410, 6}:    "hiq",                  // HiQ License Manager
	Port{1410, 17}:   "hiq",                  // HiQ License Manager
	Port{1411, 6}:    "af",                   // AudioFile
	Port{1411, 17}:   "af",                   // AudioFile
	Port{1412, 6}:    "innosys",              // Missing description for innosys
	Port{1412, 17}:   "innosys",              // Missing description for innosys
	Port{1413, 6}:    "innosys-acl",          // Missing description for innosys-acl
	Port{1413, 17}:   "innosys-acl",          // Missing description for innosys-acl
	Port{1414, 6}:    "ibm-mqseries",         // IBM MQSeries
	Port{1414, 17}:   "ibm-mqseries",         // IBM MQSeries
	Port{1415, 6}:    "dbstar",               // Missing description for dbstar
	Port{1415, 17}:   "dbstar",               // Missing description for dbstar
	Port{1416, 6}:    "novell-lu6.2",         // novell-lu6-2 | Novell LU6.2
	Port{1416, 17}:   "novell-lu6.2",         // Novell LU6.2
	Port{1417, 6}:    "timbuktu-srv1",        // Timbuktu Service 1 Port
	Port{1417, 17}:   "timbuktu-srv1",        // Timbuktu Service 1 Port
	Port{1418, 6}:    "timbuktu-srv2",        // Timbuktu Service 2 Port
	Port{1418, 17}:   "timbuktu-srv2",        // Timbuktu Service 2 Port
	Port{1419, 6}:    "timbuktu-srv3",        // Timbuktu Service 3 Port
	Port{1419, 17}:   "timbuktu-srv3",        // Timbuktu Service 3 Port
	Port{1420, 6}:    "timbuktu-srv4",        // Timbuktu Service 4 Port
	Port{1420, 17}:   "timbuktu-srv4",        // Timbuktu Service 4 Port
	Port{1421, 6}:    "gandalf-lm",           // Gandalf License Manager
	Port{1421, 17}:   "gandalf-lm",           // Gandalf License Manager
	Port{1422, 6}:    "autodesk-lm",          // Autodesk License Manager
	Port{1422, 17}:   "autodesk-lm",          // Autodesk License Manager
	Port{1423, 6}:    "essbase",              // Essbase Arbor Software
	Port{1423, 17}:   "essbase",              // Essbase Arbor Software
	Port{1424, 6}:    "hybrid",               // Hybrid Encryption Protocol
	Port{1424, 17}:   "hybrid",               // Hybrid Encryption Protocol
	Port{1425, 6}:    "zion-lm",              // Zion Software License Manager
	Port{1425, 17}:   "zion-lm",              // Zion Software License Manager
	Port{1426, 6}:    "sas-1",                // sais | Satellite-data Acquisition System 1
	Port{1426, 17}:   "sas-1",                // Satellite-data Acquisition System 1
	Port{1427, 6}:    "mloadd",               // mloadd monitoring tool
	Port{1427, 17}:   "mloadd",               // mloadd monitoring tool
	Port{1428, 6}:    "informatik-lm",        // Informatik License Manager
	Port{1428, 17}:   "informatik-lm",        // Informatik License Manager
	Port{1429, 6}:    "nms",                  // Hypercom NMS
	Port{1429, 17}:   "nms",                  // Hypercom NMS
	Port{1430, 6}:    "tpdu",                 // Hypercom TPDU
	Port{1430, 17}:   "tpdu",                 // Hypercom TPDU
	Port{1431, 6}:    "rgtp",                 // Reverse Gossip Transport
	Port{1431, 17}:   "rgtp",                 // Reverse Gossip Transport
	Port{1432, 6}:    "blueberry-lm",         // Blueberry Software License Manager
	Port{1432, 17}:   "blueberry-lm",         // Blueberry Software License Manager
	Port{1433, 6}:    "ms-sql-s",             // Microsoft-SQL-Server
	Port{1433, 17}:   "ms-sql-s",             // Microsoft-SQL-Server
	Port{1434, 6}:    "ms-sql-m",             // Microsoft-SQL-Monitor
	Port{1434, 17}:   "ms-sql-m",             // Microsoft-SQL-Monitor
	Port{1435, 6}:    "ibm-cics",             // IBM CICS
	Port{1435, 17}:   "ibm-cics",             // Missing description for ibm-cics
	Port{1436, 6}:    "sas-2",                // saism | Satellite-data Acquisition System 2
	Port{1436, 17}:   "sas-2",                // Satellite-data Acquisition System 2
	Port{1437, 6}:    "tabula",               // Missing description for tabula
	Port{1437, 17}:   "tabula",               // Missing description for tabula
	Port{1438, 6}:    "eicon-server",         // Eicon Security Agent Server
	Port{1438, 17}:   "eicon-server",         // Eicon Security Agent Server
	Port{1439, 6}:    "eicon-x25",            // Eicon X25 SNA Gateway
	Port{1439, 17}:   "eicon-x25",            // Eicon X25 SNA Gateway
	Port{1440, 6}:    "eicon-slp",            // Eicon Service Location Protocol
	Port{1440, 17}:   "eicon-slp",            // Eicon Service Location Protocol
	Port{1441, 6}:    "cadis-1",              // Cadis License Management
	Port{1441, 17}:   "cadis-1",              // Cadis License Management
	Port{1442, 6}:    "cadis-2",              // Cadis License Management
	Port{1442, 17}:   "cadis-2",              // Cadis License Management
	Port{1443, 6}:    "ies-lm",               // Integrated Engineering Software
	Port{1443, 17}:   "ies-lm",               // Integrated Engineering Software
	Port{1444, 6}:    "marcam-lm",            // Marcam License Management
	Port{1444, 17}:   "marcam-lm",            // Marcam License Management
	Port{1445, 6}:    "proxima-lm",           // Proxima License Manager
	Port{1445, 17}:   "proxima-lm",           // Proxima License Manager
	Port{1446, 6}:    "ora-lm",               // Optical Research Associates License Manager
	Port{1446, 17}:   "ora-lm",               // Optical Research Associates License Manager
	Port{1447, 6}:    "apri-lm",              // Applied Parallel Research LM
	Port{1447, 17}:   "apri-lm",              // Applied Parallel Research LM
	Port{1448, 6}:    "oc-lm",                // OpenConnect License Manager
	Port{1448, 17}:   "oc-lm",                // OpenConnect License Manager
	Port{1449, 6}:    "peport",               // Missing description for peport
	Port{1449, 17}:   "peport",               // Missing description for peport
	Port{1450, 6}:    "dwf",                  // Tandem Distributed Workbench Facility
	Port{1450, 17}:   "dwf",                  // Tandem Distributed Workbench Facility
	Port{1451, 6}:    "infoman",              // IBM Information Management
	Port{1451, 17}:   "infoman",              // IBM Information Management
	Port{1452, 6}:    "gtegsc-lm",            // GTE Government Systems License Man
	Port{1452, 17}:   "gtegsc-lm",            // GTE Government Systems License Man
	Port{1453, 6}:    "genie-lm",             // Genie License Manager
	Port{1453, 17}:   "genie-lm",             // Genie License Manager
	Port{1454, 6}:    "interhdl_elmd",        // interhdl-elmd | interHDL License Manager
	Port{1454, 17}:   "interhdl_elmd",        // interHDL License Manager
	Port{1455, 6}:    "esl-lm",               // ESL License Manager
	Port{1455, 17}:   "esl-lm",               // ESL License Manager
	Port{1456, 6}:    "dca",                  // Missing description for dca
	Port{1456, 17}:   "dca",                  // Missing description for dca
	Port{1457, 6}:    "valisys-lm",           // Valisys License Manager
	Port{1457, 17}:   "valisys-lm",           // Valisys License Manager
	Port{1458, 6}:    "nrcabq-lm",            // Nichols Research Corp.
	Port{1458, 17}:   "nrcabq-lm",            // Nichols Research Corp.
	Port{1459, 6}:    "proshare1",            // Proshare Notebook Application
	Port{1459, 17}:   "proshare1",            // Proshare Notebook Application
	Port{1460, 6}:    "proshare2",            // Proshare Notebook Application
	Port{1460, 17}:   "proshare2",            // Proshare Notebook Application
	Port{1461, 6}:    "ibm_wrless_lan",       // ibm-wrless-lan | IBM Wireless LAN
	Port{1461, 17}:   "ibm_wrless_lan",       // IBM Wireless LAN
	Port{1462, 6}:    "world-lm",             // World License Manager
	Port{1462, 17}:   "world-lm",             // World License Manager
	Port{1463, 6}:    "nucleus",              // Missing description for nucleus
	Port{1463, 17}:   "nucleus",              // Missing description for nucleus
	Port{1464, 6}:    "msl_lmd",              // msl-lmd | MSL License Manager
	Port{1464, 17}:   "msl_lmd",              // MSL License Manager
	Port{1465, 6}:    "pipes",                // Pipes Platform
	Port{1465, 17}:   "pipes",                // Missing description for pipes
	Port{1466, 6}:    "oceansoft-lm",         // Ocean Software License Manager
	Port{1466, 17}:   "oceansoft-lm",         // Ocean Software License Manager
	Port{1467, 6}:    "csdmbase",             // Missing description for csdmbase
	Port{1467, 17}:   "csdmbase",             // Missing description for csdmbase
	Port{1468, 6}:    "csdm",                 // Missing description for csdm
	Port{1468, 17}:   "csdm",                 // Missing description for csdm
	Port{1469, 6}:    "aal-lm",               // Active Analysis Limited License Manager
	Port{1469, 17}:   "aal-lm",               // Active Analysis Limited License Manager
	Port{1470, 6}:    "uaiact",               // Universal Analytics
	Port{1470, 17}:   "uaiact",               // Universal Analytics
	Port{1471, 6}:    "csdmbase",             // Missing description for csdmbase
	Port{1471, 17}:   "csdmbase",             // Missing description for csdmbase
	Port{1472, 6}:    "csdm",                 // Missing description for csdm
	Port{1472, 17}:   "csdm",                 // Missing description for csdm
	Port{1473, 6}:    "openmath",             // Missing description for openmath
	Port{1473, 17}:   "openmath",             // Missing description for openmath
	Port{1474, 6}:    "telefinder",           // Missing description for telefinder
	Port{1474, 17}:   "telefinder",           // Missing description for telefinder
	Port{1475, 6}:    "taligent-lm",          // Taligent License Manager
	Port{1475, 17}:   "taligent-lm",          // Taligent License Manager
	Port{1476, 6}:    "clvm-cfg",             // Missing description for clvm-cfg
	Port{1476, 17}:   "clvm-cfg",             // Missing description for clvm-cfg
	Port{1477, 6}:    "ms-sna-server",        // Missing description for ms-sna-server
	Port{1477, 17}:   "ms-sna-server",        // Missing description for ms-sna-server
	Port{1478, 6}:    "ms-sna-base",          // Missing description for ms-sna-base
	Port{1478, 17}:   "ms-sna-base",          // Missing description for ms-sna-base
	Port{1479, 6}:    "dberegister",          // Missing description for dberegister
	Port{1479, 17}:   "dberegister",          // Missing description for dberegister
	Port{1480, 6}:    "pacerforum",           // Missing description for pacerforum
	Port{1480, 17}:   "pacerforum",           // Missing description for pacerforum
	Port{1481, 6}:    "airs",                 // Missing description for airs
	Port{1481, 17}:   "airs",                 // Missing description for airs
	Port{1482, 6}:    "miteksys-lm",          // Miteksys License Manager
	Port{1482, 17}:   "miteksys-lm",          // Miteksys License Manager
	Port{1483, 6}:    "afs",                  // AFS License Manager
	Port{1483, 17}:   "afs",                  // AFS License Manager
	Port{1484, 6}:    "confluent",            // Confluent License Manager
	Port{1484, 17}:   "confluent",            // Confluent License Manager
	Port{1485, 6}:    "lansource",            // Missing description for lansource
	Port{1485, 17}:   "lansource",            // Missing description for lansource
	Port{1486, 6}:    "nms_topo_serv",        // nms-topo-serv
	Port{1486, 17}:   "nms_topo_serv",        // Missing description for nms_topo_serv
	Port{1487, 6}:    "localinfosrvr",        // Missing description for localinfosrvr
	Port{1487, 17}:   "localinfosrvr",        // Missing description for localinfosrvr
	Port{1488, 6}:    "docstor",              // Missing description for docstor
	Port{1488, 17}:   "docstor",              // Missing description for docstor
	Port{1489, 6}:    "dmdocbroker",          // Missing description for dmdocbroker
	Port{1489, 17}:   "dmdocbroker",          // Missing description for dmdocbroker
	Port{1490, 6}:    "insitu-conf",          // Missing description for insitu-conf
	Port{1490, 17}:   "insitu-conf",          // Missing description for insitu-conf
	Port{1491, 6}:    "anynetgateway",        // Missing description for anynetgateway
	Port{1491, 17}:   "anynetgateway",        // Missing description for anynetgateway
	Port{1492, 6}:    "stone-design-1",       // Missing description for stone-design-1
	Port{1492, 17}:   "stone-design-1",       // Missing description for stone-design-1
	Port{1493, 6}:    "netmap_lm",            // netmap-lm
	Port{1493, 17}:   "netmap_lm",            // Missing description for netmap_lm
	Port{1494, 6}:    "citrix-ica",           // ica
	Port{1494, 17}:   "citrix-ica",           // Missing description for citrix-ica
	Port{1495, 6}:    "cvc",                  // Missing description for cvc
	Port{1495, 17}:   "cvc",                  // Missing description for cvc
	Port{1496, 6}:    "liberty-lm",           // Missing description for liberty-lm
	Port{1496, 17}:   "liberty-lm",           // Missing description for liberty-lm
	Port{1497, 6}:    "rfx-lm",               // Missing description for rfx-lm
	Port{1497, 17}:   "rfx-lm",               // Missing description for rfx-lm
	Port{1498, 6}:    "watcom-sql",           // sybase-sqlany | Sybase SQL Any
	Port{1498, 17}:   "watcom-sql",           // Missing description for watcom-sql
	Port{1499, 6}:    "fhc",                  // Federico Heinz Consultora
	Port{1499, 17}:   "fhc",                  // Federico Heinz Consultora
	Port{1500, 6}:    "vlsi-lm",              // VLSI License Manager
	Port{1500, 17}:   "vlsi-lm",              // VLSI License Manager
	Port{1501, 6}:    "sas-3",                // saiscm | Satellite-data Acquisition System 3
	Port{1501, 17}:   "sas-3",                // Satellite-data Acquisition System 3
	Port{1502, 6}:    "shivadiscovery",       // Shiva
	Port{1502, 17}:   "shivadiscovery",       // Shiva
	Port{1503, 6}:    "imtc-mcs",             // Databeam
	Port{1503, 17}:   "imtc-mcs",             // Databeam
	Port{1504, 6}:    "evb-elm",              // EVB Software Engineering License Manager
	Port{1504, 17}:   "evb-elm",              // EVB Software Engineering License Manager
	Port{1505, 6}:    "funkproxy",            // Funk Software, Inc.
	Port{1505, 17}:   "funkproxy",            // Funk Software, Inc.
	Port{1506, 6}:    "utcd",                 // Universal Time daemon (utcd)
	Port{1506, 17}:   "utcd",                 // Universal Time daemon (utcd)
	Port{1507, 6}:    "symplex",              // Missing description for symplex
	Port{1507, 17}:   "symplex",              // Missing description for symplex
	Port{1508, 6}:    "diagmond",             // Missing description for diagmond
	Port{1508, 17}:   "diagmond",             // Missing description for diagmond
	Port{1509, 6}:    "robcad-lm",            // Robcad, Ltd. License Manager
	Port{1509, 17}:   "robcad-lm",            // Robcad, Ltd. License Manager
	Port{1510, 6}:    "mvx-lm",               // Midland Valley Exploration Ltd. Lic. Man.
	Port{1510, 17}:   "mvx-lm",               // Midland Valley Exploration Ltd. Lic. Man.
	Port{1511, 6}:    "3l-l1",                // Missing description for 3l-l1
	Port{1511, 17}:   "3l-l1",                // Missing description for 3l-l1
	Port{1512, 6}:    "wins",                 // Microsoft's Windows Internet Name Service
	Port{1512, 17}:   "wins",                 // Microsoft's Windows Internet Name Service
	Port{1513, 6}:    "fujitsu-dtc",          // Fujitsu Systems Business of America, Inc
	Port{1513, 17}:   "fujitsu-dtc",          // Fujitsu Systems Business of America, Inc
	Port{1514, 6}:    "fujitsu-dtcns",        // Fujitsu Systems Business of America, Inc
	Port{1514, 17}:   "fujitsu-dtcns",        // Fujitsu Systems Business of America, Inc
	Port{1515, 6}:    "ifor-protocol",        // Missing description for ifor-protocol
	Port{1515, 17}:   "ifor-protocol",        // Missing description for ifor-protocol
	Port{1516, 6}:    "vpad",                 // Virtual Places Audio data
	Port{1516, 17}:   "vpad",                 // Virtual Places Audio data
	Port{1517, 6}:    "vpac",                 // Virtual Places Audio control
	Port{1517, 17}:   "vpac",                 // Virtual Places Audio control
	Port{1518, 6}:    "vpvd",                 // Virtual Places Video data
	Port{1518, 17}:   "vpvd",                 // Virtual Places Video data
	Port{1519, 6}:    "vpvc",                 // Virtual Places Video control
	Port{1519, 17}:   "vpvc",                 // Virtual Places Video control
	Port{1520, 6}:    "atm-zip-office",       // atm zip office
	Port{1520, 17}:   "atm-zip-office",       // atm zip office
	Port{1521, 6}:    "oracle",               // ncube-lm | Oracle Database | nCube License Manager
	Port{1521, 17}:   "ncube-lm",             // nCube License Manager
	Port{1522, 6}:    "rna-lm",               // ricardo-lm | Ricardo North America License Manager
	Port{1522, 17}:   "rna-lm",               // Ricardo North America License Manager
	Port{1523, 6}:    "cichild-lm",           // cichild
	Port{1523, 17}:   "cichild-lm",           // Missing description for cichild-lm
	Port{1524, 6}:    "ingreslock",           // ingres
	Port{1524, 17}:   "ingreslock",           // ingres
	Port{1525, 6}:    "orasrv",               // prospero-np | oracle or Prospero Directory Service non-priv | oracle | Prospero Directory Service non-priv
	Port{1525, 17}:   "oracle",               // Missing description for oracle
	Port{1526, 6}:    "pdap-np",              // Prospero Data Access Prot non-priv
	Port{1526, 17}:   "pdap-np",              // Prospero Data Access Prot non-priv
	Port{1527, 6}:    "tlisrv",               // oracle
	Port{1527, 17}:   "tlisrv",               // oracle
	Port{1528, 6}:    "mciautoreg",           // ngr-t | NGR transport protocol for mobile ad-hoc networks
	Port{1528, 17}:   "mciautoreg",           // Missing description for mciautoreg
	Port{1529, 6}:    "support",              // coauthor | prmsd gnatsd cygnus bug tracker | oracle
	Port{1529, 17}:   "coauthor",             // oracle
	Port{1530, 6}:    "rap-service",          // Missing description for rap-service
	Port{1530, 17}:   "rap-service",          // Missing description for rap-service
	Port{1531, 6}:    "rap-listen",           // Missing description for rap-listen
	Port{1531, 17}:   "rap-listen",           // Missing description for rap-listen
	Port{1532, 6}:    "miroconnect",          // Missing description for miroconnect
	Port{1532, 17}:   "miroconnect",          // Missing description for miroconnect
	Port{1533, 6}:    "virtual-places",       // Virtual Places Software
	Port{1533, 17}:   "virtual-places",       // Virtual Places Software
	Port{1534, 6}:    "micromuse-lm",         // Missing description for micromuse-lm
	Port{1534, 17}:   "micromuse-lm",         // Missing description for micromuse-lm
	Port{1535, 6}:    "ampr-info",            // Missing description for ampr-info
	Port{1535, 17}:   "ampr-info",            // Missing description for ampr-info
	Port{1536, 6}:    "ampr-inter",           // Missing description for ampr-inter
	Port{1536, 17}:   "ampr-inter",           // Missing description for ampr-inter
	Port{1537, 6}:    "sdsc-lm",              // isi-lm
	Port{1537, 17}:   "sdsc-lm",              // Missing description for sdsc-lm
	Port{1538, 6}:    "3ds-lm",               // Missing description for 3ds-lm
	Port{1538, 17}:   "3ds-lm",               // Missing description for 3ds-lm
	Port{1539, 6}:    "intellistor-lm",       // Intellistor License Manager
	Port{1539, 17}:   "intellistor-lm",       // Intellistor License Manager
	Port{1540, 6}:    "rds",                  // Missing description for rds
	Port{1540, 17}:   "rds",                  // Missing description for rds
	Port{1541, 6}:    "rds2",                 // Missing description for rds2
	Port{1541, 17}:   "rds2",                 // Missing description for rds2
	Port{1542, 6}:    "gridgen-elmd",         // Missing description for gridgen-elmd
	Port{1542, 17}:   "gridgen-elmd",         // Missing description for gridgen-elmd
	Port{1543, 6}:    "simba-cs",             // Missing description for simba-cs
	Port{1543, 17}:   "simba-cs",             // Missing description for simba-cs
	Port{1544, 6}:    "aspeclmd",             // Missing description for aspeclmd
	Port{1544, 17}:   "aspeclmd",             // Missing description for aspeclmd
	Port{1545, 6}:    "vistium-share",        // Missing description for vistium-share
	Port{1545, 17}:   "vistium-share",        // Missing description for vistium-share
	Port{1546, 6}:    "abbaccuray",           // Missing description for abbaccuray
	Port{1546, 17}:   "abbaccuray",           // Missing description for abbaccuray
	Port{1547, 6}:    "laplink",              // Missing description for laplink
	Port{1547, 17}:   "laplink",              // Missing description for laplink
	Port{1548, 6}:    "axon-lm",              // Axon License Manager
	Port{1548, 17}:   "axon-lm",              // Axon License Manager
	Port{1549, 6}:    "shivahose",            // shivasound | Shiva Hose | Shiva Sound
	Port{1549, 17}:   "shivasound",           // Shiva Sound
	Port{1550, 6}:    "3m-image-lm",          // Image Storage license manager 3M Company
	Port{1550, 17}:   "3m-image-lm",          // Image Storage license manager 3M Company
	Port{1551, 6}:    "hecmtl-db",            // Missing description for hecmtl-db
	Port{1551, 17}:   "hecmtl-db",            // Missing description for hecmtl-db
	Port{1552, 6}:    "pciarray",             // Missing description for pciarray
	Port{1552, 17}:   "pciarray",             // Missing description for pciarray
	Port{1553, 6}:    "sna-cs",               // Missing description for sna-cs
	Port{1553, 17}:   "sna-cs",               // Missing description for sna-cs
	Port{1554, 6}:    "caci-lm",              // CACI Products Company License Manager
	Port{1554, 17}:   "caci-lm",              // CACI Products Company License Manager
	Port{1555, 6}:    "livelan",              // Missing description for livelan
	Port{1555, 17}:   "livelan",              // Missing description for livelan
	Port{1556, 6}:    "veritas_pbx",          // veritas-pbx | VERITAS Private Branch Exchange
	Port{1556, 17}:   "veritas_pbx",          // VERITAS Private Branch Exchange
	Port{1557, 6}:    "arbortext-lm",         // ArborText License Manager
	Port{1557, 17}:   "arbortext-lm",         // ArborText License Manager
	Port{1558, 6}:    "xingmpeg",             // Missing description for xingmpeg
	Port{1558, 17}:   "xingmpeg",             // Missing description for xingmpeg
	Port{1559, 6}:    "web2host",             // Missing description for web2host
	Port{1559, 17}:   "web2host",             // Missing description for web2host
	Port{1560, 6}:    "asci-val",             // ASCI-RemoteSHADOW
	Port{1560, 17}:   "asci-val",             // ASCI-RemoteSHADOW
	Port{1561, 6}:    "facilityview",         // Missing description for facilityview
	Port{1561, 17}:   "facilityview",         // Missing description for facilityview
	Port{1562, 6}:    "pconnectmgr",          // Missing description for pconnectmgr
	Port{1562, 17}:   "pconnectmgr",          // Missing description for pconnectmgr
	Port{1563, 6}:    "cadabra-lm",           // Cadabra License Manager
	Port{1563, 17}:   "cadabra-lm",           // Cadabra License Manager
	Port{1564, 6}:    "pay-per-view",         // Missing description for pay-per-view
	Port{1564, 17}:   "pay-per-view",         // Pay-Per-View
	Port{1565, 6}:    "winddlb",              // WinDD
	Port{1565, 17}:   "winddlb",              // WinDD
	Port{1566, 6}:    "corelvideo",           // Missing description for corelvideo
	Port{1566, 17}:   "corelvideo",           // CORELVIDEO
	Port{1567, 6}:    "jlicelmd",             // Missing description for jlicelmd
	Port{1567, 17}:   "jlicelmd",             // Missing description for jlicelmd
	Port{1568, 6}:    "tsspmap",              // Missing description for tsspmap
	Port{1568, 17}:   "tsspmap",              // Missing description for tsspmap
	Port{1569, 6}:    "ets",                  // Missing description for ets
	Port{1569, 17}:   "ets",                  // Missing description for ets
	Port{1570, 6}:    "orbixd",               // Missing description for orbixd
	Port{1570, 17}:   "orbixd",               // Missing description for orbixd
	Port{1571, 6}:    "rdb-dbs-disp",         // Oracle Remote Data Base
	Port{1571, 17}:   "rdb-dbs-disp",         // Oracle Remote Data Base
	Port{1572, 6}:    "chip-lm",              // Chipcom License Manager
	Port{1572, 17}:   "chip-lm",              // Chipcom License Manager
	Port{1573, 6}:    "itscomm-ns",           // Missing description for itscomm-ns
	Port{1573, 17}:   "itscomm-ns",           // Missing description for itscomm-ns
	Port{1574, 6}:    "mvel-lm",              // Missing description for mvel-lm
	Port{1574, 17}:   "mvel-lm",              // Missing description for mvel-lm
	Port{1575, 6}:    "oraclenames",          // Missing description for oraclenames
	Port{1575, 17}:   "oraclenames",          // Missing description for oraclenames
	Port{1576, 6}:    "moldflow-lm",          // Moldflow License Manager
	Port{1576, 17}:   "moldflow-lm",          // Moldflow License Manager
	Port{1577, 6}:    "hypercube-lm",         // Missing description for hypercube-lm
	Port{1577, 17}:   "hypercube-lm",         // Missing description for hypercube-lm
	Port{1578, 6}:    "jacobus-lm",           // Jacobus License Manager
	Port{1578, 17}:   "jacobus-lm",           // Jacobus License Manager
	Port{1579, 6}:    "ioc-sea-lm",           // Missing description for ioc-sea-lm
	Port{1579, 17}:   "ioc-sea-lm",           // Missing description for ioc-sea-lm
	Port{1580, 6}:    "tn-tl-r1",             // tn-tl-r2
	Port{1580, 17}:   "tn-tl-r2",             // Missing description for tn-tl-r2
	Port{1581, 6}:    "mil-2045-47001",       // Missing description for mil-2045-47001
	Port{1581, 17}:   "mil-2045-47001",       // MIL-2045-47001
	Port{1582, 6}:    "msims",                // Missing description for msims
	Port{1582, 17}:   "msims",                // MSIMS
	Port{1583, 6}:    "simbaexpress",         // Missing description for simbaexpress
	Port{1583, 17}:   "simbaexpress",         // Missing description for simbaexpress
	Port{1584, 6}:    "tn-tl-fd2",            // Missing description for tn-tl-fd2
	Port{1584, 17}:   "tn-tl-fd2",            // Missing description for tn-tl-fd2
	Port{1585, 6}:    "intv",                 // Missing description for intv
	Port{1585, 17}:   "intv",                 // Missing description for intv
	Port{1586, 6}:    "ibm-abtact",           // Missing description for ibm-abtact
	Port{1586, 17}:   "ibm-abtact",           // Missing description for ibm-abtact
	Port{1587, 6}:    "pra_elmd",             // pra-elmd
	Port{1587, 17}:   "pra_elmd",             // Missing description for pra_elmd
	Port{1588, 6}:    "triquest-lm",          // Missing description for triquest-lm
	Port{1588, 17}:   "triquest-lm",          // Missing description for triquest-lm
	Port{1589, 6}:    "vqp",                  // Missing description for vqp
	Port{1589, 17}:   "vqp",                  // VQP
	Port{1590, 6}:    "gemini-lm",            // Missing description for gemini-lm
	Port{1590, 17}:   "gemini-lm",            // Missing description for gemini-lm
	Port{1591, 6}:    "ncpm-pm",              // Missing description for ncpm-pm
	Port{1591, 17}:   "ncpm-pm",              // Missing description for ncpm-pm
	Port{1592, 6}:    "commonspace",          // Missing description for commonspace
	Port{1592, 17}:   "commonspace",          // Missing description for commonspace
	Port{1593, 6}:    "mainsoft-lm",          // Missing description for mainsoft-lm
	Port{1593, 17}:   "mainsoft-lm",          // Missing description for mainsoft-lm
	Port{1594, 6}:    "sixtrak",              // Missing description for sixtrak
	Port{1594, 17}:   "sixtrak",              // Missing description for sixtrak
	Port{1595, 6}:    "radio",                // Missing description for radio
	Port{1595, 17}:   "radio",                // Missing description for radio
	Port{1596, 6}:    "radio-sm",             // radio-bc
	Port{1596, 17}:   "radio-bc",             // Missing description for radio-bc
	Port{1597, 6}:    "orbplus-iiop",         // Missing description for orbplus-iiop
	Port{1597, 17}:   "orbplus-iiop",         // Missing description for orbplus-iiop
	Port{1598, 6}:    "picknfs",              // Missing description for picknfs
	Port{1598, 17}:   "picknfs",              // Missing description for picknfs
	Port{1599, 6}:    "simbaservices",        // Missing description for simbaservices
	Port{1599, 17}:   "simbaservices",        // Missing description for simbaservices
	Port{1600, 6}:    "issd",                 // Missing description for issd
	Port{1600, 17}:   "issd",                 // Missing description for issd
	Port{1601, 6}:    "aas",                  // Missing description for aas
	Port{1601, 17}:   "aas",                  // Missing description for aas
	Port{1602, 6}:    "inspect",              // Missing description for inspect
	Port{1602, 17}:   "inspect",              // Missing description for inspect
	Port{1603, 6}:    "picodbc",              // pickodbc
	Port{1603, 17}:   "picodbc",              // pickodbc
	Port{1604, 6}:    "icabrowser",           // Missing description for icabrowser
	Port{1604, 17}:   "icabrowser",           // Missing description for icabrowser
	Port{1605, 6}:    "slp",                  // Salutation Manager (Salutation Protocol)
	Port{1605, 17}:   "slp",                  // Salutation Manager (Salutation Protocol)
	Port{1606, 6}:    "slm-api",              // Salutation Manager (SLM-API)
	Port{1606, 17}:   "slm-api",              // Salutation Manager (SLM-API)
	Port{1607, 6}:    "stt",                  // Missing description for stt
	Port{1607, 17}:   "stt",                  // Missing description for stt
	Port{1608, 6}:    "smart-lm",             // Smart Corp. License Manager
	Port{1608, 17}:   "smart-lm",             // Smart Corp. License Manager
	Port{1609, 6}:    "isysg-lm",             // Missing description for isysg-lm
	Port{1609, 17}:   "isysg-lm",             // Missing description for isysg-lm
	Port{1610, 6}:    "taurus-wh",            // Missing description for taurus-wh
	Port{1610, 17}:   "taurus-wh",            // Missing description for taurus-wh
	Port{1611, 6}:    "ill",                  // Inter Library Loan
	Port{1611, 17}:   "ill",                  // Inter Library Loan
	Port{1612, 6}:    "netbill-trans",        // NetBill Transaction Server
	Port{1612, 17}:   "netbill-trans",        // NetBill Transaction Server
	Port{1613, 6}:    "netbill-keyrep",       // NetBill Key Repository
	Port{1613, 17}:   "netbill-keyrep",       // NetBill Key Repository
	Port{1614, 6}:    "netbill-cred",         // NetBill Credential Server
	Port{1614, 17}:   "netbill-cred",         // NetBill Credential Server
	Port{1615, 6}:    "netbill-auth",         // NetBill Authorization Server
	Port{1615, 17}:   "netbill-auth",         // NetBill Authorization Server
	Port{1616, 6}:    "netbill-prod",         // NetBill Product Server
	Port{1616, 17}:   "netbill-prod",         // NetBill Product Server
	Port{1617, 6}:    "nimrod-agent",         // Nimrod Inter-Agent Communication
	Port{1617, 17}:   "nimrod-agent",         // Nimrod Inter-Agent Communication
	Port{1618, 6}:    "skytelnet",            // Missing description for skytelnet
	Port{1618, 17}:   "skytelnet",            // Missing description for skytelnet
	Port{1619, 6}:    "xs-openstorage",       // Missing description for xs-openstorage
	Port{1619, 17}:   "xs-openstorage",       // Missing description for xs-openstorage
	Port{1620, 6}:    "faxportwinport",       // Missing description for faxportwinport
	Port{1620, 17}:   "faxportwinport",       // Missing description for faxportwinport
	Port{1621, 6}:    "softdataphone",        // Missing description for softdataphone
	Port{1621, 17}:   "softdataphone",        // Missing description for softdataphone
	Port{1622, 6}:    "ontime",               // Missing description for ontime
	Port{1622, 17}:   "ontime",               // Missing description for ontime
	Port{1623, 6}:    "jaleosnd",             // Missing description for jaleosnd
	Port{1623, 17}:   "jaleosnd",             // Missing description for jaleosnd
	Port{1624, 6}:    "udp-sr-port",          // Missing description for udp-sr-port
	Port{1624, 17}:   "udp-sr-port",          // Missing description for udp-sr-port
	Port{1625, 6}:    "svs-omagent",          // Missing description for svs-omagent
	Port{1625, 17}:   "svs-omagent",          // Missing description for svs-omagent
	Port{1626, 6}:    "shockwave",            // Missing description for shockwave
	Port{1626, 17}:   "shockwave",            // Shockwave
	Port{1627, 6}:    "t128-gateway",         // T.128 Gateway
	Port{1627, 17}:   "t128-gateway",         // T.128 Gateway
	Port{1628, 6}:    "lontalk-norm",         // LonTalk normal
	Port{1628, 17}:   "lontalk-norm",         // LonTalk normal
	Port{1629, 6}:    "lontalk-urgnt",        // LonTalk urgent
	Port{1629, 17}:   "lontalk-urgnt",        // LonTalk urgent
	Port{1630, 6}:    "oraclenet8cman",       // Oracle Net8 Cman
	Port{1630, 17}:   "oraclenet8cman",       // Oracle Net8 Cman
	Port{1631, 6}:    "visitview",            // Visit view
	Port{1631, 17}:   "visitview",            // Visit view
	Port{1632, 6}:    "pammratc",             // Missing description for pammratc
	Port{1632, 17}:   "pammratc",             // PAMMRATC
	Port{1633, 6}:    "pammrpc",              // Missing description for pammrpc
	Port{1633, 17}:   "pammrpc",              // PAMMRPC
	Port{1634, 6}:    "loaprobe",             // Log On America Probe
	Port{1634, 17}:   "loaprobe",             // Log On America Probe
	Port{1635, 6}:    "edb-server1",          // EDB Server 1
	Port{1635, 17}:   "edb-server1",          // EDB Server 1
	Port{1636, 6}:    "isdc",                 // ISP shared public data control
	Port{1636, 17}:   "isdc",                 // ISP shared public data control
	Port{1637, 6}:    "islc",                 // ISP shared local data control
	Port{1637, 17}:   "islc",                 // ISP shared local data control
	Port{1638, 6}:    "ismc",                 // ISP shared management control
	Port{1638, 17}:   "ismc",                 // ISP shared management control
	Port{1639, 6}:    "cert-initiator",       // Missing description for cert-initiator
	Port{1639, 17}:   "cert-initiator",       // Missing description for cert-initiator
	Port{1640, 6}:    "cert-responder",       // Missing description for cert-responder
	Port{1640, 17}:   "cert-responder",       // Missing description for cert-responder
	Port{1641, 6}:    "invision",             // Missing description for invision
	Port{1641, 17}:   "invision",             // InVision
	Port{1642, 6}:    "isis-am",              // Missing description for isis-am
	Port{1642, 17}:   "isis-am",              // Missing description for isis-am
	Port{1643, 6}:    "isis-ambc",            // Missing description for isis-ambc
	Port{1643, 17}:   "isis-ambc",            // Missing description for isis-ambc
	Port{1644, 6}:    "saiseh",               // Satellite-data Acquisition System 4
	Port{1645, 6}:    "sightline",            // Missing description for sightline
	Port{1645, 17}:   "radius",               // radius authentication
	Port{1646, 6}:    "sa-msg-port",          // Missing description for sa-msg-port
	Port{1646, 17}:   "radacct",              // radius accounting
	Port{1647, 6}:    "rsap",                 // Missing description for rsap
	Port{1647, 17}:   "rsap",                 // Missing description for rsap
	Port{1648, 6}:    "concurrent-lm",        // Missing description for concurrent-lm
	Port{1648, 17}:   "concurrent-lm",        // Missing description for concurrent-lm
	Port{1649, 6}:    "kermit",               // Missing description for kermit
	Port{1649, 17}:   "kermit",               // Missing description for kermit
	Port{1650, 6}:    "nkd",                  // nkdn
	Port{1650, 17}:   "nkd",                  // Missing description for nkd
	Port{1651, 6}:    "shiva_confsrvr",       // shiva-confsrvr
	Port{1651, 17}:   "shiva_confsrvr",       // Missing description for shiva_confsrvr
	Port{1652, 6}:    "xnmp",                 // Missing description for xnmp
	Port{1652, 17}:   "xnmp",                 // Missing description for xnmp
	Port{1653, 6}:    "alphatech-lm",         // Missing description for alphatech-lm
	Port{1653, 17}:   "alphatech-lm",         // Missing description for alphatech-lm
	Port{1654, 6}:    "stargatealerts",       // Missing description for stargatealerts
	Port{1654, 17}:   "stargatealerts",       // Missing description for stargatealerts
	Port{1655, 6}:    "dec-mbadmin",          // Missing description for dec-mbadmin
	Port{1655, 17}:   "dec-mbadmin",          // Missing description for dec-mbadmin
	Port{1656, 6}:    "dec-mbadmin-h",        // Missing description for dec-mbadmin-h
	Port{1656, 17}:   "dec-mbadmin-h",        // Missing description for dec-mbadmin-h
	Port{1657, 6}:    "fujitsu-mmpdc",        // Missing description for fujitsu-mmpdc
	Port{1657, 17}:   "fujitsu-mmpdc",        // Missing description for fujitsu-mmpdc
	Port{1658, 6}:    "sixnetudr",            // Missing description for sixnetudr
	Port{1658, 17}:   "sixnetudr",            // Missing description for sixnetudr
	Port{1659, 6}:    "sg-lm",                // Silicon Grail License Manager
	Port{1659, 17}:   "sg-lm",                // Silicon Grail License Manager
	Port{1660, 6}:    "skip-mc-gikreq",       // Missing description for skip-mc-gikreq
	Port{1660, 17}:   "skip-mc-gikreq",       // Missing description for skip-mc-gikreq
	Port{1661, 6}:    "netview-aix-1",        // Missing description for netview-aix-1
	Port{1661, 17}:   "netview-aix-1",        // Missing description for netview-aix-1
	Port{1662, 6}:    "netview-aix-2",        // Missing description for netview-aix-2
	Port{1662, 17}:   "netview-aix-2",        // Missing description for netview-aix-2
	Port{1663, 6}:    "netview-aix-3",        // Missing description for netview-aix-3
	Port{1663, 17}:   "netview-aix-3",        // Missing description for netview-aix-3
	Port{1664, 6}:    "netview-aix-4",        // Missing description for netview-aix-4
	Port{1664, 17}:   "netview-aix-4",        // Missing description for netview-aix-4
	Port{1665, 6}:    "netview-aix-5",        // Missing description for netview-aix-5
	Port{1665, 17}:   "netview-aix-5",        // Missing description for netview-aix-5
	Port{1666, 6}:    "netview-aix-6",        // Missing description for netview-aix-6
	Port{1666, 17}:   "netview-aix-6",        // Missing description for netview-aix-6
	Port{1667, 6}:    "netview-aix-7",        // Missing description for netview-aix-7
	Port{1667, 17}:   "netview-aix-7",        // Missing description for netview-aix-7
	Port{1668, 6}:    "netview-aix-8",        // Missing description for netview-aix-8
	Port{1668, 17}:   "netview-aix-8",        // Missing description for netview-aix-8
	Port{1669, 6}:    "netview-aix-9",        // Missing description for netview-aix-9
	Port{1669, 17}:   "netview-aix-9",        // Missing description for netview-aix-9
	Port{1670, 6}:    "netview-aix-10",       // Missing description for netview-aix-10
	Port{1670, 17}:   "netview-aix-10",       // Missing description for netview-aix-10
	Port{1671, 6}:    "netview-aix-11",       // Missing description for netview-aix-11
	Port{1671, 17}:   "netview-aix-11",       // Missing description for netview-aix-11
	Port{1672, 6}:    "netview-aix-12",       // Missing description for netview-aix-12
	Port{1672, 17}:   "netview-aix-12",       // Missing description for netview-aix-12
	Port{1673, 6}:    "proshare-mc-1",        // Intel Proshare Multicast
	Port{1673, 17}:   "proshare-mc-1",        // Intel Proshare Multicast
	Port{1674, 6}:    "proshare-mc-2",        // Intel Proshare Multicast
	Port{1674, 17}:   "proshare-mc-2",        // Intel Proshare Multicast
	Port{1675, 6}:    "pdp",                  // Pacific Data Products
	Port{1675, 17}:   "pdp",                  // Pacific Data Products
	Port{1676, 6}:    "netcomm1",             // netcomm2
	Port{1676, 17}:   "netcomm2",             // Missing description for netcomm2
	Port{1677, 6}:    "groupwise",            // Missing description for groupwise
	Port{1677, 17}:   "groupwise",            // Missing description for groupwise
	Port{1678, 6}:    "prolink",              // Missing description for prolink
	Port{1678, 17}:   "prolink",              // Missing description for prolink
	Port{1679, 6}:    "darcorp-lm",           // Missing description for darcorp-lm
	Port{1679, 17}:   "darcorp-lm",           // Missing description for darcorp-lm
	Port{1680, 6}:    "CarbonCopy",           // microcom-sbp
	Port{1680, 17}:   "microcom-sbp",         // Missing description for microcom-sbp
	Port{1681, 6}:    "sd-elmd",              // Missing description for sd-elmd
	Port{1681, 17}:   "sd-elmd",              // Missing description for sd-elmd
	Port{1682, 6}:    "lanyon-lantern",       // Missing description for lanyon-lantern
	Port{1682, 17}:   "lanyon-lantern",       // Missing description for lanyon-lantern
	Port{1683, 6}:    "ncpm-hip",             // Missing description for ncpm-hip
	Port{1683, 17}:   "ncpm-hip",             // Missing description for ncpm-hip
	Port{1684, 6}:    "snaresecure",          // Missing description for snaresecure
	Port{1684, 17}:   "snaresecure",          // SnareSecure
	Port{1685, 6}:    "n2nremote",            // Missing description for n2nremote
	Port{1685, 17}:   "n2nremote",            // Missing description for n2nremote
	Port{1686, 6}:    "cvmon",                // Missing description for cvmon
	Port{1686, 17}:   "cvmon",                // Missing description for cvmon
	Port{1687, 6}:    "nsjtp-ctrl",           // Missing description for nsjtp-ctrl
	Port{1687, 17}:   "nsjtp-ctrl",           // Missing description for nsjtp-ctrl
	Port{1688, 6}:    "nsjtp-data",           // Missing description for nsjtp-data
	Port{1688, 17}:   "nsjtp-data",           // Missing description for nsjtp-data
	Port{1689, 6}:    "firefox",              // Missing description for firefox
	Port{1689, 17}:   "firefox",              // Missing description for firefox
	Port{1690, 6}:    "ng-umds",              // Missing description for ng-umds
	Port{1690, 17}:   "ng-umds",              // Missing description for ng-umds
	Port{1691, 6}:    "empire-empuma",        // Missing description for empire-empuma
	Port{1691, 17}:   "empire-empuma",        // Missing description for empire-empuma
	Port{1692, 6}:    "sstsys-lm",            // Missing description for sstsys-lm
	Port{1692, 17}:   "sstsys-lm",            // Missing description for sstsys-lm
	Port{1693, 6}:    "rrirtr",               // Missing description for rrirtr
	Port{1693, 17}:   "rrirtr",               // Missing description for rrirtr
	Port{1694, 6}:    "rrimwm",               // Missing description for rrimwm
	Port{1694, 17}:   "rrimwm",               // Missing description for rrimwm
	Port{1695, 6}:    "rrilwm",               // Missing description for rrilwm
	Port{1695, 17}:   "rrilwm",               // Missing description for rrilwm
	Port{1696, 6}:    "rrifmm",               // Missing description for rrifmm
	Port{1696, 17}:   "rrifmm",               // Missing description for rrifmm
	Port{1697, 6}:    "rrisat",               // Missing description for rrisat
	Port{1697, 17}:   "rrisat",               // Missing description for rrisat
	Port{1698, 6}:    "rsvp-encap-1",         // RSVP-ENCAPSULATION-1
	Port{1698, 17}:   "rsvp-encap-1",         // RSVP-ENCAPSULATION-1
	Port{1699, 6}:    "rsvp-encap-2",         // RSVP-ENCAPSULATION-2
	Port{1699, 17}:   "rsvp-encap-2",         // RSVP-ENCAPSULATION-2
	Port{1700, 6}:    "mps-raft",             // Missing description for mps-raft
	Port{1700, 17}:   "mps-raft",             // Missing description for mps-raft
	Port{1701, 6}:    "l2f",                  // l2tp
	Port{1701, 17}:   "L2TP",                 // Missing description for L2TP
	Port{1702, 6}:    "deskshare",            // Missing description for deskshare
	Port{1702, 17}:   "deskshare",            // Missing description for deskshare
	Port{1703, 6}:    "hb-engine",            // Missing description for hb-engine
	Port{1703, 17}:   "hb-engine",            // Missing description for hb-engine
	Port{1704, 6}:    "bcs-broker",           // Missing description for bcs-broker
	Port{1704, 17}:   "bcs-broker",           // Missing description for bcs-broker
	Port{1705, 6}:    "slingshot",            // Missing description for slingshot
	Port{1705, 17}:   "slingshot",            // Missing description for slingshot
	Port{1706, 6}:    "jetform",              // Missing description for jetform
	Port{1706, 17}:   "jetform",              // Missing description for jetform
	Port{1707, 6}:    "vdmplay",              // Missing description for vdmplay
	Port{1707, 17}:   "vdmplay",              // Missing description for vdmplay
	Port{1708, 6}:    "gat-lmd",              // Missing description for gat-lmd
	Port{1708, 17}:   "gat-lmd",              // Missing description for gat-lmd
	Port{1709, 6}:    "centra",               // Missing description for centra
	Port{1709, 17}:   "centra",               // Missing description for centra
	Port{1710, 6}:    "impera",               // Missing description for impera
	Port{1710, 17}:   "impera",               // Missing description for impera
	Port{1711, 6}:    "pptconference",        // Missing description for pptconference
	Port{1711, 17}:   "pptconference",        // Missing description for pptconference
	Port{1712, 6}:    "registrar",            // resource monitoring service
	Port{1712, 17}:   "registrar",            // resource monitoring service
	Port{1713, 6}:    "conferencetalk",       // Missing description for conferencetalk
	Port{1713, 17}:   "conferencetalk",       // ConferenceTalk
	Port{1714, 6}:    "sesi-lm",              // Missing description for sesi-lm
	Port{1714, 17}:   "sesi-lm",              // Missing description for sesi-lm
	Port{1715, 6}:    "houdini-lm",           // Missing description for houdini-lm
	Port{1715, 17}:   "houdini-lm",           // Missing description for houdini-lm
	Port{1716, 6}:    "xmsg",                 // Missing description for xmsg
	Port{1716, 17}:   "xmsg",                 // Missing description for xmsg
	Port{1717, 6}:    "fj-hdnet",             // Missing description for fj-hdnet
	Port{1717, 17}:   "fj-hdnet",             // Missing description for fj-hdnet
	Port{1718, 6}:    "h323gatedisc",         // H.323 Multicast Gatekeeper Discover
	Port{1718, 17}:   "h225gatedisc",         // H.225 gatekeeper discovery
	Port{1719, 6}:    "h323gatestat",         // H.323 Unicast Gatekeeper Signaling
	Port{1719, 17}:   "h323gatestat",         // H.323 Gatestat
	Port{1720, 6}:    "h323q931",             // h323hostcall | Interactive media | H.323 Call Control Signalling | H.323 Call Control
	Port{1720, 17}:   "h323hostcall",         // Missing description for h323hostcall
	Port{1721, 6}:    "caicci",               // Missing description for caicci
	Port{1721, 17}:   "caicci",               // Missing description for caicci
	Port{1722, 6}:    "hks-lm",               // HKS License Manager
	Port{1722, 17}:   "hks-lm",               // HKS License Manager
	Port{1723, 6}:    "pptp",                 // Point-to-point tunnelling protocol
	Port{1723, 17}:   "pptp",                 // Missing description for pptp
	Port{1724, 6}:    "csbphonemaster",       // Missing description for csbphonemaster
	Port{1724, 17}:   "csbphonemaster",       // Missing description for csbphonemaster
	Port{1725, 6}:    "iden-ralp",            // Missing description for iden-ralp
	Port{1725, 17}:   "iden-ralp",            // Missing description for iden-ralp
	Port{1726, 6}:    "iberiagames",          // Missing description for iberiagames
	Port{1726, 17}:   "iberiagames",          // IBERIAGAMES
	Port{1727, 6}:    "winddx",               // Missing description for winddx
	Port{1727, 17}:   "winddx",               // Missing description for winddx
	Port{1728, 6}:    "telindus",             // Missing description for telindus
	Port{1728, 17}:   "telindus",             // TELINDUS
	Port{1729, 6}:    "citynl",               // CityNL License Management
	Port{1729, 17}:   "citynl",               // CityNL License Management
	Port{1730, 6}:    "roketz",               // Missing description for roketz
	Port{1730, 17}:   "roketz",               // Missing description for roketz
	Port{1731, 6}:    "msiccp",               // Missing description for msiccp
	Port{1731, 17}:   "msiccp",               // MSICCP
	Port{1732, 6}:    "proxim",               // Missing description for proxim
	Port{1732, 17}:   "proxim",               // Missing description for proxim
	Port{1733, 6}:    "siipat",               // SIMS - SIIPAT Protocol for Alarm Transmission
	Port{1733, 17}:   "siipat",               // SIMS - SIIPAT Protocol for Alarm Transmission
	Port{1734, 6}:    "cambertx-lm",          // Camber Corporation License Management
	Port{1734, 17}:   "cambertx-lm",          // Camber Corporation License Management
	Port{1735, 6}:    "privatechat",          // Missing description for privatechat
	Port{1735, 17}:   "privatechat",          // PrivateChat
	Port{1736, 6}:    "street-stream",        // Missing description for street-stream
	Port{1736, 17}:   "street-stream",        // Missing description for street-stream
	Port{1737, 6}:    "ultimad",              // Missing description for ultimad
	Port{1737, 17}:   "ultimad",              // Missing description for ultimad
	Port{1738, 6}:    "gamegen1",             // Missing description for gamegen1
	Port{1738, 17}:   "gamegen1",             // GameGen1
	Port{1739, 6}:    "webaccess",            // Missing description for webaccess
	Port{1739, 17}:   "webaccess",            // Missing description for webaccess
	Port{1740, 6}:    "encore",               // Missing description for encore
	Port{1740, 17}:   "encore",               // Missing description for encore
	Port{1741, 6}:    "cisco-net-mgmt",       // Missing description for cisco-net-mgmt
	Port{1741, 17}:   "cisco-net-mgmt",       // Missing description for cisco-net-mgmt
	Port{1742, 6}:    "3Com-nsd",             // Missing description for 3Com-nsd
	Port{1742, 17}:   "3Com-nsd",             // Missing description for 3Com-nsd
	Port{1743, 6}:    "cinegrfx-lm",          // Cinema Graphics License Manager
	Port{1743, 17}:   "cinegrfx-lm",          // Cinema Graphics License Manager
	Port{1744, 6}:    "ncpm-ft",              // Missing description for ncpm-ft
	Port{1744, 17}:   "ncpm-ft",              // Missing description for ncpm-ft
	Port{1745, 6}:    "remote-winsock",       // Missing description for remote-winsock
	Port{1745, 17}:   "remote-winsock",       // Missing description for remote-winsock
	Port{1746, 6}:    "ftrapid-1",            // Missing description for ftrapid-1
	Port{1746, 17}:   "ftrapid-1",            // Missing description for ftrapid-1
	Port{1747, 6}:    "ftrapid-2",            // Missing description for ftrapid-2
	Port{1747, 17}:   "ftrapid-2",            // Missing description for ftrapid-2
	Port{1748, 6}:    "oracle-em1",           // Missing description for oracle-em1
	Port{1748, 17}:   "oracle-em1",           // Missing description for oracle-em1
	Port{1749, 6}:    "aspen-services",       // Missing description for aspen-services
	Port{1749, 17}:   "aspen-services",       // Missing description for aspen-services
	Port{1750, 6}:    "sslp",                 // Simple Socket Library's PortMaster
	Port{1750, 17}:   "sslp",                 // Simple Socket Library's PortMaster
	Port{1751, 6}:    "swiftnet",             // Missing description for swiftnet
	Port{1751, 17}:   "swiftnet",             // SwiftNet
	Port{1752, 6}:    "lofr-lm",              // Leap of Faith Research License Manager
	Port{1752, 17}:   "lofr-lm",              // Leap of Faith Research License Manager
	Port{1753, 6}:    "predatar-comms",       // Predatar Comms Service
	Port{1754, 6}:    "oracle-em2",           // Missing description for oracle-em2
	Port{1754, 17}:   "oracle-em2",           // Missing description for oracle-em2
	Port{1755, 6}:    "wms",                  // Windows media service | ms-streaming
	Port{1755, 17}:   "ms-streaming",         // Missing description for ms-streaming
	Port{1756, 6}:    "capfast-lmd",          // Missing description for capfast-lmd
	Port{1756, 17}:   "capfast-lmd",          // Missing description for capfast-lmd
	Port{1757, 6}:    "cnhrp",                // Missing description for cnhrp
	Port{1757, 17}:   "cnhrp",                // Missing description for cnhrp
	Port{1758, 6}:    "tftp-mcast",           // Missing description for tftp-mcast
	Port{1758, 17}:   "tftp-mcast",           // Missing description for tftp-mcast
	Port{1759, 6}:    "spss-lm",              // SPSS License Manager
	Port{1759, 17}:   "spss-lm",              // SPSS License Manager
	Port{1760, 6}:    "www-ldap-gw",          // Missing description for www-ldap-gw
	Port{1760, 17}:   "www-ldap-gw",          // Missing description for www-ldap-gw
	Port{1761, 6}:    "landesk-rc",           // LANDesk Remote Control | cft-0
	Port{1761, 17}:   "cft-0",                // Missing description for cft-0
	Port{1762, 6}:    "landesk-rc",           // LANDesk Remote Control | cft-1
	Port{1762, 17}:   "cft-1",                // Missing description for cft-1
	Port{1763, 6}:    "landesk-rc",           // LANDesk Remote Control | cft-2
	Port{1763, 17}:   "cft-2",                // Missing description for cft-2
	Port{1764, 6}:    "landesk-rc",           // LANDesk Remote Control | cft-3
	Port{1764, 17}:   "cft-3",                // Missing description for cft-3
	Port{1765, 6}:    "cft-4",                // Missing description for cft-4
	Port{1765, 17}:   "cft-4",                // Missing description for cft-4
	Port{1766, 6}:    "cft-5",                // Missing description for cft-5
	Port{1766, 17}:   "cft-5",                // Missing description for cft-5
	Port{1767, 6}:    "cft-6",                // Missing description for cft-6
	Port{1767, 17}:   "cft-6",                // Missing description for cft-6
	Port{1768, 6}:    "cft-7",                // Missing description for cft-7
	Port{1768, 17}:   "cft-7",                // Missing description for cft-7
	Port{1769, 6}:    "bmc-net-adm",          // Missing description for bmc-net-adm
	Port{1769, 17}:   "bmc-net-adm",          // Missing description for bmc-net-adm
	Port{1770, 6}:    "bmc-net-svc",          // Missing description for bmc-net-svc
	Port{1770, 17}:   "bmc-net-svc",          // Missing description for bmc-net-svc
	Port{1771, 6}:    "vaultbase",            // Missing description for vaultbase
	Port{1771, 17}:   "vaultbase",            // Missing description for vaultbase
	Port{1772, 6}:    "essweb-gw",            // EssWeb Gateway
	Port{1772, 17}:   "essweb-gw",            // EssWeb Gateway
	Port{1773, 6}:    "kmscontrol",           // Missing description for kmscontrol
	Port{1773, 17}:   "kmscontrol",           // KMSControl
	Port{1774, 6}:    "global-dtserv",        // Missing description for global-dtserv
	Port{1774, 17}:   "global-dtserv",        // Missing description for global-dtserv
	Port{1775, 6}:    "vdab",                 // data interchange between visual processing containers
	Port{1776, 6}:    "femis",                // Federal Emergency Management Information System
	Port{1776, 17}:   "femis",                // Federal Emergency Management Information System
	Port{1777, 6}:    "powerguardian",        // Missing description for powerguardian
	Port{1777, 17}:   "powerguardian",        // Missing description for powerguardian
	Port{1778, 6}:    "prodigy-intrnet",      // prodigy-internet
	Port{1778, 17}:   "prodigy-intrnet",      // prodigy-internet
	Port{1779, 6}:    "pharmasoft",           // Missing description for pharmasoft
	Port{1779, 17}:   "pharmasoft",           // Missing description for pharmasoft
	Port{1780, 6}:    "dpkeyserv",            // Missing description for dpkeyserv
	Port{1780, 17}:   "dpkeyserv",            // Missing description for dpkeyserv
	Port{1781, 6}:    "answersoft-lm",        // Missing description for answersoft-lm
	Port{1781, 17}:   "answersoft-lm",        // Missing description for answersoft-lm
	Port{1782, 6}:    "hp-hcip",              // Missing description for hp-hcip
	Port{1782, 17}:   "hp-hcip",              // Missing description for hp-hcip
	Port{1784, 6}:    "finle-lm",             // Finle License Manager
	Port{1784, 17}:   "finle-lm",             // Finle License Manager
	Port{1785, 6}:    "windlm",               // Wind River Systems License Manager
	Port{1785, 17}:   "windlm",               // Wind River Systems License Manager
	Port{1786, 6}:    "funk-logger",          // Missing description for funk-logger
	Port{1786, 17}:   "funk-logger",          // Missing description for funk-logger
	Port{1787, 6}:    "funk-license",         // Missing description for funk-license
	Port{1787, 17}:   "funk-license",         // Missing description for funk-license
	Port{1788, 6}:    "psmond",               // Missing description for psmond
	Port{1788, 17}:   "psmond",               // Missing description for psmond
	Port{1789, 6}:    "hello",                // Missing description for hello
	Port{1789, 17}:   "hello",                // Missing description for hello
	Port{1790, 6}:    "nmsp",                 // Narrative Media Streaming Protocol
	Port{1790, 17}:   "nmsp",                 // Narrative Media Streaming Protocol
	Port{1791, 6}:    "ea1",                  // Missing description for ea1
	Port{1791, 17}:   "ea1",                  // EA1
	Port{1792, 6}:    "ibm-dt-2",             // Missing description for ibm-dt-2
	Port{1792, 17}:   "ibm-dt-2",             // Missing description for ibm-dt-2
	Port{1793, 6}:    "rsc-robot",            // Missing description for rsc-robot
	Port{1793, 17}:   "rsc-robot",            // Missing description for rsc-robot
	Port{1794, 6}:    "cera-bcm",             // Missing description for cera-bcm
	Port{1794, 17}:   "cera-bcm",             // Missing description for cera-bcm
	Port{1795, 6}:    "dpi-proxy",            // Missing description for dpi-proxy
	Port{1795, 17}:   "dpi-proxy",            // Missing description for dpi-proxy
	Port{1796, 6}:    "vocaltec-admin",       // Vocaltec Server Administration
	Port{1796, 17}:   "vocaltec-admin",       // Vocaltec Server Administration
	Port{1797, 6}:    "uma",                  // Missing description for uma
	Port{1797, 17}:   "uma",                  // UMA
	Port{1798, 6}:    "etp",                  // Event Transfer Protocol
	Port{1798, 17}:   "etp",                  // Event Transfer Protocol
	Port{1799, 6}:    "netrisk",              // Missing description for netrisk
	Port{1799, 17}:   "netrisk",              // NETRISK
	Port{1800, 6}:    "ansys-lm",             // ANSYS-License manager
	Port{1800, 17}:   "ansys-lm",             // ANSYS-License manager
	Port{1801, 6}:    "msmq",                 // Microsoft Message Queuing | Microsoft Message Que
	Port{1801, 17}:   "msmq",                 // Microsoft Message Que
	Port{1802, 6}:    "concomp1",             // Missing description for concomp1
	Port{1802, 17}:   "concomp1",             // ConComp1
	Port{1803, 6}:    "hp-hcip-gwy",          // Missing description for hp-hcip-gwy
	Port{1803, 17}:   "hp-hcip-gwy",          // HP-HCIP-GWY
	Port{1804, 6}:    "enl",                  // Missing description for enl
	Port{1804, 17}:   "enl",                  // ENL
	Port{1805, 6}:    "enl-name",             // Missing description for enl-name
	Port{1805, 17}:   "enl-name",             // ENL-Name
	Port{1806, 6}:    "musiconline",          // Missing description for musiconline
	Port{1806, 17}:   "musiconline",          // Musiconline
	Port{1807, 6}:    "fhsp",                 // Fujitsu Hot Standby Protocol
	Port{1807, 17}:   "fhsp",                 // Fujitsu Hot Standby Protocol
	Port{1808, 6}:    "oracle-vp2",           // Missing description for oracle-vp2
	Port{1808, 17}:   "oracle-vp2",           // Oracle-VP2
	Port{1809, 6}:    "oracle-vp1",           // Missing description for oracle-vp1
	Port{1809, 17}:   "oracle-vp1",           // Oracle-VP1
	Port{1810, 6}:    "jerand-lm",            // Jerand License Manager
	Port{1810, 17}:   "jerand-lm",            // Jerand License Manager
	Port{1811, 6}:    "scientia-sdb",         // Missing description for scientia-sdb
	Port{1811, 17}:   "scientia-sdb",         // Scientia-SDB
	Port{1812, 132}:  "radius",               // RADIUS authentication protocol (RFC 2138)
	Port{1812, 6}:    "radius",               // RADIUS
	Port{1812, 17}:   "radius",               // RADIUS authentication protocol (RFC 2138)
	Port{1813, 132}:  "radacct",              // radius-acct | RADIUS accounting protocol (RFC 2139) | RADIUS Accounting
	Port{1813, 6}:    "radius-acct",          // RADIUS Accounting
	Port{1813, 17}:   "radacct",              // RADIUS accounting protocol (RFC 2139)
	Port{1814, 6}:    "tdp-suite",            // TDP Suite
	Port{1814, 17}:   "tdp-suite",            // TDP Suite
	Port{1815, 6}:    "mmpft",                // Missing description for mmpft
	Port{1815, 17}:   "mmpft",                // MMPFT
	Port{1816, 6}:    "harp",                 // Missing description for harp
	Port{1816, 17}:   "harp",                 // HARP
	Port{1817, 6}:    "rkb-oscs",             // Missing description for rkb-oscs
	Port{1817, 17}:   "rkb-oscs",             // RKB-OSCS
	Port{1818, 6}:    "etftp",                // Enhanced Trivial File Transfer Protocol
	Port{1818, 17}:   "etftp",                // Enhanced Trivial File Transfer Protocol
	Port{1819, 6}:    "plato-lm",             // Plato License Manager
	Port{1819, 17}:   "plato-lm",             // Plato License Manager
	Port{1820, 6}:    "mcagent",              // Missing description for mcagent
	Port{1820, 17}:   "mcagent",              // Missing description for mcagent
	Port{1821, 6}:    "donnyworld",           // Missing description for donnyworld
	Port{1821, 17}:   "donnyworld",           // Missing description for donnyworld
	Port{1822, 6}:    "es-elmd",              // Missing description for es-elmd
	Port{1822, 17}:   "es-elmd",              // Missing description for es-elmd
	Port{1823, 6}:    "unisys-lm",            // Unisys Natural Language License Manager
	Port{1823, 17}:   "unisys-lm",            // Unisys Natural Language License Manager
	Port{1824, 6}:    "metrics-pas",          // Missing description for metrics-pas
	Port{1824, 17}:   "metrics-pas",          // Missing description for metrics-pas
	Port{1825, 6}:    "direcpc-video",        // DirecPC Video
	Port{1825, 17}:   "direcpc-video",        // DirecPC Video
	Port{1826, 6}:    "ardt",                 // Missing description for ardt
	Port{1826, 17}:   "ardt",                 // ARDT
	Port{1827, 6}:    "pcm",                  // asi | PCM Agent (AutoSecure Policy Compliance Manager | ASI
	Port{1827, 17}:   "asi",                  // ASI
	Port{1828, 6}:    "itm-mcell-u",          // Missing description for itm-mcell-u
	Port{1828, 17}:   "itm-mcell-u",          // Missing description for itm-mcell-u
	Port{1829, 6}:    "optika-emedia",        // Optika eMedia
	Port{1829, 17}:   "optika-emedia",        // Optika eMedia
	Port{1830, 6}:    "net8-cman",            // Oracle Net8 CMan Admin
	Port{1830, 17}:   "net8-cman",            // Oracle Net8 CMan Admin
	Port{1831, 6}:    "myrtle",               // Missing description for myrtle
	Port{1831, 17}:   "myrtle",               // Myrtle
	Port{1832, 6}:    "tht-treasure",         // ThoughtTreasure
	Port{1832, 17}:   "tht-treasure",         // ThoughtTreasure
	Port{1833, 6}:    "udpradio",             // Missing description for udpradio
	Port{1833, 17}:   "udpradio",             // Missing description for udpradio
	Port{1834, 6}:    "ardusuni",             // ARDUS Unicast
	Port{1834, 17}:   "ardusuni",             // ARDUS Unicast
	Port{1835, 6}:    "ardusmul",             // ARDUS Multicast
	Port{1835, 17}:   "ardusmul",             // ARDUS Multicast
	Port{1836, 6}:    "ste-smsc",             // Missing description for ste-smsc
	Port{1836, 17}:   "ste-smsc",             // Missing description for ste-smsc
	Port{1837, 6}:    "csoft1",               // Missing description for csoft1
	Port{1837, 17}:   "csoft1",               // Missing description for csoft1
	Port{1838, 6}:    "talnet",               // Missing description for talnet
	Port{1838, 17}:   "talnet",               // TALNET
	Port{1839, 6}:    "netopia-vo1",          // Missing description for netopia-vo1
	Port{1839, 17}:   "netopia-vo1",          // Missing description for netopia-vo1
	Port{1840, 6}:    "netopia-vo2",          // Missing description for netopia-vo2
	Port{1840, 17}:   "netopia-vo2",          // Missing description for netopia-vo2
	Port{1841, 6}:    "netopia-vo3",          // Missing description for netopia-vo3
	Port{1841, 17}:   "netopia-vo3",          // Missing description for netopia-vo3
	Port{1842, 6}:    "netopia-vo4",          // Missing description for netopia-vo4
	Port{1842, 17}:   "netopia-vo4",          // Missing description for netopia-vo4
	Port{1843, 6}:    "netopia-vo5",          // Missing description for netopia-vo5
	Port{1843, 17}:   "netopia-vo5",          // Missing description for netopia-vo5
	Port{1844, 6}:    "direcpc-dll",          // Missing description for direcpc-dll
	Port{1844, 17}:   "direcpc-dll",          // DirecPC-DLL
	Port{1845, 6}:    "altalink",             // Missing description for altalink
	Port{1845, 17}:   "altalink",             // Missing description for altalink
	Port{1846, 6}:    "tunstall-pnc",         // Tunstall PNC
	Port{1846, 17}:   "tunstall-pnc",         // Tunstall PNC
	Port{1847, 6}:    "slp-notify",           // SLP Notification
	Port{1847, 17}:   "slp-notify",           // SLP Notification
	Port{1848, 6}:    "fjdocdist",            // Missing description for fjdocdist
	Port{1848, 17}:   "fjdocdist",            // Missing description for fjdocdist
	Port{1849, 6}:    "alpha-sms",            // Missing description for alpha-sms
	Port{1849, 17}:   "alpha-sms",            // ALPHA-SMS
	Port{1850, 6}:    "gsi",                  // Missing description for gsi
	Port{1850, 17}:   "gsi",                  // GSI
	Port{1851, 6}:    "ctcd",                 // Missing description for ctcd
	Port{1851, 17}:   "ctcd",                 // Missing description for ctcd
	Port{1852, 6}:    "virtual-time",         // Virtual Time
	Port{1852, 17}:   "virtual-time",         // Virtual Time
	Port{1853, 6}:    "vids-avtp",            // Missing description for vids-avtp
	Port{1853, 17}:   "vids-avtp",            // VIDS-AVTP
	Port{1854, 6}:    "buddy-draw",           // Buddy Draw
	Port{1854, 17}:   "buddy-draw",           // Buddy Draw
	Port{1855, 6}:    "fiorano-rtrsvc",       // Fiorano RtrSvc
	Port{1855, 17}:   "fiorano-rtrsvc",       // Fiorano RtrSvc
	Port{1856, 6}:    "fiorano-msgsvc",       // Fiorano MsgSvc
	Port{1856, 17}:   "fiorano-msgsvc",       // Fiorano MsgSvc
	Port{1857, 6}:    "datacaptor",           // Missing description for datacaptor
	Port{1857, 17}:   "datacaptor",           // DataCaptor
	Port{1858, 6}:    "privateark",           // Missing description for privateark
	Port{1858, 17}:   "privateark",           // PrivateArk
	Port{1859, 6}:    "gammafetchsvr",        // Gamma Fetcher Server
	Port{1859, 17}:   "gammafetchsvr",        // Gamma Fetcher Server
	Port{1860, 6}:    "sunscalar-svc",        // SunSCALAR Services
	Port{1860, 17}:   "sunscalar-svc",        // SunSCALAR Services
	Port{1861, 6}:    "lecroy-vicp",          // LeCroy VICP
	Port{1861, 17}:   "lecroy-vicp",          // LeCroy VICP
	Port{1862, 6}:    "mysql-cm-agent",       // MySQL Cluster Manager Agent
	Port{1862, 17}:   "mysql-cm-agent",       // MySQL Cluster Manager Agent
	Port{1863, 6}:    "msnp",                 // MSN Messenger
	Port{1863, 17}:   "msnp",                 // MSN Messenger
	Port{1864, 6}:    "paradym-31",           // paradym-31port | Paradym 31 Port
	Port{1864, 17}:   "paradym-31",           // Missing description for paradym-31
	Port{1865, 6}:    "entp",                 // Missing description for entp
	Port{1865, 17}:   "entp",                 // ENTP
	Port{1866, 6}:    "swrmi",                // Missing description for swrmi
	Port{1866, 17}:   "swrmi",                // Missing description for swrmi
	Port{1867, 6}:    "udrive",               // Missing description for udrive
	Port{1867, 17}:   "udrive",               // UDRIVE
	Port{1868, 6}:    "viziblebrowser",       // Missing description for viziblebrowser
	Port{1868, 17}:   "viziblebrowser",       // VizibleBrowser
	Port{1869, 6}:    "transact",             // Missing description for transact
	Port{1869, 17}:   "transact",             // TransAct
	Port{1870, 6}:    "sunscalar-dns",        // SunSCALAR DNS Service
	Port{1870, 17}:   "sunscalar-dns",        // SunSCALAR DNS Service
	Port{1871, 6}:    "canocentral0",         // Cano Central 0
	Port{1871, 17}:   "canocentral0",         // Cano Central 0
	Port{1872, 6}:    "canocentral1",         // Cano Central 1
	Port{1872, 17}:   "canocentral1",         // Cano Central 1
	Port{1873, 6}:    "fjmpjps",              // Missing description for fjmpjps
	Port{1873, 17}:   "fjmpjps",              // Fjmpjps
	Port{1874, 6}:    "fjswapsnp",            // Missing description for fjswapsnp
	Port{1874, 17}:   "fjswapsnp",            // Fjswapsnp
	Port{1875, 6}:    "westell-stats",        // westell stats
	Port{1875, 17}:   "westell-stats",        // westell stats
	Port{1876, 6}:    "ewcappsrv",            // Missing description for ewcappsrv
	Port{1876, 17}:   "ewcappsrv",            // Missing description for ewcappsrv
	Port{1877, 6}:    "hp-webqosdb",          // Missing description for hp-webqosdb
	Port{1877, 17}:   "hp-webqosdb",          // Missing description for hp-webqosdb
	Port{1878, 6}:    "drmsmc",               // Missing description for drmsmc
	Port{1878, 17}:   "drmsmc",               // Missing description for drmsmc
	Port{1879, 6}:    "nettgain-nms",         // NettGain NMS
	Port{1879, 17}:   "nettgain-nms",         // NettGain NMS
	Port{1880, 6}:    "vsat-control",         // Gilat VSAT Control
	Port{1880, 17}:   "vsat-control",         // Gilat VSAT Control
	Port{1881, 6}:    "ibm-mqseries2",        // IBM WebSphere MQ Everyplace
	Port{1881, 17}:   "ibm-mqseries2",        // IBM WebSphere MQ Everyplace
	Port{1882, 6}:    "ecsqdmn",              // CA eTrust Common Services
	Port{1882, 17}:   "ecsqdmn",              // CA eTrust Common Services
	Port{1883, 6}:    "mqtt",                 // Message Queuing Telemetry Transport Protocol
	Port{1883, 17}:   "ibm-mqisdp",           // IBM MQSeries SCADA
	Port{1884, 6}:    "idmaps",               // Internet Distance Map Svc
	Port{1884, 17}:   "idmaps",               // Internet Distance Map Svc
	Port{1885, 6}:    "vrtstrapserver",       // Veritas Trap Server
	Port{1885, 17}:   "vrtstrapserver",       // Veritas Trap Server
	Port{1886, 6}:    "leoip",                // Leonardo over IP
	Port{1886, 17}:   "leoip",                // Leonardo over IP
	Port{1887, 6}:    "filex-lport",          // FileX Listening Port
	Port{1887, 17}:   "filex-lport",          // FileX Listening Port
	Port{1888, 6}:    "ncconfig",             // NC Config Port
	Port{1888, 17}:   "ncconfig",             // NC Config Port
	Port{1889, 6}:    "unify-adapter",        // Unify Web Adapter Service
	Port{1889, 17}:   "unify-adapter",        // Unify Web Adapter Service
	Port{1890, 6}:    "wilkenlistener",       // Missing description for wilkenlistener
	Port{1890, 17}:   "wilkenlistener",       // wilkenListener
	Port{1891, 6}:    "childkey-notif",       // ChildKey Notification
	Port{1891, 17}:   "childkey-notif",       // ChildKey Notification
	Port{1892, 6}:    "childkey-ctrl",        // ChildKey Control
	Port{1892, 17}:   "childkey-ctrl",        // ChildKey Control
	Port{1893, 6}:    "elad",                 // ELAD Protocol
	Port{1893, 17}:   "elad",                 // ELAD Protocol
	Port{1894, 6}:    "o2server-port",        // O2Server Port
	Port{1894, 17}:   "o2server-port",        // O2Server Port
	Port{1896, 6}:    "b-novative-ls",        // b-novative license server
	Port{1896, 17}:   "b-novative-ls",        // b-novative license server
	Port{1897, 6}:    "metaagent",            // Missing description for metaagent
	Port{1897, 17}:   "metaagent",            // MetaAgent
	Port{1898, 6}:    "cymtec-port",          // Cymtec secure management
	Port{1898, 17}:   "cymtec-port",          // Cymtec secure management
	Port{1899, 6}:    "mc2studios",           // Missing description for mc2studios
	Port{1899, 17}:   "mc2studios",           // MC2Studios
	Port{1900, 6}:    "upnp",                 // ssdp | Universal PnP | SSDP
	Port{1900, 17}:   "upnp",                 // Universal PnP
	Port{1901, 6}:    "fjicl-tep-a",          // Fujitsu ICL Terminal Emulator Program A
	Port{1901, 17}:   "fjicl-tep-a",          // Fujitsu ICL Terminal Emulator Program A
	Port{1902, 6}:    "fjicl-tep-b",          // Fujitsu ICL Terminal Emulator Program B
	Port{1902, 17}:   "fjicl-tep-b",          // Fujitsu ICL Terminal Emulator Program B
	Port{1903, 6}:    "linkname",             // Local Link Name Resolution
	Port{1903, 17}:   "linkname",             // Local Link Name Resolution
	Port{1904, 6}:    "fjicl-tep-c",          // Fujitsu ICL Terminal Emulator Program C
	Port{1904, 17}:   "fjicl-tep-c",          // Fujitsu ICL Terminal Emulator Program C
	Port{1905, 6}:    "sugp",                 // Secure UP.Link Gateway Protocol
	Port{1905, 17}:   "sugp",                 // Secure UP.Link Gateway Protocol
	Port{1906, 6}:    "tpmd",                 // TPortMapperReq
	Port{1906, 17}:   "tpmd",                 // TPortMapperReq
	Port{1907, 6}:    "intrastar",            // Missing description for intrastar
	Port{1907, 17}:   "intrastar",            // IntraSTAR
	Port{1908, 6}:    "dawn",                 // Missing description for dawn
	Port{1908, 17}:   "dawn",                 // Dawn
	Port{1909, 6}:    "global-wlink",         // Global World Link
	Port{1909, 17}:   "global-wlink",         // Global World Link
	Port{1910, 6}:    "ultrabac",             // UltraBac Software communications port
	Port{1910, 17}:   "ultrabac",             // UltraBac Software communications port
	Port{1911, 6}:    "mtp",                  // Starlight Networks Multimedia Transport Protocol
	Port{1911, 17}:   "mtp",                  // Starlight Networks Multimedia Transport Protocol
	Port{1912, 6}:    "rhp-iibp",             // Missing description for rhp-iibp
	Port{1912, 17}:   "rhp-iibp",             // Missing description for rhp-iibp
	Port{1913, 6}:    "armadp",               // Missing description for armadp
	Port{1913, 17}:   "armadp",               // Missing description for armadp
	Port{1914, 6}:    "elm-momentum",         // Missing description for elm-momentum
	Port{1914, 17}:   "elm-momentum",         // Elm-Momentum
	Port{1915, 6}:    "facelink",             // Missing description for facelink
	Port{1915, 17}:   "facelink",             // FACELINK
	Port{1916, 6}:    "persona",              // Persoft Persona
	Port{1916, 17}:   "persona",              // Persoft Persona
	Port{1917, 6}:    "noagent",              // Missing description for noagent
	Port{1917, 17}:   "noagent",              // nOAgent
	Port{1918, 6}:    "can-nds",              // IBM Tivole Directory Service - NDS
	Port{1918, 17}:   "can-nds",              // IBM Tivole Directory Service - NDS
	Port{1919, 6}:    "can-dch",              // IBM Tivoli Directory Service - DCH
	Port{1919, 17}:   "can-dch",              // IBM Tivoli Directory Service - DCH
	Port{1920, 6}:    "can-ferret",           // IBM Tivoli Directory Service - FERRET
	Port{1920, 17}:   "can-ferret",           // IBM Tivoli Directory Service - FERRET
	Port{1921, 6}:    "noadmin",              // Missing description for noadmin
	Port{1921, 17}:   "noadmin",              // NoAdmin
	Port{1922, 6}:    "tapestry",             // Missing description for tapestry
	Port{1922, 17}:   "tapestry",             // Tapestry
	Port{1923, 6}:    "spice",                // Missing description for spice
	Port{1923, 17}:   "spice",                // SPICE
	Port{1924, 6}:    "xiip",                 // Missing description for xiip
	Port{1924, 17}:   "xiip",                 // XIIP
	Port{1925, 6}:    "discovery-port",       // Surrogate Discovery Port
	Port{1925, 17}:   "discovery-port",       // Surrogate Discovery Port
	Port{1926, 6}:    "egs",                  // Evolution Game Server
	Port{1926, 17}:   "egs",                  // Evolution Game Server
	Port{1927, 6}:    "videte-cipc",          // Videte CIPC Port
	Port{1927, 17}:   "videte-cipc",          // Videte CIPC Port
	Port{1928, 6}:    "emsd-port",            // Expnd Maui Srvr Dscovr
	Port{1928, 17}:   "emsd-port",            // Expnd Maui Srvr Dscovr
	Port{1929, 6}:    "bandwiz-system",       // Bandwiz System - Server
	Port{1929, 17}:   "bandwiz-system",       // Bandwiz System - Server
	Port{1930, 6}:    "driveappserver",       // Drive AppServer
	Port{1930, 17}:   "driveappserver",       // Drive AppServer
	Port{1931, 6}:    "amdsched",             // AMD SCHED
	Port{1931, 17}:   "amdsched",             // AMD SCHED
	Port{1932, 6}:    "ctt-broker",           // CTT Broker
	Port{1932, 17}:   "ctt-broker",           // CTT Broker
	Port{1933, 6}:    "xmapi",                // IBM LM MT Agent
	Port{1933, 17}:   "xmapi",                // IBM LM MT Agent
	Port{1934, 6}:    "xaapi",                // IBM LM Appl Agent
	Port{1934, 17}:   "xaapi",                // IBM LM Appl Agent
	Port{1935, 6}:    "rtmp",                 // macromedia-fcs | Macromedia FlasComm Server | Macromedia Flash Communications Server MX | Macromedia Flash Communications server MX
	Port{1935, 17}:   "macromedia-fcs",       // Macromedia Flash Communications server MX
	Port{1936, 6}:    "jetcmeserver",         // JetCmeServer Server Port
	Port{1936, 17}:   "jetcmeserver",         // JetCmeServer Server Port
	Port{1937, 6}:    "jwserver",             // JetVWay Server Port
	Port{1937, 17}:   "jwserver",             // JetVWay Server Port
	Port{1938, 6}:    "jwclient",             // JetVWay Client Port
	Port{1938, 17}:   "jwclient",             // JetVWay Client Port
	Port{1939, 6}:    "jvserver",             // JetVision Server Port
	Port{1939, 17}:   "jvserver",             // JetVision Server Port
	Port{1940, 6}:    "jvclient",             // JetVision Client Port
	Port{1940, 17}:   "jvclient",             // JetVision Client Port
	Port{1941, 6}:    "dic-aida",             // Missing description for dic-aida
	Port{1941, 17}:   "dic-aida",             // DIC-Aida
	Port{1942, 6}:    "res",                  // Real Enterprise Service
	Port{1942, 17}:   "res",                  // Real Enterprise Service
	Port{1943, 6}:    "beeyond-media",        // Beeyond Media
	Port{1943, 17}:   "beeyond-media",        // Beeyond Media
	Port{1944, 6}:    "close-combat",         // Missing description for close-combat
	Port{1944, 17}:   "close-combat",         // Missing description for close-combat
	Port{1945, 6}:    "dialogic-elmd",        // Missing description for dialogic-elmd
	Port{1945, 17}:   "dialogic-elmd",        // Missing description for dialogic-elmd
	Port{1946, 6}:    "tekpls",               // Missing description for tekpls
	Port{1946, 17}:   "tekpls",               // Missing description for tekpls
	Port{1947, 6}:    "sentinelsrm",          // Missing description for sentinelsrm
	Port{1947, 17}:   "sentinelsrm",          // SentinelSRM
	Port{1948, 6}:    "eye2eye",              // Missing description for eye2eye
	Port{1948, 17}:   "eye2eye",              // Missing description for eye2eye
	Port{1949, 6}:    "ismaeasdaqlive",       // ISMA Easdaq Live
	Port{1949, 17}:   "ismaeasdaqlive",       // ISMA Easdaq Live
	Port{1950, 6}:    "ismaeasdaqtest",       // ISMA Easdaq Test
	Port{1950, 17}:   "ismaeasdaqtest",       // ISMA Easdaq Test
	Port{1951, 6}:    "bcs-lmserver",         // Missing description for bcs-lmserver
	Port{1951, 17}:   "bcs-lmserver",         // Missing description for bcs-lmserver
	Port{1952, 6}:    "mpnjsc",               // Missing description for mpnjsc
	Port{1952, 17}:   "mpnjsc",               // Missing description for mpnjsc
	Port{1953, 6}:    "rapidbase",            // Rapid Base
	Port{1953, 17}:   "rapidbase",            // Rapid Base
	Port{1954, 6}:    "abr-api",              // ABR-API (diskbridge)
	Port{1954, 17}:   "abr-api",              // ABR-API (diskbridge)
	Port{1955, 6}:    "abr-secure",           // ABR-Secure Data (diskbridge)
	Port{1955, 17}:   "abr-secure",           // ABR-Secure Data (diskbridge)
	Port{1956, 6}:    "vrtl-vmf-ds",          // Vertel VMF DS
	Port{1956, 17}:   "vrtl-vmf-ds",          // Vertel VMF DS
	Port{1957, 6}:    "unix-status",          // Missing description for unix-status
	Port{1957, 17}:   "unix-status",          // Missing description for unix-status
	Port{1958, 6}:    "dxadmind",             // CA Administration Daemon
	Port{1958, 17}:   "dxadmind",             // CA Administration Daemon
	Port{1959, 6}:    "simp-all",             // SIMP Channel
	Port{1959, 17}:   "simp-all",             // SIMP Channel
	Port{1960, 6}:    "nasmanager",           // Merit DAC NASmanager
	Port{1960, 17}:   "nasmanager",           // Merit DAC NASmanager
	Port{1961, 6}:    "bts-appserver",        // BTS APPSERVER
	Port{1961, 17}:   "bts-appserver",        // BTS APPSERVER
	Port{1962, 6}:    "biap-mp",              // Missing description for biap-mp
	Port{1962, 17}:   "biap-mp",              // BIAP-MP
	Port{1963, 6}:    "webmachine",           // Missing description for webmachine
	Port{1963, 17}:   "webmachine",           // WebMachine
	Port{1964, 6}:    "solid-e-engine",       // SOLID E ENGINE
	Port{1964, 17}:   "solid-e-engine",       // SOLID E ENGINE
	Port{1965, 6}:    "tivoli-npm",           // Tivoli NPM
	Port{1965, 17}:   "tivoli-npm",           // Tivoli NPM
	Port{1966, 6}:    "slush",                // Missing description for slush
	Port{1966, 17}:   "slush",                // Slush
	Port{1967, 6}:    "sns-quote",            // SNS Quote
	Port{1967, 17}:   "sns-quote",            // SNS Quote
	Port{1968, 6}:    "lipsinc",              // Missing description for lipsinc
	Port{1968, 17}:   "lipsinc",              // LIPSinc
	Port{1969, 6}:    "lipsinc1",             // LIPSinc 1
	Port{1969, 17}:   "lipsinc1",             // LIPSinc 1
	Port{1970, 6}:    "netop-rc",             // NetOp Remote Control
	Port{1970, 17}:   "netop-rc",             // NetOp Remote Control
	Port{1971, 6}:    "netop-school",         // NetOp School
	Port{1971, 17}:   "netop-school",         // NetOp School
	Port{1972, 6}:    "intersys-cache",       // Cache
	Port{1972, 17}:   "intersys-cache",       // Cache
	Port{1973, 6}:    "dlsrap",               // Data Link Switching Remote Access Protocol
	Port{1973, 17}:   "dlsrap",               // Data Link Switching Remote Access Protocol
	Port{1974, 6}:    "drp",                  // Missing description for drp
	Port{1974, 17}:   "drp",                  // DRP
	Port{1975, 6}:    "tcoflashagent",        // TCO Flash Agent
	Port{1975, 17}:   "tcoflashagent",        // TCO Flash Agent
	Port{1976, 6}:    "tcoregagent",          // TCO Reg Agent
	Port{1976, 17}:   "tcoregagent",          // TCO Reg Agent
	Port{1977, 6}:    "tcoaddressbook",       // TCO Address Book
	Port{1977, 17}:   "tcoaddressbook",       // TCO Address Book
	Port{1978, 6}:    "unisql",               // Missing description for unisql
	Port{1978, 17}:   "unisql",               // UniSQL
	Port{1979, 6}:    "unisql-java",          // UniSQL Java
	Port{1979, 17}:   "unisql-java",          // UniSQL Java
	Port{1980, 6}:    "pearldoc-xact",        // PearlDoc XACT
	Port{1980, 17}:   "pearldoc-xact",        // PearlDoc XACT
	Port{1981, 6}:    "p2pq",                 // Missing description for p2pq
	Port{1981, 17}:   "p2pq",                 // p2pQ
	Port{1982, 6}:    "estamp",               // Evidentiary Timestamp
	Port{1982, 17}:   "estamp",               // Evidentiary Timestamp
	Port{1983, 6}:    "lhtp",                 // Loophole Test Protocol
	Port{1983, 17}:   "lhtp",                 // Loophole Test Protocol
	Port{1984, 6}:    "bigbrother",           // bb | Big Brother monitoring server - www.bb4.com | BB
	Port{1984, 17}:   "bb",                   // BB
	Port{1985, 6}:    "hsrp",                 // Hot Standby Router Protocol
	Port{1985, 17}:   "hsrp",                 // Hot Standby Router Protocol
	Port{1986, 6}:    "licensedaemon",        // cisco license management
	Port{1986, 17}:   "licensedaemon",        // cisco license management
	Port{1987, 6}:    "tr-rsrb-p1",           // cisco RSRB Priority 1 port
	Port{1987, 17}:   "tr-rsrb-p1",           // cisco RSRB Priority 1 port
	Port{1988, 6}:    "tr-rsrb-p2",           // cisco RSRB Priority 2 port
	Port{1988, 17}:   "tr-rsrb-p2",           // cisco RSRB Priority 2 port
	Port{1989, 6}:    "tr-rsrb-p3",           // mshnet | cisco RSRB Priority 3 port | MHSnet system
	Port{1989, 17}:   "tr-rsrb-p3",           // cisco RSRB Priority 3 port
	Port{1990, 6}:    "stun-p1",              // cisco STUN Priority 1 port
	Port{1990, 17}:   "stun-p1",              // cisco STUN Priority 1 port
	Port{1991, 6}:    "stun-p2",              // cisco STUN Priority 2 port
	Port{1991, 17}:   "stun-p2",              // cisco STUN Priority 2 port
	Port{1992, 6}:    "stun-p3",              // ipsendmsg | cisco STUN Priority 3 port | IPsendmsg
	Port{1992, 17}:   "stun-p3",              // cisco STUN Priority 3 port
	Port{1993, 6}:    "snmp-tcp-port",        // cisco SNMP TCP port
	Port{1993, 17}:   "snmp-tcp-port",        // cisco SNMP TCP port
	Port{1994, 6}:    "stun-port",            // cisco serial tunnel port
	Port{1994, 17}:   "stun-port",            // cisco serial tunnel port
	Port{1995, 6}:    "perf-port",            // cisco perf port
	Port{1995, 17}:   "perf-port",            // cisco perf port
	Port{1996, 6}:    "tr-rsrb-port",         // cisco Remote SRB port
	Port{1996, 17}:   "tr-rsrb-port",         // cisco Remote SRB port
	Port{1997, 6}:    "gdp-port",             // cisco Gateway Discovery Protocol
	Port{1997, 17}:   "gdp-port",             // cisco Gateway Discovery Protocol
	Port{1998, 6}:    "x25-svc-port",         // cisco X.25 service (XOT)
	Port{1998, 17}:   "x25-svc-port",         // cisco X.25 service (XOT)
	Port{1999, 6}:    "tcp-id-port",          // cisco identification port
	Port{1999, 17}:   "tcp-id-port",          // cisco identification port
	Port{2000, 6}:    "cisco-sccp",           // cisco SCCP (Skinny Client Control Protocol) | Cisco SCCP | Cisco SCCp
	Port{2000, 17}:   "cisco-sccp",           // cisco SCCP (Skinny Client Control Protocol)
	Port{2001, 6}:    "dc",                   // wizard | or nfr20 web queries | curry
	Port{2001, 17}:   "wizard",               // curry
	Port{2002, 6}:    "globe",                // Missing description for globe
	Port{2002, 17}:   "globe",                // Missing description for globe
	Port{2003, 6}:    "finger",               // brutus | GNU finger (cfingerd) | Brutus Server
	Port{2003, 17}:   "brutus",               // Brutus Server
	Port{2004, 6}:    "mailbox",              // emce | CCWS mm conf
	Port{2004, 17}:   "emce",                 // CCWS mm conf
	Port{2005, 6}:    "deslogin",             // oracle | berknet | encrypted symmetric telnet login
	Port{2005, 17}:   "oracle",               // Missing description for oracle
	Port{2006, 6}:    "invokator",            // raid-cd | raid
	Port{2006, 17}:   "raid-cc",              // raid
	Port{2007, 6}:    "dectalk",              // raid-am
	Port{2007, 17}:   "raid-am",              // Missing description for raid-am
	Port{2008, 6}:    "conf",                 // terminaldb
	Port{2008, 17}:   "terminaldb",           // Missing description for terminaldb
	Port{2009, 6}:    "news",                 // whosockami
	Port{2009, 17}:   "whosockami",           // Missing description for whosockami
	Port{2010, 6}:    "search",               // pipe_server | pipe-server | Or nfr411
	Port{2010, 17}:   "pipe_server",          // Also used by NFR
	Port{2011, 6}:    "raid-cc",              // servserv | raid
	Port{2011, 17}:   "servserv",             // Missing description for servserv
	Port{2012, 6}:    "ttyinfo",              // raid-ac
	Port{2012, 17}:   "raid-ac",              // Missing description for raid-ac
	Port{2013, 6}:    "raid-am",              // raid-cd
	Port{2013, 17}:   "raid-cd",              // Missing description for raid-cd
	Port{2014, 6}:    "troff",                // raid-sf
	Port{2014, 17}:   "raid-sf",              // Missing description for raid-sf
	Port{2015, 6}:    "cypress",              // raid-cs
	Port{2015, 17}:   "raid-cs",              // Missing description for raid-cs
	Port{2016, 6}:    "bootserver",           // Missing description for bootserver
	Port{2016, 17}:   "bootserver",           // Missing description for bootserver
	Port{2017, 6}:    "cypress-stat",         // bootclient
	Port{2017, 17}:   "bootclient",           // Missing description for bootclient
	Port{2018, 6}:    "terminaldb",           // rellpack
	Port{2018, 17}:   "rellpack",             // Missing description for rellpack
	Port{2019, 6}:    "whosockami",           // about
	Port{2019, 17}:   "about",                // Missing description for about
	Port{2020, 6}:    "xinupageserver",       // Missing description for xinupageserver
	Port{2020, 17}:   "xinupageserver",       // Missing description for xinupageserver
	Port{2021, 6}:    "servexec",             // xinuexpansion1
	Port{2021, 17}:   "xinuexpansion1",       // Missing description for xinuexpansion1
	Port{2022, 6}:    "down",                 // xinuexpansion2
	Port{2022, 17}:   "xinuexpansion2",       // Missing description for xinuexpansion2
	Port{2023, 6}:    "xinuexpansion3",       // Missing description for xinuexpansion3
	Port{2023, 17}:   "xinuexpansion3",       // Missing description for xinuexpansion3
	Port{2024, 6}:    "xinuexpansion4",       // Missing description for xinuexpansion4
	Port{2024, 17}:   "xinuexpansion4",       // Missing description for xinuexpansion4
	Port{2025, 6}:    "ellpack",              // xribs
	Port{2025, 17}:   "xribs",                // Missing description for xribs
	Port{2026, 6}:    "scrabble",             // Missing description for scrabble
	Port{2026, 17}:   "scrabble",             // Missing description for scrabble
	Port{2027, 6}:    "shadowserver",         // Missing description for shadowserver
	Port{2027, 17}:   "shadowserver",         // Missing description for shadowserver
	Port{2028, 6}:    "submitserver",         // Missing description for submitserver
	Port{2028, 17}:   "submitserver",         // Missing description for submitserver
	Port{2029, 6}:    "hsrpv6",               // Hot Standby Router Protocol IPv6
	Port{2029, 17}:   "hsrpv6",               // Hot Standby Router Protocol IPv6
	Port{2030, 6}:    "device2",              // Missing description for device2
	Port{2030, 17}:   "device2",              // Missing description for device2
	Port{2031, 6}:    "mobrien-chat",         // Missing description for mobrien-chat
	Port{2031, 17}:   "mobrien-chat",         // Missing description for mobrien-chat
	Port{2032, 6}:    "blackboard",           // Missing description for blackboard
	Port{2032, 17}:   "blackboard",           // Missing description for blackboard
	Port{2033, 6}:    "glogger",              // Missing description for glogger
	Port{2033, 17}:   "glogger",              // Missing description for glogger
	Port{2034, 6}:    "scoremgr",             // Missing description for scoremgr
	Port{2034, 17}:   "scoremgr",             // Missing description for scoremgr
	Port{2035, 6}:    "imsldoc",              // Missing description for imsldoc
	Port{2035, 17}:   "imsldoc",              // Missing description for imsldoc
	Port{2036, 6}:    "e-dpnet",              // Ethernet WS DP network
	Port{2036, 17}:   "e-dpnet",              // Ethernet WS DP network
	Port{2037, 6}:    "applus",               // APplus Application Server
	Port{2037, 17}:   "applus",               // APplus Application Server
	Port{2038, 6}:    "objectmanager",        // Missing description for objectmanager
	Port{2038, 17}:   "objectmanager",        // Missing description for objectmanager
	Port{2039, 6}:    "prizma",               // Prizma Monitoring Service
	Port{2039, 17}:   "prizma",               // Prizma Monitoring Service
	Port{2040, 6}:    "lam",                  // Missing description for lam
	Port{2040, 17}:   "lam",                  // Missing description for lam
	Port{2041, 6}:    "interbase",            // Missing description for interbase
	Port{2041, 17}:   "interbase",            // Missing description for interbase
	Port{2042, 6}:    "isis",                 // Missing description for isis
	Port{2042, 17}:   "isis",                 // Missing description for isis
	Port{2043, 6}:    "isis-bcast",           // Missing description for isis-bcast
	Port{2043, 17}:   "isis-bcast",           // Missing description for isis-bcast
	Port{2044, 6}:    "rimsl",                // Missing description for rimsl
	Port{2044, 17}:   "rimsl",                // Missing description for rimsl
	Port{2045, 6}:    "cdfunc",               // Missing description for cdfunc
	Port{2045, 17}:   "cdfunc",               // Missing description for cdfunc
	Port{2046, 6}:    "sdfunc",               // Missing description for sdfunc
	Port{2046, 17}:   "sdfunc",               // Missing description for sdfunc
	Port{2047, 6}:    "dls",                  // Missing description for dls
	Port{2047, 17}:   "dls",                  // Missing description for dls
	Port{2048, 6}:    "dls-monitor",          // Missing description for dls-monitor
	Port{2048, 17}:   "dls-monitor",          // Missing description for dls-monitor
	Port{2049, 132}:  "nfs",                  // shilp | Network File System | Network File System - Sun Microsystems
	Port{2049, 6}:    "nfs",                  // networked file system
	Port{2049, 17}:   "nfs",                  // networked file system
	Port{2050, 6}:    "av-emb-config",        // Avaya EMB Config Port
	Port{2050, 17}:   "av-emb-config",        // Avaya EMB Config Port
	Port{2051, 6}:    "epnsdp",               // Missing description for epnsdp
	Port{2051, 17}:   "epnsdp",               // EPNSDP
	Port{2052, 6}:    "clearvisn",            // clearVisn Services Port
	Port{2052, 17}:   "clearvisn",            // clearVisn Services Port
	Port{2053, 6}:    "knetd",                // lot105-ds-upd | Lot105 DSuper Updates
	Port{2053, 17}:   "lot105-ds-upd",        // Lot105 DSuper Updates
	Port{2054, 6}:    "weblogin",             // Weblogin Port
	Port{2054, 17}:   "weblogin",             // Weblogin Port
	Port{2055, 6}:    "iop",                  // Iliad-Odyssey Protocol
	Port{2055, 17}:   "iop",                  // Iliad-Odyssey Protocol
	Port{2056, 6}:    "omnisky",              // OmniSky Port
	Port{2056, 17}:   "omnisky",              // OmniSky Port
	Port{2057, 6}:    "rich-cp",              // Rich Content Protocol
	Port{2057, 17}:   "rich-cp",              // Rich Content Protocol
	Port{2058, 6}:    "newwavesearch",        // NewWaveSearchables RMI
	Port{2058, 17}:   "newwavesearch",        // NewWaveSearchables RMI
	Port{2059, 6}:    "bmc-messaging",        // BMC Messaging Service
	Port{2059, 17}:   "bmc-messaging",        // BMC Messaging Service
	Port{2060, 6}:    "teleniumdaemon",       // Telenium Daemon IF
	Port{2060, 17}:   "teleniumdaemon",       // Telenium Daemon IF
	Port{2061, 6}:    "netmount",             // Missing description for netmount
	Port{2061, 17}:   "netmount",             // NetMount
	Port{2062, 6}:    "icg-swp",              // ICG SWP Port
	Port{2062, 17}:   "icg-swp",              // ICG SWP Port
	Port{2063, 6}:    "icg-bridge",           // ICG Bridge Port
	Port{2063, 17}:   "icg-bridge",           // ICG Bridge Port
	Port{2064, 6}:    "dnet-keyproxy",        // icg-iprelay | A closed-source client for solving the RSA cryptographic challenge. This is the keyblock proxy port. | ICG IP Relay Port
	Port{2064, 17}:   "icg-iprelay",          // ICG IP Relay Port
	Port{2065, 6}:    "dlsrpn",               // Data Link Switch Read Port Number
	Port{2065, 17}:   "dlsrpn",               // Data Link Switch Read Port Number
	Port{2066, 6}:    "aura",                 // AVM USB Remote Architecture
	Port{2066, 17}:   "aura",                 // AVM USB Remote Architecture
	Port{2067, 6}:    "dlswpn",               // Data Link Switch Write Port Number
	Port{2067, 17}:   "dlswpn",               // Data Link Switch Write Port Number
	Port{2068, 6}:    "avocentkvm",           // avauthsrvprtcl | Avocent KVM Server | Avocent AuthSrv Protocol
	Port{2068, 17}:   "avauthsrvprtcl",       // Avocent AuthSrv Protocol
	Port{2069, 6}:    "event-port",           // HTTP Event Port
	Port{2069, 17}:   "event-port",           // HTTP Event Port
	Port{2070, 6}:    "ah-esp-encap",         // AH and ESP Encapsulated in UDP packet
	Port{2070, 17}:   "ah-esp-encap",         // AH and ESP Encapsulated in UDP packet
	Port{2071, 6}:    "acp-port",             // Axon Control Protocol
	Port{2071, 17}:   "acp-port",             // Axon Control Protocol
	Port{2072, 6}:    "msync",                // GlobeCast mSync
	Port{2072, 17}:   "msync",                // GlobeCast mSync
	Port{2073, 6}:    "gxs-data-port",        // DataReel Database Socket
	Port{2073, 17}:   "gxs-data-port",        // DataReel Database Socket
	Port{2074, 6}:    "vrtl-vmf-sa",          // Vertel VMF SA
	Port{2074, 17}:   "vrtl-vmf-sa",          // Vertel VMF SA
	Port{2075, 6}:    "newlixengine",         // Newlix ServerWare Engine
	Port{2075, 17}:   "newlixengine",         // Newlix ServerWare Engine
	Port{2076, 6}:    "newlixconfig",         // Newlix JSPConfig
	Port{2076, 17}:   "newlixconfig",         // Newlix JSPConfig
	Port{2077, 6}:    "tsrmagt",              // Old Tivoli Storage Manager
	Port{2077, 17}:   "tsrmagt",              // Old Tivoli Storage Manager
	Port{2078, 6}:    "tpcsrvr",              // IBM Total Productivity Center Server
	Port{2078, 17}:   "tpcsrvr",              // IBM Total Productivity Center Server
	Port{2079, 6}:    "idware-router",        // IDWARE Router Port
	Port{2079, 17}:   "idware-router",        // IDWARE Router Port
	Port{2080, 6}:    "autodesk-nlm",         // Autodesk NLM (FLEXlm)
	Port{2080, 17}:   "autodesk-nlm",         // Autodesk NLM (FLEXlm)
	Port{2081, 6}:    "kme-trap-port",        // KME PRINTER TRAP PORT
	Port{2081, 17}:   "kme-trap-port",        // KME PRINTER TRAP PORT
	Port{2082, 6}:    "infowave",             // Infowave Mobility Server
	Port{2082, 17}:   "infowave",             // Infowave Mobiltiy Server
	Port{2083, 6}:    "radsec",               // Secure Radius Service
	Port{2083, 17}:   "radsec",               // Secure Radius Service
	Port{2084, 6}:    "sunclustergeo",        // SunCluster Geographic
	Port{2084, 17}:   "sunclustergeo",        // SunCluster Geographic
	Port{2085, 6}:    "ada-cip",              // ADA Control
	Port{2085, 17}:   "ada-cip",              // ADA Control
	Port{2086, 6}:    "gnunet",               // Missing description for gnunet
	Port{2086, 17}:   "gnunet",               // GNUnet
	Port{2087, 6}:    "eli",                  // ELI - Event Logging Integration
	Port{2087, 17}:   "eli",                  // ELI - Event Logging Integration
	Port{2088, 6}:    "ip-blf",               // IP Busy Lamp Field
	Port{2088, 17}:   "ip-blf",               // IP Busy Lamp Field
	Port{2089, 6}:    "sep",                  // Security Encapsulation Protocol - SEP
	Port{2089, 17}:   "sep",                  // Security Encapsulation Protocol - SEP
	Port{2090, 6}:    "lrp",                  // Load Report Protocol
	Port{2090, 17}:   "lrp",                  // Load Report Protocol
	Port{2091, 6}:    "prp",                  // Missing description for prp
	Port{2091, 17}:   "prp",                  // PRP
	Port{2092, 6}:    "descent3",             // Descent 3
	Port{2092, 17}:   "descent3",             // Descent 3
	Port{2093, 6}:    "nbx-cc",               // NBX CC
	Port{2093, 17}:   "nbx-cc",               // NBX CC
	Port{2094, 6}:    "nbx-au",               // NBX AU
	Port{2094, 17}:   "nbx-au",               // NBX AU
	Port{2095, 6}:    "nbx-ser",              // NBX SER
	Port{2095, 17}:   "nbx-ser",              // NBX SER
	Port{2096, 6}:    "nbx-dir",              // NBX DIR
	Port{2096, 17}:   "nbx-dir",              // NBX DIR
	Port{2097, 6}:    "jetformpreview",       // Jet Form Preview
	Port{2097, 17}:   "jetformpreview",       // Jet Form Preview
	Port{2098, 6}:    "dialog-port",          // Dialog Port
	Port{2098, 17}:   "dialog-port",          // Dialog Port
	Port{2099, 6}:    "h2250-annex-g",        // H.225.0 Annex G | H.225.0 Annex G Signalling
	Port{2099, 17}:   "h2250-annex-g",        // H.225.0 Annex G
	Port{2100, 6}:    "amiganetfs",           // Amiga Network Filesystem
	Port{2100, 17}:   "amiganetfs",           // Amiga Network Filesystem
	Port{2101, 6}:    "rtcm-sc104",           // Missing description for rtcm-sc104
	Port{2101, 17}:   "rtcm-sc104",           // Missing description for rtcm-sc104
	Port{2102, 6}:    "zephyr-srv",           // Zephyr server
	Port{2102, 17}:   "zephyr-srv",           // Zephyr server
	Port{2103, 6}:    "zephyr-clt",           // Zephyr serv-hm connection
	Port{2103, 17}:   "zephyr-clt",           // Zephyr serv-hm connection
	Port{2104, 6}:    "zephyr-hm",            // Zephyr hostmanager
	Port{2104, 17}:   "zephyr-hm",            // Zephyr hostmanager
	Port{2105, 6}:    "eklogin",              // minipay | Kerberos (v4) encrypted rlogin | MiniPay
	Port{2105, 17}:   "eklogin",              // Kerberos (v4) encrypted rlogin
	Port{2106, 6}:    "ekshell",              // mzap | Kerberos (v4) encrypted rshell | MZAP
	Port{2106, 17}:   "ekshell",              // Kerberos (v4) encrypted rshell
	Port{2107, 6}:    "msmq-mgmt",            // bintec-admin | Microsoft Message Queuing (IANA calls this bintec-admin) | BinTec Admin
	Port{2107, 17}:   "bintec-admin",         // BinTec Admin
	Port{2108, 6}:    "rkinit",               // comcam | Kerberos (v4) remote initialization | Comcam
	Port{2108, 17}:   "rkinit",               // Kerberos (v4) remote initialization
	Port{2109, 6}:    "ergolight",            // Missing description for ergolight
	Port{2109, 17}:   "ergolight",            // Ergolight
	Port{2110, 6}:    "umsp",                 // Missing description for umsp
	Port{2110, 17}:   "umsp",                 // UMSP
	Port{2111, 6}:    "kx",                   // dsatp | X over kerberos | OPNET Dynamic Sampling Agent Transaction Protocol
	Port{2111, 17}:   "dsatp",                // DSATP
	Port{2112, 6}:    "kip",                  // idonix-metanet | IP over kerberos | Idonix MetaNet
	Port{2112, 17}:   "idonix-metanet",       // Idonix MetaNet
	Port{2113, 6}:    "hsl-storm",            // HSL StoRM
	Port{2113, 17}:   "hsl-storm",            // HSL StoRM
	Port{2114, 6}:    "newheights",           // ariascribe | Classical Music Meta-Data Access and Enhancement
	Port{2114, 17}:   "newheights",           // NEWHEIGHTS
	Port{2115, 6}:    "kdm",                  // Key Distribution Manager
	Port{2115, 17}:   "kdm",                  // Key Distribution Manager
	Port{2116, 6}:    "ccowcmr",              // Missing description for ccowcmr
	Port{2116, 17}:   "ccowcmr",              // CCOWCMR
	Port{2117, 6}:    "mentaclient",          // Missing description for mentaclient
	Port{2117, 17}:   "mentaclient",          // MENTACLIENT
	Port{2118, 6}:    "mentaserver",          // Missing description for mentaserver
	Port{2118, 17}:   "mentaserver",          // MENTASERVER
	Port{2119, 6}:    "gsigatekeeper",        // Missing description for gsigatekeeper
	Port{2119, 17}:   "gsigatekeeper",        // GSIGATEKEEPER
	Port{2120, 6}:    "kauth",                // qencp | Remote kauth | Quick Eagle Networks CP
	Port{2120, 17}:   "qencp",                // Quick Eagle Networks CP
	Port{2121, 6}:    "ccproxy-ftp",          // scientia-ssdb | CCProxy FTP Proxy | SCIENTIA-SSDB
	Port{2121, 17}:   "scientia-ssdb",        // SCIENTIA-SSDB
	Port{2122, 6}:    "caupc-remote",         // CauPC Remote Control
	Port{2122, 17}:   "caupc-remote",         // CauPC Remote Control
	Port{2123, 6}:    "gtp-control",          // GTP-Control Plane (3GPP)
	Port{2123, 17}:   "gtp-control",          // GTP-Control Plane (3GPP)
	Port{2124, 6}:    "elatelink",            // Missing description for elatelink
	Port{2124, 17}:   "elatelink",            // ELATELINK
	Port{2125, 6}:    "lockstep",             // Missing description for lockstep
	Port{2125, 17}:   "lockstep",             // LOCKSTEP
	Port{2126, 6}:    "pktcable-cops",        // Missing description for pktcable-cops
	Port{2126, 17}:   "pktcable-cops",        // PktCable-COPS
	Port{2127, 6}:    "index-pc-wb",          // Missing description for index-pc-wb
	Port{2127, 17}:   "index-pc-wb",          // INDEX-PC-WB
	Port{2128, 6}:    "net-steward",          // Net Steward Control
	Port{2128, 17}:   "net-steward",          // Net Steward Control
	Port{2129, 6}:    "cs-live",              // cs-live.com
	Port{2129, 17}:   "cs-live",              // cs-live.com
	Port{2130, 6}:    "xds",                  // Missing description for xds
	Port{2130, 17}:   "xds",                  // XDS
	Port{2131, 6}:    "avantageb2b",          // Missing description for avantageb2b
	Port{2131, 17}:   "avantageb2b",          // Avantageb2b
	Port{2132, 6}:    "solera-epmap",         // SoleraTec End Point Map
	Port{2132, 17}:   "solera-epmap",         // SoleraTec End Point Map
	Port{2133, 6}:    "zymed-zpp",            // Missing description for zymed-zpp
	Port{2133, 17}:   "zymed-zpp",            // ZYMED-ZPP
	Port{2134, 6}:    "avenue",               // Missing description for avenue
	Port{2134, 17}:   "avenue",               // AVENUE
	Port{2135, 6}:    "gris",                 // Grid Resource Information Server
	Port{2135, 17}:   "gris",                 // Grid Resource Information Server
	Port{2136, 6}:    "appworxsrv",           // Missing description for appworxsrv
	Port{2136, 17}:   "appworxsrv",           // APPWORXSRV
	Port{2137, 6}:    "connect",              // Missing description for connect
	Port{2137, 17}:   "connect",              // CONNECT
	Port{2138, 6}:    "unbind-cluster",       // Missing description for unbind-cluster
	Port{2138, 17}:   "unbind-cluster",       // UNBIND-CLUSTER
	Port{2139, 6}:    "ias-auth",             // Missing description for ias-auth
	Port{2139, 17}:   "ias-auth",             // IAS-AUTH
	Port{2140, 6}:    "ias-reg",              // Missing description for ias-reg
	Port{2140, 17}:   "ias-reg",              // IAS-REG
	Port{2141, 6}:    "ias-admind",           // Missing description for ias-admind
	Port{2141, 17}:   "ias-admind",           // IAS-ADMIND
	Port{2142, 6}:    "tdmoip",               // TDM OVER IP
	Port{2142, 17}:   "tdmoip",               // TDM OVER IP
	Port{2143, 6}:    "lv-jc",                // Live Vault Job Control
	Port{2143, 17}:   "lv-jc",                // Live Vault Job Control
	Port{2144, 6}:    "lv-ffx",               // Live Vault Fast Object Transfer
	Port{2144, 17}:   "lv-ffx",               // Live Vault Fast Object Transfer
	Port{2145, 6}:    "lv-pici",              // Live Vault Remote Diagnostic Console Support
	Port{2145, 17}:   "lv-pici",              // Live Vault Remote Diagnostic Console Support
	Port{2146, 6}:    "lv-not",               // Live Vault Admin Event Notification
	Port{2146, 17}:   "lv-not",               // Live Vault Admin Event Notification
	Port{2147, 6}:    "lv-auth",              // Live Vault Authentication
	Port{2147, 17}:   "lv-auth",              // Live Vault Authentication
	Port{2148, 6}:    "veritas-ucl",          // Veritas Universal Communication Layer | VERITAS UNIVERSAL COMMUNICATION LAYER
	Port{2148, 17}:   "veritas-ucl",          // Veritas Universal Communication Layer
	Port{2149, 6}:    "acptsys",              // Missing description for acptsys
	Port{2149, 17}:   "acptsys",              // ACPTSYS
	Port{2150, 6}:    "dynamic3d",            // Missing description for dynamic3d
	Port{2150, 17}:   "dynamic3d",            // DYNAMIC3D
	Port{2151, 6}:    "docent",               // Missing description for docent
	Port{2151, 17}:   "docent",               // DOCENT
	Port{2152, 6}:    "gtp-user",             // GTP-User Plane (3GPP)
	Port{2152, 17}:   "gtp-user",             // GTP-User Plane (3GPP)
	Port{2153, 6}:    "ctlptc",               // Control Protocol
	Port{2153, 17}:   "ctlptc",               // Control Protocol
	Port{2154, 6}:    "stdptc",               // Standard Protocol
	Port{2154, 17}:   "stdptc",               // Standard Protocol
	Port{2155, 6}:    "brdptc",               // Bridge Protocol
	Port{2155, 17}:   "brdptc",               // Bridge Protocol
	Port{2156, 6}:    "trp",                  // Talari Reliable Protocol
	Port{2156, 17}:   "trp",                  // Talari Reliable Protocol
	Port{2157, 6}:    "xnds",                 // Xerox Network Document Scan Protocol
	Port{2157, 17}:   "xnds",                 // Xerox Network Document Scan Protocol
	Port{2158, 6}:    "touchnetplus",         // TouchNetPlus Service
	Port{2158, 17}:   "touchnetplus",         // TouchNetPlus Service
	Port{2159, 6}:    "gdbremote",            // GDB Remote Debug Port
	Port{2159, 17}:   "gdbremote",            // GDB Remote Debug Port
	Port{2160, 6}:    "apc-2160",             // APC 2160
	Port{2160, 17}:   "apc-2160",             // APC 2160
	Port{2161, 6}:    "apc-agent",            // apc-2161 | American Power Conversion | APC 2161
	Port{2161, 17}:   "apc-2161",             // APC 2161
	Port{2162, 6}:    "navisphere",           // Missing description for navisphere
	Port{2162, 17}:   "navisphere",           // Navisphere
	Port{2163, 6}:    "navisphere-sec",       // Navisphere Secure
	Port{2163, 17}:   "navisphere-sec",       // Navisphere Secure
	Port{2164, 6}:    "ddns-v3",              // Dynamic DNS Version 3
	Port{2164, 17}:   "ddns-v3",              // Dynamic DNS Version 3
	Port{2165, 6}:    "x-bone-api",           // X-Bone API
	Port{2165, 17}:   "x-bone-api",           // X-Bone API
	Port{2166, 6}:    "iwserver",             // Missing description for iwserver
	Port{2166, 17}:   "iwserver",             // Missing description for iwserver
	Port{2167, 6}:    "raw-serial",           // Raw Async Serial Link
	Port{2167, 17}:   "raw-serial",           // Raw Async Serial Link
	Port{2168, 6}:    "easy-soft-mux",        // easy-soft Multiplexer
	Port{2168, 17}:   "easy-soft-mux",        // easy-soft Multiplexer
	Port{2169, 6}:    "brain",                // Backbone for Academic Information Notification (BRAIN)
	Port{2169, 17}:   "brain",                // Backbone for Academic Information Notification (BRAIN)
	Port{2170, 6}:    "eyetv",                // EyeTV Server Port
	Port{2170, 17}:   "eyetv",                // EyeTV Server Port
	Port{2171, 6}:    "msfw-storage",         // MS Firewall Storage
	Port{2171, 17}:   "msfw-storage",         // MS Firewall Storage
	Port{2172, 6}:    "msfw-s-storage",       // MS Firewall SecureStorage
	Port{2172, 17}:   "msfw-s-storage",       // MS Firewall SecureStorage
	Port{2173, 6}:    "msfw-replica",         // MS Firewall Replication
	Port{2173, 17}:   "msfw-replica",         // MS Firewall Replication
	Port{2174, 6}:    "msfw-array",           // MS Firewall Intra Array
	Port{2174, 17}:   "msfw-array",           // MS Firewall Intra Array
	Port{2175, 6}:    "airsync",              // Microsoft Desktop AirSync Protocol
	Port{2175, 17}:   "airsync",              // Microsoft Desktop AirSync Protocol
	Port{2176, 6}:    "rapi",                 // Microsoft ActiveSync Remote API
	Port{2176, 17}:   "rapi",                 // Microsoft ActiveSync Remote API
	Port{2177, 6}:    "qwave",                // qWAVE Bandwidth Estimate
	Port{2177, 17}:   "qwave",                // qWAVE Bandwidth Estimate
	Port{2178, 6}:    "bitspeer",             // Peer Services for BITS
	Port{2178, 17}:   "bitspeer",             // Peer Services for BITS
	Port{2179, 6}:    "vmrdp",                // Microsoft RDP for virtual machines
	Port{2179, 17}:   "vmrdp",                // Microsoft RDP for virtual machines
	Port{2180, 6}:    "mc-gt-srv",            // Millicent Vendor Gateway Server
	Port{2180, 17}:   "mc-gt-srv",            // Millicent Vendor Gateway Server
	Port{2181, 6}:    "eforward",             // Missing description for eforward
	Port{2181, 17}:   "eforward",             // Missing description for eforward
	Port{2182, 6}:    "cgn-stat",             // CGN status
	Port{2182, 17}:   "cgn-stat",             // CGN status
	Port{2183, 6}:    "cgn-config",           // Code Green configuration
	Port{2183, 17}:   "cgn-config",           // Code Green configuration
	Port{2184, 6}:    "nvd",                  // NVD User
	Port{2184, 17}:   "nvd",                  // NVD User
	Port{2185, 6}:    "onbase-dds",           // OnBase Distributed Disk Services
	Port{2185, 17}:   "onbase-dds",           // OnBase Distributed Disk Services
	Port{2186, 6}:    "gtaua",                // Guy-Tek Automated Update Applications
	Port{2186, 17}:   "gtaua",                // Guy-Tek Automated Update Applications
	Port{2187, 6}:    "ssmc",                 // ssmd | Sepehr System Management Control | Sepehr System Management Data
	Port{2187, 17}:   "ssmd",                 // Sepehr System Management Data
	Port{2188, 6}:    "radware-rpm",          // Radware Resource Pool Manager
	Port{2189, 6}:    "radware-rpm-s",        // Secure Radware Resource Pool Manager
	Port{2190, 6}:    "tivoconnect",          // TiVoConnect Beacon
	Port{2190, 17}:   "tivoconnect",          // TiVoConnect Beacon
	Port{2191, 6}:    "tvbus",                // TvBus Messaging
	Port{2191, 17}:   "tvbus",                // TvBus Messaging
	Port{2192, 6}:    "asdis",                // ASDIS software management
	Port{2192, 17}:   "asdis",                // ASDIS software management
	Port{2193, 6}:    "drwcs",                // Dr.Web Enterprise Management Service
	Port{2193, 17}:   "drwcs",                // Dr.Web Enterprise Management Service
	Port{2197, 6}:    "mnp-exchange",         // MNP data exchange
	Port{2197, 17}:   "mnp-exchange",         // MNP data exchange
	Port{2198, 6}:    "onehome-remote",       // OneHome Remote Access
	Port{2198, 17}:   "onehome-remote",       // OneHome Remote Access
	Port{2199, 6}:    "onehome-help",         // OneHome Service Port
	Port{2199, 17}:   "onehome-help",         // OneHome Service Port
	Port{2200, 6}:    "ici",                  // Missing description for ici
	Port{2200, 17}:   "ici",                  // ICI
	Port{2201, 6}:    "ats",                  // Advanced Training System Program
	Port{2201, 17}:   "ats",                  // Advanced Training System Program
	Port{2202, 6}:    "imtc-map",             // Int. Multimedia Teleconferencing Cosortium
	Port{2202, 17}:   "imtc-map",             // Int. Multimedia Teleconferencing Cosortium
	Port{2203, 6}:    "b2-runtime",           // b2 Runtime Protocol
	Port{2203, 17}:   "b2-runtime",           // b2 Runtime Protocol
	Port{2204, 6}:    "b2-license",           // b2 License Server
	Port{2204, 17}:   "b2-license",           // b2 License Server
	Port{2205, 6}:    "jps",                  // Java Presentation Server
	Port{2205, 17}:   "jps",                  // Java Presentation Server
	Port{2206, 6}:    "hpocbus",              // HP OpenCall bus
	Port{2206, 17}:   "hpocbus",              // HP OpenCall bus
	Port{2207, 6}:    "hpssd",                // HP Status and Services
	Port{2207, 17}:   "hpssd",                // HP Status and Services
	Port{2208, 6}:    "hpiod",                // HP I O Backend
	Port{2208, 17}:   "hpiod",                // HP I O Backend
	Port{2209, 6}:    "rimf-ps",              // HP RIM for Files Portal Service
	Port{2209, 17}:   "rimf-ps",              // HP RIM for Files Portal Service
	Port{2210, 6}:    "noaaport",             // NOAAPORT Broadcast Network
	Port{2210, 17}:   "noaaport",             // NOAAPORT Broadcast Network
	Port{2211, 6}:    "emwin",                // Missing description for emwin
	Port{2211, 17}:   "emwin",                // EMWIN
	Port{2212, 6}:    "leecoposserver",       // LeeCO POS Server Service
	Port{2212, 17}:   "leecoposserver",       // LeeCO POS Server Service
	Port{2213, 6}:    "kali",                 // Missing description for kali
	Port{2213, 17}:   "kali",                 // Kali
	Port{2214, 6}:    "rpi",                  // RDQ Protocol Interface
	Port{2214, 17}:   "rpi",                  // RDQ Protocol Interface
	Port{2215, 6}:    "ipcore",               // IPCore.co.za GPRS
	Port{2215, 17}:   "ipcore",               // IPCore.co.za GPRS
	Port{2216, 6}:    "vtu-comms",            // VTU data service
	Port{2216, 17}:   "vtu-comms",            // VTU data service
	Port{2217, 6}:    "gotodevice",           // GoToDevice Device Management
	Port{2217, 17}:   "gotodevice",           // GoToDevice Device Management
	Port{2218, 6}:    "bounzza",              // Bounzza IRC Proxy
	Port{2218, 17}:   "bounzza",              // Bounzza IRC Proxy
	Port{2219, 6}:    "netiq-ncap",           // NetIQ NCAP Protocol
	Port{2219, 17}:   "netiq-ncap",           // NetIQ NCAP Protocol
	Port{2220, 6}:    "netiq",                // NetIQ End2End
	Port{2220, 17}:   "netiq",                // NetIQ End2End
	Port{2221, 6}:    "rockwell-csp1",        // ethernet-ip-s | Rockwell CSP1 | EtherNet IP over TLS | EtherNet IP over DTLS
	Port{2221, 17}:   "rockwell-csp1",        // Rockwell CSP1
	Port{2222, 6}:    "EtherNetIP-1",         // EtherNet IP-1 | EtherNet-IP-1 | EtherNet IP I O
	Port{2222, 17}:   "msantipiracy",         // Microsoft Office OS X antipiracy network monitor
	Port{2223, 6}:    "rockwell-csp2",        // Rockwell CSP2
	Port{2223, 17}:   "rockwell-csp2",        // Rockwell CSP2
	Port{2224, 6}:    "efi-mg",               // Easy Flexible Internet Multiplayer Games
	Port{2224, 17}:   "efi-mg",               // Easy Flexible Internet Multiplayer Games
	Port{2225, 132}:  "rcip-itu",             // Resource Connection Initiation Protocol
	Port{2225, 6}:    "rcip-itu",             // Resource Connection Initiation Protocol
	Port{2225, 17}:   "rcip-itu",             // Resource Connection Initiation Protocol
	Port{2226, 6}:    "di-drm",               // Digital Instinct DRM
	Port{2226, 17}:   "di-drm",               // Digital Instinct DRM
	Port{2227, 6}:    "di-msg",               // DI Messaging Service
	Port{2227, 17}:   "di-msg",               // DI Messaging Service
	Port{2228, 6}:    "ehome-ms",             // eHome Message Server
	Port{2228, 17}:   "ehome-ms",             // eHome Message Server
	Port{2229, 6}:    "datalens",             // DataLens Service
	Port{2229, 17}:   "datalens",             // DataLens Service
	Port{2230, 6}:    "queueadm",             // MetaSoft Job Queue Administration Service
	Port{2230, 17}:   "queueadm",             // MetaSoft Job Queue Administration Service
	Port{2231, 6}:    "wimaxasncp",           // WiMAX ASN Control Plane Protocol
	Port{2231, 17}:   "wimaxasncp",           // WiMAX ASN Control Plane Protocol
	Port{2232, 6}:    "ivs-video",            // IVS Video default
	Port{2232, 17}:   "ivs-video",            // IVS Video default
	Port{2233, 6}:    "infocrypt",            // Missing description for infocrypt
	Port{2233, 17}:   "infocrypt",            // INFOCRYPT
	Port{2234, 6}:    "directplay",           // Missing description for directplay
	Port{2234, 17}:   "directplay",           // DirectPlay
	Port{2235, 6}:    "sercomm-wlink",        // Missing description for sercomm-wlink
	Port{2235, 17}:   "sercomm-wlink",        // Sercomm-WLink
	Port{2236, 6}:    "nani",                 // Missing description for nani
	Port{2236, 17}:   "nani",                 // Nani
	Port{2237, 6}:    "optech-port1-lm",      // Optech Port1 License Manager
	Port{2237, 17}:   "optech-port1-lm",      // Optech Port1 License Manager
	Port{2238, 6}:    "aviva-sna",            // AVIVA SNA SERVER
	Port{2238, 17}:   "aviva-sna",            // AVIVA SNA SERVER
	Port{2239, 6}:    "imagequery",           // Image Query
	Port{2239, 17}:   "imagequery",           // Image Query
	Port{2240, 6}:    "recipe",               // Missing description for recipe
	Port{2240, 17}:   "recipe",               // RECIPe
	Port{2241, 6}:    "ivsd",                 // IVS Daemon
	Port{2241, 17}:   "ivsd",                 // IVS Daemon
	Port{2242, 6}:    "foliocorp",            // Folio Remote Server
	Port{2242, 17}:   "foliocorp",            // Folio Remote Server
	Port{2243, 6}:    "magicom",              // Magicom Protocol
	Port{2243, 17}:   "magicom",              // Magicom Protocol
	Port{2244, 6}:    "nmsserver",            // NMS Server
	Port{2244, 17}:   "nmsserver",            // NMS Server
	Port{2245, 6}:    "hao",                  // Missing description for hao
	Port{2245, 17}:   "hao",                  // HaO
	Port{2246, 6}:    "pc-mta-addrmap",       // PacketCable MTA Addr Map
	Port{2246, 17}:   "pc-mta-addrmap",       // PacketCable MTA Addr Map
	Port{2247, 6}:    "antidotemgrsvr",       // Antidote Deployment Manager Service
	Port{2247, 17}:   "antidotemgrsvr",       // Antidote Deployment Manager Service
	Port{2248, 6}:    "ums",                  // User Management Service
	Port{2248, 17}:   "ums",                  // User Management Service
	Port{2249, 6}:    "rfmp",                 // RISO File Manager Protocol
	Port{2249, 17}:   "rfmp",                 // RISO File Manager Protocol
	Port{2250, 6}:    "remote-collab",        // Missing description for remote-collab
	Port{2250, 17}:   "remote-collab",        // Missing description for remote-collab
	Port{2251, 6}:    "dif-port",             // Distributed Framework Port
	Port{2251, 17}:   "dif-port",             // Distributed Framework Port
	Port{2252, 6}:    "njenet-ssl",           // NJENET using SSL
	Port{2252, 17}:   "njenet-ssl",           // NJENET using SSL
	Port{2253, 6}:    "dtv-chan-req",         // DTV Channel Request
	Port{2253, 17}:   "dtv-chan-req",         // DTV Channel Request
	Port{2254, 6}:    "seispoc",              // Seismic P.O.C. Port
	Port{2254, 17}:   "seispoc",              // Seismic P.O.C. Port
	Port{2255, 6}:    "vrtp",                 // VRTP - ViRtue Transfer Protocol
	Port{2255, 17}:   "vrtp",                 // VRTP - ViRtue Transfer Protocol
	Port{2256, 6}:    "pcc-mfp",              // PCC MFP
	Port{2256, 17}:   "pcc-mfp",              // PCC MFP
	Port{2257, 6}:    "simple-tx-rx",         // simple text file transfer
	Port{2257, 17}:   "simple-tx-rx",         // simple text file transfer
	Port{2258, 6}:    "rcts",                 // Rotorcraft Communications Test System
	Port{2258, 17}:   "rcts",                 // Rotorcraft Communications Test System
	Port{2259, 6}:    "acd-pm",               // Accedian Performance Measurement
	Port{2259, 17}:   "acd-pm",               // Accedian Performance Measurement
	Port{2260, 6}:    "apc-2260",             // APC 2260
	Port{2260, 17}:   "apc-2260",             // APC 2260
	Port{2261, 6}:    "comotionmaster",       // CoMotion Master Server
	Port{2261, 17}:   "comotionmaster",       // CoMotion Master Server
	Port{2262, 6}:    "comotionback",         // CoMotion Backup Server
	Port{2262, 17}:   "comotionback",         // CoMotion Backup Server
	Port{2263, 6}:    "ecwcfg",               // ECweb Configuration Service
	Port{2263, 17}:   "ecwcfg",               // ECweb Configuration Service
	Port{2264, 6}:    "apx500api-1",          // Audio Precision Apx500 API Port 1
	Port{2264, 17}:   "apx500api-1",          // Audio Precision Apx500 API Port 1
	Port{2265, 6}:    "apx500api-2",          // Audio Precision Apx500 API Port 2
	Port{2265, 17}:   "apx500api-2",          // Audio Precision Apx500 API Port 2
	Port{2266, 6}:    "mfserver",             // M-Files Server | M-files Server
	Port{2266, 17}:   "mfserver",             // M-files Server
	Port{2267, 6}:    "ontobroker",           // Missing description for ontobroker
	Port{2267, 17}:   "ontobroker",           // OntoBroker
	Port{2268, 6}:    "amt",                  // Missing description for amt
	Port{2268, 17}:   "amt",                  // AMT
	Port{2269, 6}:    "mikey",                // Missing description for mikey
	Port{2269, 17}:   "mikey",                // MIKEY
	Port{2270, 6}:    "starschool",           // Missing description for starschool
	Port{2270, 17}:   "starschool",           // starSchool
	Port{2271, 6}:    "mmcals",               // Secure Meeting Maker Scheduling
	Port{2271, 17}:   "mmcals",               // Secure Meeting Maker Scheduling
	Port{2272, 6}:    "mmcal",                // Meeting Maker Scheduling
	Port{2272, 17}:   "mmcal",                // Meeting Maker Scheduling
	Port{2273, 6}:    "mysql-im",             // MySQL Instance Manager
	Port{2273, 17}:   "mysql-im",             // MySQL Instance Manager
	Port{2274, 6}:    "pcttunnell",           // PCTTunneller
	Port{2274, 17}:   "pcttunnell",           // PCTTunneller
	Port{2275, 6}:    "ibridge-data",         // iBridge Conferencing
	Port{2275, 17}:   "ibridge-data",         // iBridge Conferencing
	Port{2276, 6}:    "ibridge-mgmt",         // iBridge Management
	Port{2276, 17}:   "ibridge-mgmt",         // iBridge Management
	Port{2277, 6}:    "bluectrlproxy",        // Bt device control proxy
	Port{2277, 17}:   "bluectrlproxy",        // Bt device control proxy
	Port{2278, 6}:    "s3db",                 // Simple Stacked Sequences Database
	Port{2278, 17}:   "s3db",                 // Simple Stacked Sequences Database
	Port{2279, 6}:    "xmquery",              // Missing description for xmquery
	Port{2279, 17}:   "xmquery",              // Missing description for xmquery
	Port{2280, 6}:    "lnvpoller",            // Missing description for lnvpoller
	Port{2280, 17}:   "lnvpoller",            // LNVPOLLER
	Port{2281, 6}:    "lnvconsole",           // Missing description for lnvconsole
	Port{2281, 17}:   "lnvconsole",           // LNVCONSOLE
	Port{2282, 6}:    "lnvalarm",             // Missing description for lnvalarm
	Port{2282, 17}:   "lnvalarm",             // LNVALARM
	Port{2283, 6}:    "lnvstatus",            // Missing description for lnvstatus
	Port{2283, 17}:   "lnvstatus",            // LNVSTATUS
	Port{2284, 6}:    "lnvmaps",              // Missing description for lnvmaps
	Port{2284, 17}:   "lnvmaps",              // LNVMAPS
	Port{2285, 6}:    "lnvmailmon",           // Missing description for lnvmailmon
	Port{2285, 17}:   "lnvmailmon",           // LNVMAILMON
	Port{2286, 6}:    "nas-metering",         // Missing description for nas-metering
	Port{2286, 17}:   "nas-metering",         // NAS-Metering
	Port{2287, 6}:    "dna",                  // Missing description for dna
	Port{2287, 17}:   "dna",                  // DNA
	Port{2288, 6}:    "netml",                // Missing description for netml
	Port{2288, 17}:   "netml",                // NETML
	Port{2289, 6}:    "dict-lookup",          // Lookup dict server
	Port{2289, 17}:   "dict-lookup",          // Lookup dict server
	Port{2290, 6}:    "sonus-logging",        // Sonus Logging Services
	Port{2290, 17}:   "sonus-logging",        // Sonus Logging Services
	Port{2291, 6}:    "eapsp",                // EPSON Advanced Printer Share Protocol
	Port{2291, 17}:   "eapsp",                // EPSON Advanced Printer Share Protocol
	Port{2292, 6}:    "mib-streaming",        // Sonus Element Management Services
	Port{2292, 17}:   "mib-streaming",        // Sonus Element Management Services
	Port{2293, 6}:    "npdbgmngr",            // Network Platform Debug Manager
	Port{2293, 17}:   "npdbgmngr",            // Network Platform Debug Manager
	Port{2294, 6}:    "konshus-lm",           // Konshus License Manager (FLEX)
	Port{2294, 17}:   "konshus-lm",           // Konshus License Manager (FLEX)
	Port{2295, 6}:    "advant-lm",            // Advant License Manager
	Port{2295, 17}:   "advant-lm",            // Advant License Manager
	Port{2296, 6}:    "theta-lm",             // Theta License Manager (Rainbow)
	Port{2296, 17}:   "theta-lm",             // Theta License Manager (Rainbow)
	Port{2297, 6}:    "d2k-datamover1",       // D2K DataMover 1
	Port{2297, 17}:   "d2k-datamover1",       // D2K DataMover 1
	Port{2298, 6}:    "d2k-datamover2",       // D2K DataMover 2
	Port{2298, 17}:   "d2k-datamover2",       // D2K DataMover 2
	Port{2299, 6}:    "pc-telecommute",       // PC Telecommute
	Port{2299, 17}:   "pc-telecommute",       // PC Telecommute
	Port{2300, 6}:    "cvmmon",               // Missing description for cvmmon
	Port{2300, 17}:   "cvmmon",               // CVMMON
	Port{2301, 6}:    "compaqdiag",           // cpq-wbem | Compaq remote diagnostic management | Compaq HTTP
	Port{2301, 17}:   "cpq-wbem",             // Compaq HTTP
	Port{2302, 6}:    "binderysupport",       // Bindery Support
	Port{2302, 17}:   "binderysupport",       // Bindery Support
	Port{2303, 6}:    "proxy-gateway",        // Proxy Gateway
	Port{2303, 17}:   "proxy-gateway",        // Proxy Gateway
	Port{2304, 6}:    "attachmate-uts",       // Attachmate UTS
	Port{2304, 17}:   "attachmate-uts",       // Attachmate UTS
	Port{2305, 6}:    "mt-scaleserver",       // MT ScaleServer
	Port{2305, 17}:   "mt-scaleserver",       // MT ScaleServer
	Port{2306, 6}:    "tappi-boxnet",         // TAPPI BoxNet
	Port{2306, 17}:   "tappi-boxnet",         // TAPPI BoxNet
	Port{2307, 6}:    "pehelp",               // Missing description for pehelp
	Port{2307, 17}:   "pehelp",               // Missing description for pehelp
	Port{2308, 6}:    "sdhelp",               // Missing description for sdhelp
	Port{2308, 17}:   "sdhelp",               // Missing description for sdhelp
	Port{2309, 6}:    "sdserver",             // SD Server
	Port{2309, 17}:   "sdserver",             // SD Server
	Port{2310, 6}:    "sdclient",             // SD Client
	Port{2310, 17}:   "sdclient",             // SD Client
	Port{2311, 6}:    "messageservice",       // Message Service
	Port{2311, 17}:   "messageservice",       // Message Service
	Port{2312, 6}:    "wanscaler",            // WANScaler Communication Service
	Port{2312, 17}:   "wanscaler",            // WANScaler Communication Service
	Port{2313, 6}:    "iapp",                 // IAPP (Inter Access Point Protocol)
	Port{2313, 17}:   "iapp",                 // IAPP (Inter Access Point Protocol)
	Port{2314, 6}:    "cr-websystems",        // CR WebSystems
	Port{2314, 17}:   "cr-websystems",        // CR WebSystems
	Port{2315, 6}:    "precise-sft",          // Precise Sft.
	Port{2315, 17}:   "precise-sft",          // Precise Sft.
	Port{2316, 6}:    "sent-lm",              // SENT License Manager
	Port{2316, 17}:   "sent-lm",              // SENT License Manager
	Port{2317, 6}:    "attachmate-g32",       // Attachmate G32
	Port{2317, 17}:   "attachmate-g32",       // Attachmate G32
	Port{2318, 6}:    "cadencecontrol",       // Cadence Control
	Port{2318, 17}:   "cadencecontrol",       // Cadence Control
	Port{2319, 6}:    "infolibria",           // Missing description for infolibria
	Port{2319, 17}:   "infolibria",           // InfoLibria
	Port{2320, 6}:    "siebel-ns",            // Siebel NS
	Port{2320, 17}:   "siebel-ns",            // Siebel NS
	Port{2321, 6}:    "rdlap",                // Missing description for rdlap
	Port{2321, 17}:   "rdlap",                // RDLAP
	Port{2322, 6}:    "ofsd",                 // Missing description for ofsd
	Port{2322, 17}:   "ofsd",                 // Missing description for ofsd
	Port{2323, 6}:    "3d-nfsd",              // Missing description for 3d-nfsd
	Port{2323, 17}:   "3d-nfsd",              // Missing description for 3d-nfsd
	Port{2324, 6}:    "cosmocall",            // Missing description for cosmocall
	Port{2324, 17}:   "cosmocall",            // Cosmocall
	Port{2325, 6}:    "ansysli",              // ANSYS Licensing Interconnect
	Port{2325, 17}:   "ansysli",              // ANSYS Licensing Interconnect
	Port{2326, 6}:    "idcp",                 // Missing description for idcp
	Port{2326, 17}:   "idcp",                 // IDCP
	Port{2327, 6}:    "xingcsm",              // Missing description for xingcsm
	Port{2327, 17}:   "xingcsm",              // Missing description for xingcsm
	Port{2328, 6}:    "netrix-sftm",          // Netrix SFTM
	Port{2328, 17}:   "netrix-sftm",          // Netrix SFTM
	Port{2329, 6}:    "nvd",                  // Missing description for nvd
	Port{2329, 17}:   "nvd",                  // NVD
	Port{2330, 6}:    "tscchat",              // Missing description for tscchat
	Port{2330, 17}:   "tscchat",              // TSCCHAT
	Port{2331, 6}:    "agentview",            // Missing description for agentview
	Port{2331, 17}:   "agentview",            // AGENTVIEW
	Port{2332, 6}:    "rcc-host",             // RCC Host
	Port{2332, 17}:   "rcc-host",             // RCC Host
	Port{2333, 6}:    "snapp",                // Missing description for snapp
	Port{2333, 17}:   "snapp",                // SNAPP
	Port{2334, 6}:    "ace-client",           // ACE Client Auth
	Port{2334, 17}:   "ace-client",           // ACE Client Auth
	Port{2335, 6}:    "ace-proxy",            // ACE Proxy
	Port{2335, 17}:   "ace-proxy",            // ACE Proxy
	Port{2336, 6}:    "appleugcontrol",       // Apple UG Control
	Port{2336, 17}:   "appleugcontrol",       // Apple UG Control
	Port{2337, 6}:    "ideesrv",              // Missing description for ideesrv
	Port{2337, 17}:   "ideesrv",              // Missing description for ideesrv
	Port{2338, 6}:    "norton-lambert",       // Norton Lambert
	Port{2338, 17}:   "norton-lambert",       // Norton Lambert
	Port{2339, 6}:    "3com-webview",         // 3Com WebView
	Port{2339, 17}:   "3com-webview",         // 3Com WebView
	Port{2340, 6}:    "wrs_registry",         // wrs-registry | WRS Registry
	Port{2340, 17}:   "wrs_registry",         // WRS Registry
	Port{2341, 6}:    "xiostatus",            // XIO Status
	Port{2341, 17}:   "xiostatus",            // XIO Status
	Port{2342, 6}:    "manage-exec",          // Seagate Manage Exec
	Port{2342, 17}:   "manage-exec",          // Seagate Manage Exec
	Port{2343, 6}:    "nati-logos",           // nati logos
	Port{2343, 17}:   "nati-logos",           // nati logos
	Port{2344, 6}:    "fcmsys",               // Missing description for fcmsys
	Port{2344, 17}:   "fcmsys",               // Missing description for fcmsys
	Port{2345, 6}:    "dbm",                  // Missing description for dbm
	Port{2345, 17}:   "dbm",                  // Missing description for dbm
	Port{2346, 6}:    "redstorm_join",        // redstorm-join | Game Connection Port
	Port{2346, 17}:   "redstorm_join",        // Game Connection Port
	Port{2347, 6}:    "redstorm_find",        // redstorm-find | Game Announcement and Location
	Port{2347, 17}:   "redstorm_find",        // Game Announcement and Location
	Port{2348, 6}:    "redstorm_info",        // redstorm-info | Information to query for game status
	Port{2348, 17}:   "redstorm_info",        // Information to query for game status
	Port{2349, 6}:    "redstorm_diag",        // redstorm-diag | Diagnostics Port
	Port{2349, 17}:   "redstorm_diag",        // Diagnostics Port
	Port{2350, 6}:    "psbserver",            // Pharos Booking Server
	Port{2350, 17}:   "psbserver",            // Pharos Booking Server
	Port{2351, 6}:    "psrserver",            // Missing description for psrserver
	Port{2351, 17}:   "psrserver",            // Missing description for psrserver
	Port{2352, 6}:    "pslserver",            // Missing description for pslserver
	Port{2352, 17}:   "pslserver",            // Missing description for pslserver
	Port{2353, 6}:    "pspserver",            // Missing description for pspserver
	Port{2353, 17}:   "pspserver",            // Missing description for pspserver
	Port{2354, 6}:    "psprserver",           // Missing description for psprserver
	Port{2354, 17}:   "psprserver",           // Missing description for psprserver
	Port{2355, 6}:    "psdbserver",           // Missing description for psdbserver
	Port{2355, 17}:   "psdbserver",           // Missing description for psdbserver
	Port{2356, 6}:    "gxtelmd",              // GXT License Managemant
	Port{2356, 17}:   "gxtelmd",              // GXT License Managemant
	Port{2357, 6}:    "unihub-server",        // UniHub Server
	Port{2357, 17}:   "unihub-server",        // UniHub Server
	Port{2358, 6}:    "futrix",               // Missing description for futrix
	Port{2358, 17}:   "futrix",               // Futrix
	Port{2359, 6}:    "flukeserver",          // Missing description for flukeserver
	Port{2359, 17}:   "flukeserver",          // FlukeServer
	Port{2360, 6}:    "nexstorindltd",        // Missing description for nexstorindltd
	Port{2360, 17}:   "nexstorindltd",        // NexstorIndLtd
	Port{2361, 6}:    "tl1",                  // Missing description for tl1
	Port{2361, 17}:   "tl1",                  // TL1
	Port{2362, 6}:    "digiman",              // Missing description for digiman
	Port{2362, 17}:   "digiman",              // Missing description for digiman
	Port{2363, 6}:    "mediacntrlnfsd",       // Media Central NFSD
	Port{2363, 17}:   "mediacntrlnfsd",       // Media Central NFSD
	Port{2364, 6}:    "oi-2000",              // Missing description for oi-2000
	Port{2364, 17}:   "oi-2000",              // OI-2000
	Port{2365, 6}:    "dbref",                // Missing description for dbref
	Port{2365, 17}:   "dbref",                // Missing description for dbref
	Port{2366, 6}:    "qip-login",            // Missing description for qip-login
	Port{2366, 17}:   "qip-login",            // Missing description for qip-login
	Port{2367, 6}:    "service-ctrl",         // Service Control
	Port{2367, 17}:   "service-ctrl",         // Service Control
	Port{2368, 6}:    "opentable",            // Missing description for opentable
	Port{2368, 17}:   "opentable",            // OpenTable
	Port{2370, 6}:    "l3-hbmon",             // Missing description for l3-hbmon
	Port{2370, 17}:   "l3-hbmon",             // L3-HBMon
	Port{2371, 6}:    "worldwire",            // hp-rda | Compaq WorldWire Port | HP Remote Device Access
	Port{2371, 17}:   "worldwire",            // Compaq WorldWire Port
	Port{2372, 6}:    "lanmessenger",         // Missing description for lanmessenger
	Port{2372, 17}:   "lanmessenger",         // LanMessenger
	Port{2373, 6}:    "remographlm",          // Remograph License Manager
	Port{2374, 6}:    "hydra",                // Hydra RPC
	Port{2375, 6}:    "docker",               // docker.com | Docker REST API (plain text)
	Port{2376, 6}:    "docker",               // docker-s | docker.com | Docker REST API (ssl)
	Port{2377, 6}:    "swarm",                // RPC interface for Docker Swarm
	Port{2379, 6}:    "etcd-client",          // etcd client communication
	Port{2380, 6}:    "etcd-server",          // etcd server to server communication
	Port{2381, 6}:    "compaq-https",         // Compaq HTTPS
	Port{2381, 17}:   "compaq-https",         // Compaq HTTPS
	Port{2382, 6}:    "ms-olap3",             // Microsoft OLAP
	Port{2382, 17}:   "ms-olap3",             // Microsoft OLAP
	Port{2383, 6}:    "ms-olap4",             // MS OLAP 4 | Microsoft OLAP
	Port{2383, 17}:   "ms-olap4",             // Microsoft OLAP
	Port{2384, 6}:    "sd-request",           // sd-capacity | SD-CAPACITY
	Port{2384, 17}:   "sd-capacity",          // SD-CAPACITY
	Port{2385, 6}:    "sd-data",              // Missing description for sd-data
	Port{2385, 17}:   "sd-data",              // SD-DATA
	Port{2386, 6}:    "virtualtape",          // Virtual Tape
	Port{2386, 17}:   "virtualtape",          // Virtual Tape
	Port{2387, 6}:    "vsamredirector",       // VSAM Redirector
	Port{2387, 17}:   "vsamredirector",       // VSAM Redirector
	Port{2388, 6}:    "mynahautostart",       // MYNAH AutoStart
	Port{2388, 17}:   "mynahautostart",       // MYNAH AutoStart
	Port{2389, 6}:    "ovsessionmgr",         // OpenView Session Mgr
	Port{2389, 17}:   "ovsessionmgr",         // OpenView Session Mgr
	Port{2390, 6}:    "rsmtp",                // Missing description for rsmtp
	Port{2390, 17}:   "rsmtp",                // RSMTP
	Port{2391, 6}:    "3com-net-mgmt",        // 3COM Net Management
	Port{2391, 17}:   "3com-net-mgmt",        // 3COM Net Management
	Port{2392, 6}:    "tacticalauth",         // Tactical Auth
	Port{2392, 17}:   "tacticalauth",         // Tactical Auth
	Port{2393, 6}:    "ms-olap1",             // SQL Server Downlevel OLAP Client Support | MS OLAP 1
	Port{2393, 17}:   "ms-olap1",             // MS OLAP 1
	Port{2394, 6}:    "ms-olap2",             // SQL Server Downlevel OLAP Client Support | MS OLAP 2
	Port{2394, 17}:   "ms-olap2",             // MS OLAP 2
	Port{2395, 6}:    "lan900_remote",        // lan900-remote | LAN900 Remote
	Port{2395, 17}:   "lan900_remote",        // LAN900 Remote
	Port{2396, 6}:    "wusage",               // Missing description for wusage
	Port{2396, 17}:   "wusage",               // Wusage
	Port{2397, 6}:    "ncl",                  // Missing description for ncl
	Port{2397, 17}:   "ncl",                  // NCL
	Port{2398, 6}:    "orbiter",              // Missing description for orbiter
	Port{2398, 17}:   "orbiter",              // Orbiter
	Port{2399, 6}:    "fmpro-fdal",           // FileMaker, Inc. - Data Access Layer
	Port{2399, 17}:   "fmpro-fdal",           // FileMaker, Inc. - Data Access Layer
	Port{2400, 6}:    "opequus-server",       // OpEquus Server
	Port{2400, 17}:   "opequus-server",       // OpEquus Server
	Port{2401, 6}:    "cvspserver",           // CVS network server
	Port{2401, 17}:   "cvspserver",           // CVS network server
	Port{2402, 6}:    "taskmaster2000",       // TaskMaster 2000 Server
	Port{2402, 17}:   "taskmaster2000",       // TaskMaster 2000 Server
	Port{2403, 6}:    "taskmaster2000",       // TaskMaster 2000 Web
	Port{2403, 17}:   "taskmaster2000",       // TaskMaster 2000 Web
	Port{2404, 6}:    "iec-104",              // IEC 60870-5-104 process control over IP
	Port{2404, 17}:   "iec-104",              // IEC 60870-5-104 process control over IP
	Port{2405, 6}:    "trc-netpoll",          // TRC Netpoll
	Port{2405, 17}:   "trc-netpoll",          // TRC Netpoll
	Port{2406, 6}:    "jediserver",           // Missing description for jediserver
	Port{2406, 17}:   "jediserver",           // JediServer
	Port{2407, 6}:    "orion",                // Missing description for orion
	Port{2407, 17}:   "orion",                // Orion
	Port{2408, 6}:    "optimanet",            // railgun-webaccl | CloudFlare Railgun Web Acceleration Protocol
	Port{2408, 17}:   "optimanet",            // OptimaNet
	Port{2409, 6}:    "sns-protocol",         // SNS Protocol
	Port{2409, 17}:   "sns-protocol",         // SNS Protocol
	Port{2410, 6}:    "vrts-registry",        // VRTS Registry
	Port{2410, 17}:   "vrts-registry",        // VRTS Registry
	Port{2411, 6}:    "netwave-ap-mgmt",      // Netwave AP Management
	Port{2411, 17}:   "netwave-ap-mgmt",      // Netwave AP Management
	Port{2412, 6}:    "cdn",                  // Missing description for cdn
	Port{2412, 17}:   "cdn",                  // CDN
	Port{2413, 6}:    "orion-rmi-reg",        // Missing description for orion-rmi-reg
	Port{2413, 17}:   "orion-rmi-reg",        // Missing description for orion-rmi-reg
	Port{2414, 6}:    "beeyond",              // Missing description for beeyond
	Port{2414, 17}:   "beeyond",              // Beeyond
	Port{2415, 6}:    "codima-rtp",           // Codima Remote Transaction Protocol
	Port{2415, 17}:   "codima-rtp",           // Codima Remote Transaction Protocol
	Port{2416, 6}:    "rmtserver",            // RMT Server
	Port{2416, 17}:   "rmtserver",            // RMT Server
	Port{2417, 6}:    "composit-server",      // Composit Server
	Port{2417, 17}:   "composit-server",      // Composit Server
	Port{2418, 6}:    "cas",                  // Missing description for cas
	Port{2418, 17}:   "cas",                  // Missing description for cas
	Port{2419, 6}:    "attachmate-s2s",       // Attachmate S2S
	Port{2419, 17}:   "attachmate-s2s",       // Attachmate S2S
	Port{2420, 6}:    "dslremote-mgmt",       // DSL Remote Management
	Port{2420, 17}:   "dslremote-mgmt",       // DSL Remote Management
	Port{2421, 6}:    "g-talk",               // Missing description for g-talk
	Port{2421, 17}:   "g-talk",               // G-Talk
	Port{2422, 6}:    "crmsbits",             // Missing description for crmsbits
	Port{2422, 17}:   "crmsbits",             // CRMSBITS
	Port{2423, 6}:    "rnrp",                 // Missing description for rnrp
	Port{2423, 17}:   "rnrp",                 // RNRP
	Port{2424, 6}:    "kofax-svr",            // Missing description for kofax-svr
	Port{2424, 17}:   "kofax-svr",            // KOFAX-SVR
	Port{2425, 6}:    "fjitsuappmgr",         // Fujitsu App Manager
	Port{2425, 17}:   "fjitsuappmgr",         // Fujitsu App Manager
	Port{2426, 6}:    "vcmp",                 // VeloCloud MultiPath Protocol
	Port{2427, 132}:  "mgcp-gateway",         // Media Gateway Control Protocol Gateway
	Port{2427, 6}:    "mgcp-gateway",         // Media Gateway Control Protocol Gateway
	Port{2427, 17}:   "mgcp-gateway",         // Media Gateway Control Protocol Gateway
	Port{2428, 6}:    "ott",                  // One Way Trip Time
	Port{2428, 17}:   "ott",                  // One Way Trip Time
	Port{2429, 6}:    "ft-role",              // Missing description for ft-role
	Port{2429, 17}:   "ft-role",              // FT-ROLE
	Port{2430, 6}:    "venus",                // Missing description for venus
	Port{2430, 17}:   "venus",                // Missing description for venus
	Port{2431, 6}:    "venus-se",             // Missing description for venus-se
	Port{2431, 17}:   "venus-se",             // Missing description for venus-se
	Port{2432, 6}:    "codasrv",              // Missing description for codasrv
	Port{2432, 17}:   "codasrv",              // Missing description for codasrv
	Port{2433, 6}:    "codasrv-se",           // Missing description for codasrv-se
	Port{2433, 17}:   "codasrv-se",           // Missing description for codasrv-se
	Port{2434, 6}:    "pxc-epmap",            // Missing description for pxc-epmap
	Port{2434, 17}:   "pxc-epmap",            // Missing description for pxc-epmap
	Port{2435, 6}:    "optilogic",            // Missing description for optilogic
	Port{2435, 17}:   "optilogic",            // OptiLogic
	Port{2436, 6}:    "topx",                 // TOP X
	Port{2436, 17}:   "topx",                 // TOP X
	Port{2437, 6}:    "unicontrol",           // Missing description for unicontrol
	Port{2437, 17}:   "unicontrol",           // UniControl
	Port{2438, 6}:    "msp",                  // Missing description for msp
	Port{2438, 17}:   "msp",                  // MSP
	Port{2439, 6}:    "sybasedbsynch",        // Missing description for sybasedbsynch
	Port{2439, 17}:   "sybasedbsynch",        // SybaseDBSynch
	Port{2440, 6}:    "spearway",             // Spearway Lockers
	Port{2440, 17}:   "spearway",             // Spearway Lockers
	Port{2441, 6}:    "pvsw-inet",            // Pervasive I*net Data Server
	Port{2441, 17}:   "pvsw-inet",            // Pervasive I*net Data Server
	Port{2442, 6}:    "netangel",             // Missing description for netangel
	Port{2442, 17}:   "netangel",             // Netangel
	Port{2443, 6}:    "powerclientcsf",       // PowerClient Central Storage Facility
	Port{2443, 17}:   "powerclientcsf",       // PowerClient Central Storage Facility
	Port{2444, 6}:    "btpp2sectrans",        // BT PP2 Sectrans
	Port{2444, 17}:   "btpp2sectrans",        // BT PP2 Sectrans
	Port{2445, 6}:    "dtn1",                 // Missing description for dtn1
	Port{2445, 17}:   "dtn1",                 // DTN1
	Port{2446, 6}:    "bues_service",         // bues-service
	Port{2446, 17}:   "bues_service",         // Missing description for bues_service
	Port{2447, 6}:    "ovwdb",                // OpenView NNM daemon
	Port{2447, 17}:   "ovwdb",                // OpenView NNM daemon
	Port{2448, 6}:    "hpppssvr",             // hpppsvr
	Port{2448, 17}:   "hpppssvr",             // hpppsvr
	Port{2449, 6}:    "ratl",                 // Missing description for ratl
	Port{2449, 17}:   "ratl",                 // RATL
	Port{2450, 6}:    "netadmin",             // Missing description for netadmin
	Port{2450, 17}:   "netadmin",             // Missing description for netadmin
	Port{2451, 6}:    "netchat",              // Missing description for netchat
	Port{2451, 17}:   "netchat",              // Missing description for netchat
	Port{2452, 6}:    "snifferclient",        // Missing description for snifferclient
	Port{2452, 17}:   "snifferclient",        // SnifferClient
	Port{2453, 6}:    "madge-ltd",            // madge ltd
	Port{2453, 17}:   "madge-ltd",            // madge ltd
	Port{2454, 6}:    "indx-dds",             // Missing description for indx-dds
	Port{2454, 17}:   "indx-dds",             // IndX-DDS
	Port{2455, 6}:    "wago-io-system",       // Missing description for wago-io-system
	Port{2455, 17}:   "wago-io-system",       // WAGO-IO-SYSTEM
	Port{2456, 6}:    "altav-remmgt",         // Missing description for altav-remmgt
	Port{2456, 17}:   "altav-remmgt",         // Missing description for altav-remmgt
	Port{2457, 6}:    "rapido-ip",            // Rapido_IP
	Port{2457, 17}:   "rapido-ip",            // Rapido_IP
	Port{2458, 6}:    "griffin",              // Missing description for griffin
	Port{2458, 17}:   "griffin",              // Missing description for griffin
	Port{2459, 6}:    "community",            // Missing description for community
	Port{2459, 17}:   "community",            // Community
	Port{2460, 6}:    "ms-theater",           // Missing description for ms-theater
	Port{2460, 17}:   "ms-theater",           // Missing description for ms-theater
	Port{2461, 6}:    "qadmifoper",           // Missing description for qadmifoper
	Port{2461, 17}:   "qadmifoper",           // Missing description for qadmifoper
	Port{2462, 6}:    "qadmifevent",          // Missing description for qadmifevent
	Port{2462, 17}:   "qadmifevent",          // Missing description for qadmifevent
	Port{2463, 6}:    "lsi-raid-mgmt",        // LSI RAID Management
	Port{2463, 17}:   "lsi-raid-mgmt",        // LSI RAID Management
	Port{2464, 6}:    "direcpc-si",           // DirecPC SI
	Port{2464, 17}:   "direcpc-si",           // DirecPC SI
	Port{2465, 6}:    "lbm",                  // Load Balance Management
	Port{2465, 17}:   "lbm",                  // Load Balance Management
	Port{2466, 6}:    "lbf",                  // Load Balance Forwarding
	Port{2466, 17}:   "lbf",                  // Load Balance Forwarding
	Port{2467, 6}:    "high-criteria",        // High Criteria
	Port{2467, 17}:   "high-criteria",        // High Criteria
	Port{2468, 6}:    "qip-msgd",             // qip_msgd
	Port{2468, 17}:   "qip-msgd",             // qip_msgd
	Port{2469, 6}:    "mti-tcs-comm",         // Missing description for mti-tcs-comm
	Port{2469, 17}:   "mti-tcs-comm",         // MTI-TCS-COMM
	Port{2470, 6}:    "taskman-port",         // taskman port
	Port{2470, 17}:   "taskman-port",         // taskman port
	Port{2471, 6}:    "seaodbc",              // Missing description for seaodbc
	Port{2471, 17}:   "seaodbc",              // SeaODBC
	Port{2472, 6}:    "c3",                   // Missing description for c3
	Port{2472, 17}:   "c3",                   // C3
	Port{2473, 6}:    "aker-cdp",             // Missing description for aker-cdp
	Port{2473, 17}:   "aker-cdp",             // Aker-cdp
	Port{2474, 6}:    "vitalanalysis",        // Vital Analysis
	Port{2474, 17}:   "vitalanalysis",        // Vital Analysis
	Port{2475, 6}:    "ace-server",           // ACE Server
	Port{2475, 17}:   "ace-server",           // ACE Server
	Port{2476, 6}:    "ace-svr-prop",         // ACE Server Propagation
	Port{2476, 17}:   "ace-svr-prop",         // ACE Server Propagation
	Port{2477, 6}:    "ssm-cvs",              // SecurSight Certificate Valifation Service
	Port{2477, 17}:   "ssm-cvs",              // SecurSight Certificate Valifation Service
	Port{2478, 6}:    "ssm-cssps",            // SecurSight Authentication Server (SSL)
	Port{2478, 17}:   "ssm-cssps",            // SecurSight Authentication Server (SSL)
	Port{2479, 6}:    "ssm-els",              // SecurSight Event Logging Server (SSL)
	Port{2479, 17}:   "ssm-els",              // SecurSight Event Logging Server (SSL)
	Port{2480, 6}:    "powerexchange",        // Informatica PowerExchange Listener
	Port{2480, 17}:   "powerexchange",        // Informatica PowerExchange Listener
	Port{2481, 6}:    "giop",                 // Oracle GIOP
	Port{2481, 17}:   "giop",                 // Oracle GIOP
	Port{2482, 6}:    "giop-ssl",             // Oracle GIOP SSL
	Port{2482, 17}:   "giop-ssl",             // Oracle GIOP SSL
	Port{2483, 6}:    "ttc",                  // Oracle TTC
	Port{2483, 17}:   "ttc",                  // Oracle TTC
	Port{2484, 6}:    "ttc-ssl",              // Oracle TTC SSL
	Port{2484, 17}:   "ttc-ssl",              // Oracle TTC SSL
	Port{2485, 6}:    "netobjects1",          // Net Objects1
	Port{2485, 17}:   "netobjects1",          // Net Objects1
	Port{2486, 6}:    "netobjects2",          // Net Objects2
	Port{2486, 17}:   "netobjects2",          // Net Objects2
	Port{2487, 6}:    "pns",                  // Policy Notice Service
	Port{2487, 17}:   "pns",                  // Policy Notice Service
	Port{2488, 6}:    "moy-corp",             // Moy Corporation
	Port{2488, 17}:   "moy-corp",             // Moy Corporation
	Port{2489, 6}:    "tsilb",                // Missing description for tsilb
	Port{2489, 17}:   "tsilb",                // TSILB
	Port{2490, 6}:    "qip-qdhcp",            // qip_qdhcp
	Port{2490, 17}:   "qip-qdhcp",            // qip_qdhcp
	Port{2491, 6}:    "conclave-cpp",         // Conclave CPP
	Port{2491, 17}:   "conclave-cpp",         // Conclave CPP
	Port{2492, 6}:    "groove",               // Missing description for groove
	Port{2492, 17}:   "groove",               // GROOVE
	Port{2493, 6}:    "talarian-mqs",         // Talarian MQS
	Port{2493, 17}:   "talarian-mqs",         // Talarian MQS
	Port{2494, 6}:    "bmc-ar",               // BMC AR
	Port{2494, 17}:   "bmc-ar",               // BMC AR
	Port{2495, 6}:    "fast-rem-serv",        // Fast Remote Services
	Port{2495, 17}:   "fast-rem-serv",        // Fast Remote Services
	Port{2496, 6}:    "dirgis",               // Missing description for dirgis
	Port{2496, 17}:   "dirgis",               // DIRGIS
	Port{2497, 6}:    "quaddb",               // Quad DB
	Port{2497, 17}:   "quaddb",               // Quad DB
	Port{2498, 6}:    "odn-castraq",          // Missing description for odn-castraq
	Port{2498, 17}:   "odn-castraq",          // ODN-CasTraq
	Port{2499, 6}:    "unicontrol",           // Missing description for unicontrol
	Port{2499, 17}:   "unicontrol",           // UniControl
	Port{2500, 6}:    "rtsserv",              // Resource Tracking system server
	Port{2500, 17}:   "rtsserv",              // Resource Tracking system server
	Port{2501, 6}:    "rtsclient",            // Resource Tracking system client
	Port{2501, 17}:   "rtsclient",            // Resource Tracking system client
	Port{2502, 6}:    "kentrox-prot",         // Kentrox Protocol
	Port{2502, 17}:   "kentrox-prot",         // Kentrox Protocol
	Port{2503, 6}:    "nms-dpnss",            // Missing description for nms-dpnss
	Port{2503, 17}:   "nms-dpnss",            // NMS-DPNSS
	Port{2504, 6}:    "wlbs",                 // Missing description for wlbs
	Port{2504, 17}:   "wlbs",                 // WLBS
	Port{2505, 6}:    "ppcontrol",            // PowerPlay Control
	Port{2505, 17}:   "ppcontrol",            // PowerPlay Control
	Port{2506, 6}:    "jbroker",              // Missing description for jbroker
	Port{2506, 17}:   "jbroker",              // Missing description for jbroker
	Port{2507, 6}:    "spock",                // Missing description for spock
	Port{2507, 17}:   "spock",                // Missing description for spock
	Port{2508, 6}:    "jdatastore",           // Missing description for jdatastore
	Port{2508, 17}:   "jdatastore",           // JDataStore
	Port{2509, 6}:    "fjmpss",               // Missing description for fjmpss
	Port{2509, 17}:   "fjmpss",               // Missing description for fjmpss
	Port{2510, 6}:    "fjappmgrbulk",         // Missing description for fjappmgrbulk
	Port{2510, 17}:   "fjappmgrbulk",         // Missing description for fjappmgrbulk
	Port{2511, 6}:    "metastorm",            // Missing description for metastorm
	Port{2511, 17}:   "metastorm",            // Metastorm
	Port{2512, 6}:    "citrixima",            // Citrix IMA
	Port{2512, 17}:   "citrixima",            // Citrix IMA
	Port{2513, 6}:    "citrixadmin",          // Citrix ADMIN
	Port{2513, 17}:   "citrixadmin",          // Citrix ADMIN
	Port{2514, 6}:    "facsys-ntp",           // Facsys NTP
	Port{2514, 17}:   "facsys-ntp",           // Facsys NTP
	Port{2515, 6}:    "facsys-router",        // Facsys Router
	Port{2515, 17}:   "facsys-router",        // Facsys Router
	Port{2516, 6}:    "maincontrol",          // Main Control
	Port{2516, 17}:   "maincontrol",          // Main Control
	Port{2517, 6}:    "call-sig-trans",       // H.323 Annex E call signaling transport | H.323 Annex E Call Control Signalling Transport
	Port{2517, 17}:   "call-sig-trans",       // H.323 Annex E call signaling transport
	Port{2518, 6}:    "willy",                // Missing description for willy
	Port{2518, 17}:   "willy",                // Willy
	Port{2519, 6}:    "globmsgsvc",           // Missing description for globmsgsvc
	Port{2519, 17}:   "globmsgsvc",           // Missing description for globmsgsvc
	Port{2520, 6}:    "pvsw",                 // Pervasive Listener
	Port{2520, 17}:   "pvsw",                 // Pervasive Listener
	Port{2521, 6}:    "adaptecmgr",           // Adaptec Manager
	Port{2521, 17}:   "adaptecmgr",           // Adaptec Manager
	Port{2522, 6}:    "windb",                // Missing description for windb
	Port{2522, 17}:   "windb",                // WinDb
	Port{2523, 6}:    "qke-llc-v3",           // Qke LLC V.3
	Port{2523, 17}:   "qke-llc-v3",           // Qke LLC V.3
	Port{2524, 6}:    "optiwave-lm",          // Optiwave License Management
	Port{2524, 17}:   "optiwave-lm",          // Optiwave License Management
	Port{2525, 6}:    "ms-v-worlds",          // MS V-Worlds
	Port{2525, 17}:   "ms-v-worlds",          // MS V-Worlds
	Port{2526, 6}:    "ema-sent-lm",          // EMA License Manager
	Port{2526, 17}:   "ema-sent-lm",          // EMA License Manager
	Port{2527, 6}:    "iqserver",             // IQ Server
	Port{2527, 17}:   "iqserver",             // IQ Server
	Port{2528, 6}:    "ncr_ccl",              // ncr-ccl | NCR CCL
	Port{2528, 17}:   "ncr_ccl",              // NCR CCL
	Port{2529, 6}:    "utsftp",               // UTS FTP
	Port{2529, 17}:   "utsftp",               // UTS FTP
	Port{2530, 6}:    "vrcommerce",           // VR Commerce
	Port{2530, 17}:   "vrcommerce",           // VR Commerce
	Port{2531, 6}:    "ito-e-gui",            // ITO-E GUI
	Port{2531, 17}:   "ito-e-gui",            // ITO-E GUI
	Port{2532, 6}:    "ovtopmd",              // Missing description for ovtopmd
	Port{2532, 17}:   "ovtopmd",              // OVTOPMD
	Port{2533, 6}:    "snifferserver",        // Missing description for snifferserver
	Port{2533, 17}:   "snifferserver",        // SnifferServer
	Port{2534, 6}:    "combox-web-acc",       // Combox Web Access
	Port{2534, 17}:   "combox-web-acc",       // Combox Web Access
	Port{2535, 6}:    "madcap",               // Missing description for madcap
	Port{2535, 17}:   "madcap",               // MADCAP
	Port{2536, 6}:    "btpp2audctr1",         // Missing description for btpp2audctr1
	Port{2536, 17}:   "btpp2audctr1",         // Missing description for btpp2audctr1
	Port{2537, 6}:    "upgrade",              // Upgrade Protocol
	Port{2537, 17}:   "upgrade",              // Upgrade Protocol
	Port{2538, 6}:    "vnwk-prapi",           // Missing description for vnwk-prapi
	Port{2538, 17}:   "vnwk-prapi",           // Missing description for vnwk-prapi
	Port{2539, 6}:    "vsiadmin",             // VSI Admin
	Port{2539, 17}:   "vsiadmin",             // VSI Admin
	Port{2540, 6}:    "lonworks",             // Missing description for lonworks
	Port{2540, 17}:   "lonworks",             // LonWorks
	Port{2541, 6}:    "lonworks2",            // Missing description for lonworks2
	Port{2541, 17}:   "lonworks2",            // LonWorks2
	Port{2542, 6}:    "udrawgraph",           // uDraw(Graph)
	Port{2542, 17}:   "udrawgraph",           // uDraw(Graph)
	Port{2543, 6}:    "reftek",               // Missing description for reftek
	Port{2543, 17}:   "reftek",               // REFTEK
	Port{2544, 6}:    "novell-zen",           // Management Daemon Refresh
	Port{2544, 17}:   "novell-zen",           // Management Daemon Refresh
	Port{2545, 6}:    "sis-emt",              // Missing description for sis-emt
	Port{2545, 17}:   "sis-emt",              // Missing description for sis-emt
	Port{2546, 6}:    "vytalvaultbrtp",       // Missing description for vytalvaultbrtp
	Port{2546, 17}:   "vytalvaultbrtp",       // Missing description for vytalvaultbrtp
	Port{2547, 6}:    "vytalvaultvsmp",       // Missing description for vytalvaultvsmp
	Port{2547, 17}:   "vytalvaultvsmp",       // Missing description for vytalvaultvsmp
	Port{2548, 6}:    "vytalvaultpipe",       // Missing description for vytalvaultpipe
	Port{2548, 17}:   "vytalvaultpipe",       // Missing description for vytalvaultpipe
	Port{2549, 6}:    "ipass",                // Missing description for ipass
	Port{2549, 17}:   "ipass",                // IPASS
	Port{2550, 6}:    "ads",                  // Missing description for ads
	Port{2550, 17}:   "ads",                  // ADS
	Port{2551, 6}:    "isg-uda-server",       // ISG UDA Server
	Port{2551, 17}:   "isg-uda-server",       // ISG UDA Server
	Port{2552, 6}:    "call-logging",         // Call Logging
	Port{2552, 17}:   "call-logging",         // Call Logging
	Port{2553, 6}:    "efidiningport",        // Missing description for efidiningport
	Port{2553, 17}:   "efidiningport",        // Missing description for efidiningport
	Port{2554, 6}:    "vcnet-link-v10",       // VCnet-Link v10
	Port{2554, 17}:   "vcnet-link-v10",       // VCnet-Link v10
	Port{2555, 6}:    "compaq-wcp",           // Compaq WCP
	Port{2555, 17}:   "compaq-wcp",           // Compaq WCP
	Port{2556, 6}:    "nicetec-nmsvc",        // Missing description for nicetec-nmsvc
	Port{2556, 17}:   "nicetec-nmsvc",        // Missing description for nicetec-nmsvc
	Port{2557, 6}:    "nicetec-mgmt",         // Missing description for nicetec-mgmt
	Port{2557, 17}:   "nicetec-mgmt",         // Missing description for nicetec-mgmt
	Port{2558, 6}:    "pclemultimedia",       // PCLE Multi Media
	Port{2558, 17}:   "pclemultimedia",       // PCLE Multi Media
	Port{2559, 6}:    "lstp",                 // Missing description for lstp
	Port{2559, 17}:   "lstp",                 // LSTP
	Port{2560, 6}:    "labrat",               // Missing description for labrat
	Port{2560, 17}:   "labrat",               // Missing description for labrat
	Port{2561, 6}:    "mosaixcc",             // Missing description for mosaixcc
	Port{2561, 17}:   "mosaixcc",             // MosaixCC
	Port{2562, 6}:    "delibo",               // Missing description for delibo
	Port{2562, 17}:   "delibo",               // Delibo
	Port{2563, 6}:    "cti-redwood",          // CTI Redwood
	Port{2563, 17}:   "cti-redwood",          // CTI Redwood
	Port{2564, 6}:    "hp-3000-telnet",       // HP 3000 NS VT block mode telnet
	Port{2565, 6}:    "coord-svr",            // Coordinator Server
	Port{2565, 17}:   "coord-svr",            // Coordinator Server
	Port{2566, 6}:    "pcs-pcw",              // Missing description for pcs-pcw
	Port{2566, 17}:   "pcs-pcw",              // Missing description for pcs-pcw
	Port{2567, 6}:    "clp",                  // Cisco Line Protocol
	Port{2567, 17}:   "clp",                  // Cisco Line Protocol
	Port{2568, 6}:    "spamtrap",             // SPAM TRAP
	Port{2568, 17}:   "spamtrap",             // SPAM TRAP
	Port{2569, 6}:    "sonuscallsig",         // Sonus Call Signal
	Port{2569, 17}:   "sonuscallsig",         // Sonus Call Signal
	Port{2570, 6}:    "hs-port",              // HS Port
	Port{2570, 17}:   "hs-port",              // HS Port
	Port{2571, 6}:    "cecsvc",               // Missing description for cecsvc
	Port{2571, 17}:   "cecsvc",               // CECSVC
	Port{2572, 6}:    "ibp",                  // Missing description for ibp
	Port{2572, 17}:   "ibp",                  // IBP
	Port{2573, 6}:    "trustestablish",       // Trust Establish
	Port{2573, 17}:   "trustestablish",       // Trust Establish
	Port{2574, 6}:    "blockade-bpsp",        // Blockade BPSP
	Port{2574, 17}:   "blockade-bpsp",        // Blockade BPSP
	Port{2575, 6}:    "hl7",                  // Missing description for hl7
	Port{2575, 17}:   "hl7",                  // HL7
	Port{2576, 6}:    "tclprodebugger",       // TCL Pro Debugger
	Port{2576, 17}:   "tclprodebugger",       // TCL Pro Debugger
	Port{2577, 6}:    "scipticslsrvr",        // Scriptics Lsrvr
	Port{2577, 17}:   "scipticslsrvr",        // Scriptics Lsrvr
	Port{2578, 6}:    "rvs-isdn-dcp",         // RVS ISDN DCP
	Port{2578, 17}:   "rvs-isdn-dcp",         // RVS ISDN DCP
	Port{2579, 6}:    "mpfoncl",              // Missing description for mpfoncl
	Port{2579, 17}:   "mpfoncl",              // Missing description for mpfoncl
	Port{2580, 6}:    "tributary",            // Missing description for tributary
	Port{2580, 17}:   "tributary",            // Tributary
	Port{2581, 6}:    "argis-te",             // ARGIS TE
	Port{2581, 17}:   "argis-te",             // ARGIS TE
	Port{2582, 6}:    "argis-ds",             // ARGIS DS
	Port{2582, 17}:   "argis-ds",             // ARGIS DS
	Port{2583, 6}:    "mon",                  // Missing description for mon
	Port{2583, 17}:   "mon",                  // MON
	Port{2584, 6}:    "cyaserv",              // Missing description for cyaserv
	Port{2584, 17}:   "cyaserv",              // Missing description for cyaserv
	Port{2585, 6}:    "netx-server",          // NETX Server
	Port{2585, 17}:   "netx-server",          // NETX Server
	Port{2586, 6}:    "netx-agent",           // NETX Agent
	Port{2586, 17}:   "netx-agent",           // NETX Agent
	Port{2587, 6}:    "masc",                 // Missing description for masc
	Port{2587, 17}:   "masc",                 // MASC
	Port{2588, 6}:    "privilege",            // Missing description for privilege
	Port{2588, 17}:   "privilege",            // Privilege
	Port{2589, 6}:    "quartus-tcl",          // quartus tcl
	Port{2589, 17}:   "quartus-tcl",          // quartus tcl
	Port{2590, 6}:    "idotdist",             // Missing description for idotdist
	Port{2590, 17}:   "idotdist",             // Missing description for idotdist
	Port{2591, 6}:    "maytagshuffle",        // Maytag Shuffle
	Port{2591, 17}:   "maytagshuffle",        // Maytag Shuffle
	Port{2592, 6}:    "netrek",               // Missing description for netrek
	Port{2592, 17}:   "netrek",               // Missing description for netrek
	Port{2593, 6}:    "mns-mail",             // MNS Mail Notice Service
	Port{2593, 17}:   "mns-mail",             // MNS Mail Notice Service
	Port{2594, 6}:    "dts",                  // Data Base Server
	Port{2594, 17}:   "dts",                  // Data Base Server
	Port{2595, 6}:    "worldfusion1",         // World Fusion 1
	Port{2595, 17}:   "worldfusion1",         // World Fusion 1
	Port{2596, 6}:    "worldfusion2",         // World Fusion 2
	Port{2596, 17}:   "worldfusion2",         // World Fusion 2
	Port{2597, 6}:    "homesteadglory",       // Homestead Glory
	Port{2597, 17}:   "homesteadglory",       // Homestead Glory
	Port{2598, 6}:    "citriximaclient",      // Citrix MA Client
	Port{2598, 17}:   "citriximaclient",      // Citrix MA Client
	Port{2599, 6}:    "snapd",                // Snap Discovery
	Port{2599, 17}:   "snapd",                // Snap Discovery
	Port{2600, 6}:    "zebrasrv",             // hpstgmgr | zebra service | HPSTGMGR
	Port{2600, 17}:   "hpstgmgr",             // HPSTGMGR
	Port{2601, 6}:    "zebra",                // discp-client | zebra vty | discp client
	Port{2601, 17}:   "discp-client",         // discp client
	Port{2602, 6}:    "ripd",                 // discp-server | RIPd vty | discp server
	Port{2602, 17}:   "discp-server",         // discp server
	Port{2603, 6}:    "ripngd",               // servicemeter | RIPngd vty | Service Meter
	Port{2603, 17}:   "servicemeter",         // Service Meter
	Port{2604, 6}:    "ospfd",                // nsc-ccs | OSPFd vty | NSC CCS
	Port{2604, 17}:   "nsc-ccs",              // NSC CCS
	Port{2605, 6}:    "bgpd",                 // nsc-posa | BGPd vty | NSC POSA
	Port{2605, 17}:   "nsc-posa",             // NSC POSA
	Port{2606, 6}:    "netmon",               // Dell Netmon
	Port{2606, 17}:   "netmon",               // Dell Netmon
	Port{2607, 6}:    "connection",           // Dell Connection
	Port{2607, 17}:   "connection",           // Dell Connection
	Port{2608, 6}:    "wag-service",          // Wag Service
	Port{2608, 17}:   "wag-service",          // Wag Service
	Port{2609, 6}:    "system-monitor",       // System Monitor
	Port{2609, 17}:   "system-monitor",       // System Monitor
	Port{2610, 6}:    "versa-tek",            // VersaTek
	Port{2610, 17}:   "versa-tek",            // VersaTek
	Port{2611, 6}:    "lionhead",             // Missing description for lionhead
	Port{2611, 17}:   "lionhead",             // LIONHEAD
	Port{2612, 6}:    "qpasa-agent",          // Qpasa Agent
	Port{2612, 17}:   "qpasa-agent",          // Qpasa Agent
	Port{2613, 6}:    "smntubootstrap",       // Missing description for smntubootstrap
	Port{2613, 17}:   "smntubootstrap",       // SMNTUBootstrap
	Port{2614, 6}:    "neveroffline",         // Never Offline
	Port{2614, 17}:   "neveroffline",         // Never Offline
	Port{2615, 6}:    "firepower",            // Missing description for firepower
	Port{2615, 17}:   "firepower",            // Missing description for firepower
	Port{2616, 6}:    "appswitch-emp",        // Missing description for appswitch-emp
	Port{2616, 17}:   "appswitch-emp",        // Missing description for appswitch-emp
	Port{2617, 6}:    "cmadmin",              // Clinical Context Managers
	Port{2617, 17}:   "cmadmin",              // Clinical Context Managers
	Port{2618, 6}:    "priority-e-com",       // Priority E-Com
	Port{2618, 17}:   "priority-e-com",       // Priority E-Com
	Port{2619, 6}:    "bruce",                // Missing description for bruce
	Port{2619, 17}:   "bruce",                // Missing description for bruce
	Port{2620, 6}:    "lpsrecommender",       // Missing description for lpsrecommender
	Port{2620, 17}:   "lpsrecommender",       // LPSRecommender
	Port{2621, 6}:    "miles-apart",          // Miles Apart Jukebox Server
	Port{2621, 17}:   "miles-apart",          // Miles Apart Jukebox Server
	Port{2622, 6}:    "metricadbc",           // Missing description for metricadbc
	Port{2622, 17}:   "metricadbc",           // MetricaDBC
	Port{2623, 6}:    "lmdp",                 // Missing description for lmdp
	Port{2623, 17}:   "lmdp",                 // LMDP
	Port{2624, 6}:    "aria",                 // Missing description for aria
	Port{2624, 17}:   "aria",                 // Aria
	Port{2625, 6}:    "blwnkl-port",          // Blwnkl Port
	Port{2625, 17}:   "blwnkl-port",          // Blwnkl Port
	Port{2626, 6}:    "gbjd816",              // Missing description for gbjd816
	Port{2626, 17}:   "gbjd816",              // Missing description for gbjd816
	Port{2627, 6}:    "webster",              // moshebeeri | Network dictionary | Moshe Beeri
	Port{2627, 17}:   "webster",              // Missing description for webster
	Port{2628, 6}:    "dict",                 // Dictionary service (RFC2229)
	Port{2628, 17}:   "dict",                 // DICT
	Port{2629, 6}:    "sitaraserver",         // Sitara Server
	Port{2629, 17}:   "sitaraserver",         // Sitara Server
	Port{2630, 6}:    "sitaramgmt",           // Sitara Management
	Port{2630, 17}:   "sitaramgmt",           // Sitara Management
	Port{2631, 6}:    "sitaradir",            // Sitara Dir
	Port{2631, 17}:   "sitaradir",            // Sitara Dir
	Port{2632, 6}:    "irdg-post",            // IRdg Post
	Port{2632, 17}:   "irdg-post",            // IRdg Post
	Port{2633, 6}:    "interintelli",         // Missing description for interintelli
	Port{2633, 17}:   "interintelli",         // InterIntelli
	Port{2634, 6}:    "pk-electronics",       // PK Electronics
	Port{2634, 17}:   "pk-electronics",       // PK Electronics
	Port{2635, 6}:    "backburner",           // Back Burner
	Port{2635, 17}:   "backburner",           // Back Burner
	Port{2636, 6}:    "solve",                // Missing description for solve
	Port{2636, 17}:   "solve",                // Solve
	Port{2637, 6}:    "imdocsvc",             // Import Document Service
	Port{2637, 17}:   "imdocsvc",             // Import Document Service
	Port{2638, 6}:    "sybase",               // sybaseanywhere | Sybase database | Sybase Anywhere
	Port{2638, 17}:   "sybaseanywhere",       // Sybase Anywhere
	Port{2639, 6}:    "aminet",               // Missing description for aminet
	Port{2639, 17}:   "aminet",               // AMInet
	Port{2640, 6}:    "sai_sentlm",           // ami-control | Sabbagh Associates Licence Manager | Alcorn McBride Inc protocol used for device control
	Port{2640, 17}:   "sai_sentlm",           // Sabbagh Associates Licence Manager
	Port{2641, 6}:    "hdl-srv",              // HDL Server
	Port{2641, 17}:   "hdl-srv",              // HDL Server
	Port{2642, 6}:    "tragic",               // Missing description for tragic
	Port{2642, 17}:   "tragic",               // Tragic
	Port{2643, 6}:    "gte-samp",             // Missing description for gte-samp
	Port{2643, 17}:   "gte-samp",             // GTE-SAMP
	Port{2644, 6}:    "travsoft-ipx-t",       // Travsoft IPX Tunnel
	Port{2644, 17}:   "travsoft-ipx-t",       // Travsoft IPX Tunnel
	Port{2645, 6}:    "novell-ipx-cmd",       // Novell IPX CMD
	Port{2645, 17}:   "novell-ipx-cmd",       // Novell IPX CMD
	Port{2646, 6}:    "and-lm",               // AND License Manager
	Port{2646, 17}:   "and-lm",               // AND License Manager
	Port{2647, 6}:    "syncserver",           // Missing description for syncserver
	Port{2647, 17}:   "syncserver",           // SyncServer
	Port{2648, 6}:    "upsnotifyprot",        // Missing description for upsnotifyprot
	Port{2648, 17}:   "upsnotifyprot",        // Upsnotifyprot
	Port{2649, 6}:    "vpsipport",            // Missing description for vpsipport
	Port{2649, 17}:   "vpsipport",            // VPSIPPORT
	Port{2650, 6}:    "eristwoguns",          // Missing description for eristwoguns
	Port{2650, 17}:   "eristwoguns",          // Missing description for eristwoguns
	Port{2651, 6}:    "ebinsite",             // Missing description for ebinsite
	Port{2651, 17}:   "ebinsite",             // EBInSite
	Port{2652, 6}:    "interpathpanel",       // Missing description for interpathpanel
	Port{2652, 17}:   "interpathpanel",       // InterPathPanel
	Port{2653, 6}:    "sonus",                // Missing description for sonus
	Port{2653, 17}:   "sonus",                // Sonus
	Port{2654, 6}:    "corel_vncadmin",       // corel-vncadmin | Corel VNC Admin
	Port{2654, 17}:   "corel_vncadmin",       // Corel VNC Admin
	Port{2655, 6}:    "unglue",               // UNIX Nt Glue
	Port{2655, 17}:   "unglue",               // UNIX Nt Glue
	Port{2656, 6}:    "kana",                 // Missing description for kana
	Port{2656, 17}:   "kana",                 // Kana
	Port{2657, 6}:    "sns-dispatcher",       // SNS Dispatcher
	Port{2657, 17}:   "sns-dispatcher",       // SNS Dispatcher
	Port{2658, 6}:    "sns-admin",            // SNS Admin
	Port{2658, 17}:   "sns-admin",            // SNS Admin
	Port{2659, 6}:    "sns-query",            // SNS Query
	Port{2659, 17}:   "sns-query",            // SNS Query
	Port{2660, 6}:    "gcmonitor",            // GC Monitor
	Port{2660, 17}:   "gcmonitor",            // GC Monitor
	Port{2661, 6}:    "olhost",               // Missing description for olhost
	Port{2661, 17}:   "olhost",               // OLHOST
	Port{2662, 6}:    "bintec-capi",          // Missing description for bintec-capi
	Port{2662, 17}:   "bintec-capi",          // BinTec-CAPI
	Port{2663, 6}:    "bintec-tapi",          // Missing description for bintec-tapi
	Port{2663, 17}:   "bintec-tapi",          // BinTec-TAPI
	Port{2664, 6}:    "patrol-mq-gm",         // Patrol for MQ GM
	Port{2664, 17}:   "patrol-mq-gm",         // Patrol for MQ GM
	Port{2665, 6}:    "patrol-mq-nm",         // Patrol for MQ NM
	Port{2665, 17}:   "patrol-mq-nm",         // Patrol for MQ NM
	Port{2666, 6}:    "extensis",             // Missing description for extensis
	Port{2666, 17}:   "extensis",             // Missing description for extensis
	Port{2667, 6}:    "alarm-clock-s",        // Alarm Clock Server
	Port{2667, 17}:   "alarm-clock-s",        // Alarm Clock Server
	Port{2668, 6}:    "alarm-clock-c",        // Alarm Clock Client
	Port{2668, 17}:   "alarm-clock-c",        // Alarm Clock Client
	Port{2669, 6}:    "toad",                 // Missing description for toad
	Port{2669, 17}:   "toad",                 // TOAD
	Port{2670, 6}:    "tve-announce",         // TVE Announce
	Port{2670, 17}:   "tve-announce",         // TVE Announce
	Port{2671, 6}:    "newlixreg",            // Missing description for newlixreg
	Port{2671, 17}:   "newlixreg",            // Missing description for newlixreg
	Port{2672, 6}:    "nhserver",             // Missing description for nhserver
	Port{2672, 17}:   "nhserver",             // Missing description for nhserver
	Port{2673, 6}:    "firstcall42",          // First Call 42
	Port{2673, 17}:   "firstcall42",          // First Call 42
	Port{2674, 6}:    "ewnn",                 // Missing description for ewnn
	Port{2674, 17}:   "ewnn",                 // Missing description for ewnn
	Port{2675, 6}:    "ttc-etap",             // TTC ETAP
	Port{2675, 17}:   "ttc-etap",             // TTC ETAP
	Port{2676, 6}:    "simslink",             // Missing description for simslink
	Port{2676, 17}:   "simslink",             // SIMSLink
	Port{2677, 6}:    "gadgetgate1way",       // Gadget Gate 1 Way
	Port{2677, 17}:   "gadgetgate1way",       // Gadget Gate 1 Way
	Port{2678, 6}:    "gadgetgate2way",       // Gadget Gate 2 Way
	Port{2678, 17}:   "gadgetgate2way",       // Gadget Gate 2 Way
	Port{2679, 6}:    "syncserverssl",        // Sync Server SSL
	Port{2679, 17}:   "syncserverssl",        // Sync Server SSL
	Port{2680, 6}:    "pxc-sapxom",           // Missing description for pxc-sapxom
	Port{2680, 17}:   "pxc-sapxom",           // Missing description for pxc-sapxom
	Port{2681, 6}:    "mpnjsomb",             // Missing description for mpnjsomb
	Port{2681, 17}:   "mpnjsomb",             // Missing description for mpnjsomb
	Port{2683, 6}:    "ncdloadbalance",       // Missing description for ncdloadbalance
	Port{2683, 17}:   "ncdloadbalance",       // NCDLoadBalance
	Port{2684, 6}:    "mpnjsosv",             // Missing description for mpnjsosv
	Port{2684, 17}:   "mpnjsosv",             // Missing description for mpnjsosv
	Port{2685, 6}:    "mpnjsocl",             // Missing description for mpnjsocl
	Port{2685, 17}:   "mpnjsocl",             // Missing description for mpnjsocl
	Port{2686, 6}:    "mpnjsomg",             // Missing description for mpnjsomg
	Port{2686, 17}:   "mpnjsomg",             // Missing description for mpnjsomg
	Port{2687, 6}:    "pq-lic-mgmt",          // Missing description for pq-lic-mgmt
	Port{2687, 17}:   "pq-lic-mgmt",          // Missing description for pq-lic-mgmt
	Port{2688, 6}:    "md-cg-http",           // md-cf-http
	Port{2688, 17}:   "md-cg-http",           // md-cf-http
	Port{2689, 6}:    "fastlynx",             // Missing description for fastlynx
	Port{2689, 17}:   "fastlynx",             // FastLynx
	Port{2690, 6}:    "hp-nnm-data",          // HP NNM Embedded Database
	Port{2690, 17}:   "hp-nnm-data",          // HP NNM Embedded Database
	Port{2691, 6}:    "itinternet",           // ITInternet ISM Server
	Port{2691, 17}:   "itinternet",           // ITInternet ISM Server
	Port{2692, 6}:    "admins-lms",           // Admins LMS
	Port{2692, 17}:   "admins-lms",           // Admins LMS
	Port{2694, 6}:    "pwrsevent",            // Missing description for pwrsevent
	Port{2694, 17}:   "pwrsevent",            // Missing description for pwrsevent
	Port{2695, 6}:    "vspread",              // Missing description for vspread
	Port{2695, 17}:   "vspread",              // VSPREAD
	Port{2696, 6}:    "unifyadmin",           // Unify Admin
	Port{2696, 17}:   "unifyadmin",           // Unify Admin
	Port{2697, 6}:    "oce-snmp-trap",        // Oce SNMP Trap Port
	Port{2697, 17}:   "oce-snmp-trap",        // Oce SNMP Trap Port
	Port{2698, 6}:    "mck-ivpip",            // Missing description for mck-ivpip
	Port{2698, 17}:   "mck-ivpip",            // MCK-IVPIP
	Port{2699, 6}:    "csoft-plusclnt",       // Csoft Plus Client
	Port{2699, 17}:   "csoft-plusclnt",       // Csoft Plus Client
	Port{2700, 6}:    "tqdata",               // Missing description for tqdata
	Port{2700, 17}:   "tqdata",               // Missing description for tqdata
	Port{2701, 6}:    "sms-rcinfo",           // SMS RCINFO
	Port{2701, 17}:   "sms-rcinfo",           // Missing description for sms-rcinfo
	Port{2702, 6}:    "sms-xfer",             // SMS XFER
	Port{2702, 17}:   "sms-xfer",             // Missing description for sms-xfer
	Port{2703, 6}:    "sms-chat",             // SMS CHAT
	Port{2703, 17}:   "sms-chat",             // SMS CHAT
	Port{2704, 6}:    "sms-remctrl",          // SMS REMCTRL
	Port{2704, 17}:   "sms-remctrl",          // SMS REMCTRL
	Port{2705, 6}:    "sds-admin",            // SDS Admin
	Port{2705, 17}:   "sds-admin",            // SDS Admin
	Port{2706, 6}:    "ncdmirroring",         // NCD Mirroring
	Port{2706, 17}:   "ncdmirroring",         // NCD Mirroring
	Port{2707, 6}:    "emcsymapiport",        // Missing description for emcsymapiport
	Port{2707, 17}:   "emcsymapiport",        // EMCSYMAPIPORT
	Port{2708, 6}:    "banyan-net",           // Missing description for banyan-net
	Port{2708, 17}:   "banyan-net",           // Banyan-Net
	Port{2709, 6}:    "supermon",             // Missing description for supermon
	Port{2709, 17}:   "supermon",             // Supermon
	Port{2710, 6}:    "sso-service",          // SSO Service
	Port{2710, 17}:   "sso-service",          // SSO Service
	Port{2711, 6}:    "sso-control",          // SSO Control
	Port{2711, 17}:   "sso-control",          // SSO Control
	Port{2712, 6}:    "aocp",                 // Axapta Object Communication Protocol
	Port{2712, 17}:   "aocp",                 // Axapta Object Communication Protocol
	Port{2713, 6}:    "raventbs",             // Raven Trinity Broker Service
	Port{2713, 17}:   "raventbs",             // Raven Trinity Broker Service
	Port{2714, 6}:    "raventdm",             // Raven Trinity Data Mover
	Port{2714, 17}:   "raventdm",             // Raven Trinity Data Mover
	Port{2715, 6}:    "hpstgmgr2",            // Missing description for hpstgmgr2
	Port{2715, 17}:   "hpstgmgr2",            // HPSTGMGR2
	Port{2716, 6}:    "inova-ip-disco",       // Inova IP Disco
	Port{2716, 17}:   "inova-ip-disco",       // Inova IP Disco
	Port{2717, 6}:    "pn-requester",         // PN REQUESTER
	Port{2717, 17}:   "pn-requester",         // PN REQUESTER
	Port{2718, 6}:    "pn-requester2",        // PN REQUESTER 2
	Port{2718, 17}:   "pn-requester2",        // PN REQUESTER 2
	Port{2719, 6}:    "scan-change",          // Scan & Change
	Port{2719, 17}:   "scan-change",          // Scan & Change
	Port{2720, 6}:    "wkars",                // Missing description for wkars
	Port{2720, 17}:   "wkars",                // Missing description for wkars
	Port{2721, 6}:    "smart-diagnose",       // Smart Diagnose
	Port{2721, 17}:   "smart-diagnose",       // Smart Diagnose
	Port{2722, 6}:    "proactivesrvr",        // Proactive Server
	Port{2722, 17}:   "proactivesrvr",        // Proactive Server
	Port{2723, 6}:    "watchdog-nt",          // WatchDog NT Protocol
	Port{2723, 17}:   "watchdog-nt",          // WatchDog NT Protocol
	Port{2724, 6}:    "qotps",                // Missing description for qotps
	Port{2724, 17}:   "qotps",                // Missing description for qotps
	Port{2725, 6}:    "msolap-ptp2",          // SQL Analysis Server | MSOLAP PTP2
	Port{2725, 17}:   "msolap-ptp2",          // MSOLAP PTP2
	Port{2726, 6}:    "tams",                 // Missing description for tams
	Port{2726, 17}:   "tams",                 // TAMS
	Port{2727, 6}:    "mgcp-callagent",       // Media Gateway Control Protocol Call Agent
	Port{2727, 17}:   "mgcp-callagent",       // Media Gateway Control Protocol Call Agent
	Port{2728, 6}:    "sqdr",                 // Missing description for sqdr
	Port{2728, 17}:   "sqdr",                 // SQDR
	Port{2729, 6}:    "tcim-control",         // TCIM Control
	Port{2729, 17}:   "tcim-control",         // TCIM Control
	Port{2730, 6}:    "nec-raidplus",         // NEC RaidPlus
	Port{2730, 17}:   "nec-raidplus",         // NEC RaidPlus
	Port{2731, 6}:    "fyre-messanger",       // Fyre Messanger | Fyre Messagner
	Port{2731, 17}:   "fyre-messanger",       // Fyre Messagner
	Port{2732, 6}:    "g5m",                  // Missing description for g5m
	Port{2732, 17}:   "g5m",                  // G5M
	Port{2733, 6}:    "signet-ctf",           // Signet CTF
	Port{2733, 17}:   "signet-ctf",           // Signet CTF
	Port{2734, 6}:    "ccs-software",         // CCS Software
	Port{2734, 17}:   "ccs-software",         // CCS Software
	Port{2735, 6}:    "netiq-mc",             // NetIQ Monitor Console
	Port{2735, 17}:   "netiq-mc",             // NetIQ Monitor Console
	Port{2736, 6}:    "radwiz-nms-srv",       // RADWIZ NMS SRV
	Port{2736, 17}:   "radwiz-nms-srv",       // RADWIZ NMS SRV
	Port{2737, 6}:    "srp-feedback",         // SRP Feedback
	Port{2737, 17}:   "srp-feedback",         // SRP Feedback
	Port{2738, 6}:    "ndl-tcp-ois-gw",       // NDL TCP-OSI Gateway
	Port{2738, 17}:   "ndl-tcp-ois-gw",       // NDL TCP-OSI Gateway
	Port{2739, 6}:    "tn-timing",            // TN Timing
	Port{2739, 17}:   "tn-timing",            // TN Timing
	Port{2740, 6}:    "alarm",                // Missing description for alarm
	Port{2740, 17}:   "alarm",                // Alarm
	Port{2741, 6}:    "tsb",                  // Missing description for tsb
	Port{2741, 17}:   "tsb",                  // TSB
	Port{2742, 6}:    "tsb2",                 // Missing description for tsb2
	Port{2742, 17}:   "tsb2",                 // TSB2
	Port{2743, 6}:    "murx",                 // Missing description for murx
	Port{2743, 17}:   "murx",                 // Missing description for murx
	Port{2744, 6}:    "honyaku",              // Missing description for honyaku
	Port{2744, 17}:   "honyaku",              // Missing description for honyaku
	Port{2745, 6}:    "urbisnet",             // Missing description for urbisnet
	Port{2745, 17}:   "urbisnet",             // URBISNET
	Port{2746, 6}:    "cpudpencap",           // Missing description for cpudpencap
	Port{2746, 17}:   "cpudpencap",           // CPUDPENCAP
	Port{2747, 6}:    "fjippol-swrly",        // Missing description for fjippol-swrly
	Port{2747, 17}:   "fjippol-swrly",        // Missing description for fjippol-swrly
	Port{2748, 6}:    "fjippol-polsvr",       // Missing description for fjippol-polsvr
	Port{2748, 17}:   "fjippol-polsvr",       // Missing description for fjippol-polsvr
	Port{2749, 6}:    "fjippol-cnsl",         // Missing description for fjippol-cnsl
	Port{2749, 17}:   "fjippol-cnsl",         // Missing description for fjippol-cnsl
	Port{2750, 6}:    "fjippol-port1",        // Missing description for fjippol-port1
	Port{2750, 17}:   "fjippol-port1",        // Missing description for fjippol-port1
	Port{2751, 6}:    "fjippol-port2",        // Missing description for fjippol-port2
	Port{2751, 17}:   "fjippol-port2",        // Missing description for fjippol-port2
	Port{2752, 6}:    "rsisysaccess",         // RSISYS ACCESS
	Port{2752, 17}:   "rsisysaccess",         // RSISYS ACCESS
	Port{2753, 6}:    "de-spot",              // Missing description for de-spot
	Port{2753, 17}:   "de-spot",              // Missing description for de-spot
	Port{2754, 6}:    "apollo-cc",            // APOLLO CC
	Port{2754, 17}:   "apollo-cc",            // APOLLO CC
	Port{2755, 6}:    "expresspay",           // Express Pay
	Port{2755, 17}:   "expresspay",           // Express Pay
	Port{2756, 6}:    "simplement-tie",       // Missing description for simplement-tie
	Port{2756, 17}:   "simplement-tie",       // Missing description for simplement-tie
	Port{2757, 6}:    "cnrp",                 // Missing description for cnrp
	Port{2757, 17}:   "cnrp",                 // CNRP
	Port{2758, 6}:    "apollo-status",        // APOLLO Status
	Port{2758, 17}:   "apollo-status",        // APOLLO Status
	Port{2759, 6}:    "apollo-gms",           // APOLLO GMS
	Port{2759, 17}:   "apollo-gms",           // APOLLO GMS
	Port{2760, 6}:    "sabams",               // Saba MS
	Port{2760, 17}:   "sabams",               // Saba MS
	Port{2761, 6}:    "dicom-iscl",           // DICOM ISCL
	Port{2761, 17}:   "dicom-iscl",           // DICOM ISCL
	Port{2762, 6}:    "dicom-tls",            // DICOM TLS
	Port{2762, 17}:   "dicom-tls",            // DICOM TLS
	Port{2763, 6}:    "desktop-dna",          // Desktop DNA
	Port{2763, 17}:   "desktop-dna",          // Desktop DNA
	Port{2764, 6}:    "data-insurance",       // Data Insurance
	Port{2764, 17}:   "data-insurance",       // Data Insurance
	Port{2765, 6}:    "qip-audup",            // Missing description for qip-audup
	Port{2765, 17}:   "qip-audup",            // Missing description for qip-audup
	Port{2766, 6}:    "listen",               // compaq-scp | System V listener port | Compaq SCP
	Port{2766, 17}:   "compaq-scp",           // Compaq SCP
	Port{2767, 6}:    "uadtc",                // Missing description for uadtc
	Port{2767, 17}:   "uadtc",                // UADTC
	Port{2768, 6}:    "uacs",                 // Missing description for uacs
	Port{2768, 17}:   "uacs",                 // UACS
	Port{2769, 6}:    "exce",                 // Missing description for exce
	Port{2769, 17}:   "exce",                 // eXcE
	Port{2770, 6}:    "veronica",             // Missing description for veronica
	Port{2770, 17}:   "veronica",             // Veronica
	Port{2771, 6}:    "vergencecm",           // Vergence CM
	Port{2771, 17}:   "vergencecm",           // Vergence CM
	Port{2772, 6}:    "auris",                // Missing description for auris
	Port{2772, 17}:   "auris",                // Missing description for auris
	Port{2773, 6}:    "rbakcup1",             // RBackup Remote Backup
	Port{2773, 17}:   "rbakcup1",             // RBackup Remote Backup
	Port{2774, 6}:    "rbakcup2",             // RBackup Remote Backup
	Port{2774, 17}:   "rbakcup2",             // RBackup Remote Backup
	Port{2775, 6}:    "smpp",                 // Missing description for smpp
	Port{2775, 17}:   "smpp",                 // SMPP
	Port{2776, 6}:    "ridgeway1",            // Ridgeway Systems & Software
	Port{2776, 17}:   "ridgeway1",            // Ridgeway Systems & Software
	Port{2777, 6}:    "ridgeway2",            // Ridgeway Systems & Software
	Port{2777, 17}:   "ridgeway2",            // Ridgeway Systems & Software
	Port{2778, 6}:    "gwen-sonya",           // Missing description for gwen-sonya
	Port{2778, 17}:   "gwen-sonya",           // Gwen-Sonya
	Port{2779, 6}:    "lbc-sync",             // LBC Sync
	Port{2779, 17}:   "lbc-sync",             // LBC Sync
	Port{2780, 6}:    "lbc-control",          // LBC Control
	Port{2780, 17}:   "lbc-control",          // LBC Control
	Port{2781, 6}:    "whosells",             // Missing description for whosells
	Port{2781, 17}:   "whosells",             // Missing description for whosells
	Port{2782, 6}:    "everydayrc",           // Missing description for everydayrc
	Port{2782, 17}:   "everydayrc",           // Missing description for everydayrc
	Port{2783, 6}:    "aises",                // Missing description for aises
	Port{2783, 17}:   "aises",                // AISES
	Port{2784, 6}:    "www-dev",              // world wide web - development
	Port{2784, 17}:   "www-dev",              // world wide web - development
	Port{2785, 6}:    "aic-np",               // Missing description for aic-np
	Port{2785, 17}:   "aic-np",               // Missing description for aic-np
	Port{2786, 6}:    "aic-oncrpc",           // aic-oncrpc - Destiny MCD database
	Port{2786, 17}:   "aic-oncrpc",           // aic-oncrpc - Destiny MCD database
	Port{2787, 6}:    "piccolo",              // piccolo - Cornerstone Software
	Port{2787, 17}:   "piccolo",              // piccolo - Cornerstone Software
	Port{2788, 6}:    "fryeserv",             // NetWare Loadable Module - Seagate Software
	Port{2788, 17}:   "fryeserv",             // NetWare Loadable Module - Seagate Software
	Port{2789, 6}:    "media-agent",          // Media Agent
	Port{2789, 17}:   "media-agent",          // Media Agent
	Port{2790, 6}:    "plgproxy",             // PLG Proxy
	Port{2790, 17}:   "plgproxy",             // PLG Proxy
	Port{2791, 6}:    "mtport-regist",        // MT Port Registrator
	Port{2791, 17}:   "mtport-regist",        // MT Port Registrator
	Port{2792, 6}:    "f5-globalsite",        // Missing description for f5-globalsite
	Port{2792, 17}:   "f5-globalsite",        // Missing description for f5-globalsite
	Port{2793, 6}:    "initlsmsad",           // Missing description for initlsmsad
	Port{2793, 17}:   "initlsmsad",           // Missing description for initlsmsad
	Port{2795, 6}:    "livestats",            // Missing description for livestats
	Port{2795, 17}:   "livestats",            // LiveStats
	Port{2796, 6}:    "ac-tech",              // Missing description for ac-tech
	Port{2796, 17}:   "ac-tech",              // Missing description for ac-tech
	Port{2797, 6}:    "esp-encap",            // Missing description for esp-encap
	Port{2797, 17}:   "esp-encap",            // Missing description for esp-encap
	Port{2798, 6}:    "tmesis-upshot",        // Missing description for tmesis-upshot
	Port{2798, 17}:   "tmesis-upshot",        // TMESIS-UPShot
	Port{2799, 6}:    "icon-discover",        // ICON Discover
	Port{2799, 17}:   "icon-discover",        // ICON Discover
	Port{2800, 6}:    "acc-raid",             // ACC RAID
	Port{2800, 17}:   "acc-raid",             // ACC RAID
	Port{2801, 6}:    "igcp",                 // Missing description for igcp
	Port{2801, 17}:   "igcp",                 // IGCP
	Port{2802, 6}:    "veritas-tcp1",         // veritas-udp1 | Veritas TCP1 | Veritas UDP1
	Port{2802, 17}:   "veritas-udp1",         // Veritas UDP1
	Port{2803, 6}:    "btprjctrl",            // Missing description for btprjctrl
	Port{2803, 17}:   "btprjctrl",            // Missing description for btprjctrl
	Port{2804, 6}:    "dvr-esm",              // March Networks Digital Video Recorders and Enterprise Service Manager products
	Port{2804, 17}:   "dvr-esm",              // March Networks Digital Video Recorders and Enterprise Service Manager products
	Port{2805, 6}:    "wta-wsp-s",            // WTA WSP-S
	Port{2805, 17}:   "wta-wsp-s",            // WTA WSP-S
	Port{2806, 6}:    "cspuni",               // Missing description for cspuni
	Port{2806, 17}:   "cspuni",               // Missing description for cspuni
	Port{2807, 6}:    "cspmulti",             // Missing description for cspmulti
	Port{2807, 17}:   "cspmulti",             // Missing description for cspmulti
	Port{2808, 6}:    "j-lan-p",              // Missing description for j-lan-p
	Port{2808, 17}:   "j-lan-p",              // J-LAN-P
	Port{2809, 6}:    "corbaloc",             // Corba | CORBA LOC
	Port{2809, 17}:   "corbaloc",             // CORBA LOC
	Port{2810, 6}:    "netsteward",           // Active Net Steward
	Port{2810, 17}:   "netsteward",           // Active Net Steward
	Port{2811, 6}:    "gsiftp",               // GSI FTP
	Port{2811, 17}:   "gsiftp",               // GSI FTP
	Port{2812, 6}:    "atmtcp",               // commonly Monit httpd - http:  mmonit.com monit documentation monit.html#monit_httpd
	Port{2812, 17}:   "atmtcp",               // Missing description for atmtcp
	Port{2813, 6}:    "llm-pass",             // Missing description for llm-pass
	Port{2813, 17}:   "llm-pass",             // Missing description for llm-pass
	Port{2814, 6}:    "llm-csv",              // Missing description for llm-csv
	Port{2814, 17}:   "llm-csv",              // Missing description for llm-csv
	Port{2815, 6}:    "lbc-measure",          // LBC Measurement
	Port{2815, 17}:   "lbc-measure",          // LBC Measurement
	Port{2816, 6}:    "lbc-watchdog",         // LBC Watchdog
	Port{2816, 17}:   "lbc-watchdog",         // LBC Watchdog
	Port{2817, 6}:    "nmsigport",            // NMSig Port
	Port{2817, 17}:   "nmsigport",            // NMSig Port
	Port{2818, 6}:    "rmlnk",                // Missing description for rmlnk
	Port{2818, 17}:   "rmlnk",                // Missing description for rmlnk
	Port{2819, 6}:    "fc-faultnotify",       // FC Fault Notification
	Port{2819, 17}:   "fc-faultnotify",       // FC Fault Notification
	Port{2820, 6}:    "univision",            // Missing description for univision
	Port{2820, 17}:   "univision",            // UniVision
	Port{2821, 6}:    "vrts-at-port",         // VERITAS Authentication Service
	Port{2821, 17}:   "vrts-at-port",         // VERITAS Authentication Service
	Port{2822, 6}:    "ka0wuc",               // Missing description for ka0wuc
	Port{2822, 17}:   "ka0wuc",               // Missing description for ka0wuc
	Port{2823, 6}:    "cqg-netlan",           // CQG Net LAN
	Port{2823, 17}:   "cqg-netlan",           // CQG Net LAN
	Port{2824, 6}:    "cqg-netlan-1",         // CQG Net LAN 1 | CQG Net Lan 1
	Port{2824, 17}:   "cqg-netlan-1",         // CQG Net Lan 1
	Port{2826, 6}:    "slc-systemlog",        // slc systemlog
	Port{2826, 17}:   "slc-systemlog",        // slc systemlog
	Port{2827, 6}:    "slc-ctrlrloops",       // slc ctrlrloops
	Port{2827, 17}:   "slc-ctrlrloops",       // slc ctrlrloops
	Port{2828, 6}:    "itm-lm",               // ITM License Manager
	Port{2828, 17}:   "itm-lm",               // ITM License Manager
	Port{2829, 6}:    "silkp1",               // Missing description for silkp1
	Port{2829, 17}:   "silkp1",               // Missing description for silkp1
	Port{2830, 6}:    "silkp2",               // Missing description for silkp2
	Port{2830, 17}:   "silkp2",               // Missing description for silkp2
	Port{2831, 6}:    "silkp3",               // Missing description for silkp3
	Port{2831, 17}:   "silkp3",               // Missing description for silkp3
	Port{2832, 6}:    "silkp4",               // Missing description for silkp4
	Port{2832, 17}:   "silkp4",               // Missing description for silkp4
	Port{2833, 6}:    "glishd",               // Missing description for glishd
	Port{2833, 17}:   "glishd",               // Missing description for glishd
	Port{2834, 6}:    "evtp",                 // Missing description for evtp
	Port{2834, 17}:   "evtp",                 // EVTP
	Port{2835, 6}:    "evtp-data",            // Missing description for evtp-data
	Port{2835, 17}:   "evtp-data",            // EVTP-DATA
	Port{2836, 6}:    "catalyst",             // Missing description for catalyst
	Port{2836, 17}:   "catalyst",             // Missing description for catalyst
	Port{2837, 6}:    "repliweb",             // Missing description for repliweb
	Port{2837, 17}:   "repliweb",             // Repliweb
	Port{2838, 6}:    "starbot",              // Missing description for starbot
	Port{2838, 17}:   "starbot",              // Starbot
	Port{2839, 6}:    "nmsigport",            // Missing description for nmsigport
	Port{2839, 17}:   "nmsigport",            // NMSigPort
	Port{2840, 6}:    "l3-exprt",             // Missing description for l3-exprt
	Port{2840, 17}:   "l3-exprt",             // Missing description for l3-exprt
	Port{2841, 6}:    "l3-ranger",            // Missing description for l3-ranger
	Port{2841, 17}:   "l3-ranger",            // Missing description for l3-ranger
	Port{2842, 6}:    "l3-hawk",              // Missing description for l3-hawk
	Port{2842, 17}:   "l3-hawk",              // Missing description for l3-hawk
	Port{2843, 6}:    "pdnet",                // Missing description for pdnet
	Port{2843, 17}:   "pdnet",                // PDnet
	Port{2844, 6}:    "bpcp-poll",            // BPCP POLL
	Port{2844, 17}:   "bpcp-poll",            // BPCP POLL
	Port{2845, 6}:    "bpcp-trap",            // BPCP TRAP
	Port{2845, 17}:   "bpcp-trap",            // BPCP TRAP
	Port{2846, 6}:    "aimpp-hello",          // AIMPP Hello
	Port{2846, 17}:   "aimpp-hello",          // AIMPP Hello
	Port{2847, 6}:    "aimpp-port-req",       // AIMPP Port Req
	Port{2847, 17}:   "aimpp-port-req",       // AIMPP Port Req
	Port{2848, 6}:    "amt-blc-port",         // Missing description for amt-blc-port
	Port{2848, 17}:   "amt-blc-port",         // AMT-BLC-PORT
	Port{2849, 6}:    "fxp",                  // Missing description for fxp
	Port{2849, 17}:   "fxp",                  // FXP
	Port{2850, 6}:    "metaconsole",          // Missing description for metaconsole
	Port{2850, 17}:   "metaconsole",          // MetaConsole
	Port{2851, 6}:    "webemshttp",           // Missing description for webemshttp
	Port{2851, 17}:   "webemshttp",           // Missing description for webemshttp
	Port{2852, 6}:    "bears-01",             // Missing description for bears-01
	Port{2852, 17}:   "bears-01",             // Missing description for bears-01
	Port{2853, 6}:    "ispipes",              // Missing description for ispipes
	Port{2853, 17}:   "ispipes",              // ISPipes
	Port{2854, 6}:    "infomover",            // Missing description for infomover
	Port{2854, 17}:   "infomover",            // InfoMover
	Port{2855, 6}:    "msrp",                 // MSRP over TCP
	Port{2855, 17}:   "msrp",                 // MSRP
	Port{2856, 6}:    "cesdinv",              // Missing description for cesdinv
	Port{2856, 17}:   "cesdinv",              // Missing description for cesdinv
	Port{2857, 6}:    "simctlp",              // SimCtIP
	Port{2857, 17}:   "simctlp",              // SimCtIP
	Port{2858, 6}:    "ecnp",                 // Missing description for ecnp
	Port{2858, 17}:   "ecnp",                 // ECNP
	Port{2859, 6}:    "activememory",         // Active Memory
	Port{2859, 17}:   "activememory",         // Active Memory
	Port{2860, 6}:    "dialpad-voice1",       // Dialpad Voice 1
	Port{2860, 17}:   "dialpad-voice1",       // Dialpad Voice 1
	Port{2861, 6}:    "dialpad-voice2",       // Dialpad Voice 2
	Port{2861, 17}:   "dialpad-voice2",       // Dialpad Voice 2
	Port{2862, 6}:    "ttg-protocol",         // TTG Protocol
	Port{2862, 17}:   "ttg-protocol",         // TTG Protocol
	Port{2863, 6}:    "sonardata",            // Sonar Data
	Port{2863, 17}:   "sonardata",            // Sonar Data
	Port{2864, 6}:    "astromed-main",        // main 5001 cmd
	Port{2864, 17}:   "astromed-main",        // main 5001 cmd
	Port{2865, 6}:    "pit-vpn",              // Missing description for pit-vpn
	Port{2865, 17}:   "pit-vpn",              // Missing description for pit-vpn
	Port{2866, 6}:    "iwlistener",           // Missing description for iwlistener
	Port{2866, 17}:   "iwlistener",           // Missing description for iwlistener
	Port{2867, 6}:    "esps-portal",          // Missing description for esps-portal
	Port{2867, 17}:   "esps-portal",          // Missing description for esps-portal
	Port{2868, 6}:    "npep-messaging",       // NPEP Messaging | Norman Proprietaqry Events Protocol
	Port{2868, 17}:   "npep-messaging",       // NPEP Messaging
	Port{2869, 6}:    "icslap",               // Universal Plug and Play Device Host, SSDP Discovery Service
	Port{2869, 17}:   "icslap",               // ICSLAP
	Port{2870, 6}:    "daishi",               // Missing description for daishi
	Port{2870, 17}:   "daishi",               // Missing description for daishi
	Port{2871, 6}:    "msi-selectplay",       // MSI Select Play
	Port{2871, 17}:   "msi-selectplay",       // MSI Select Play
	Port{2872, 6}:    "radix",                // Missing description for radix
	Port{2872, 17}:   "radix",                // RADIX
	Port{2874, 6}:    "dxmessagebase1",       // DX Message Base Transport Protocol
	Port{2874, 17}:   "dxmessagebase1",       // DX Message Base Transport Protocol
	Port{2875, 6}:    "dxmessagebase2",       // DX Message Base Transport Protocol
	Port{2875, 17}:   "dxmessagebase2",       // DX Message Base Transport Protocol
	Port{2876, 6}:    "sps-tunnel",           // SPS Tunnel
	Port{2876, 17}:   "sps-tunnel",           // SPS Tunnel
	Port{2877, 6}:    "bluelance",            // Missing description for bluelance
	Port{2877, 17}:   "bluelance",            // BLUELANCE
	Port{2878, 6}:    "aap",                  // Missing description for aap
	Port{2878, 17}:   "aap",                  // AAP
	Port{2879, 6}:    "ucentric-ds",          // Missing description for ucentric-ds
	Port{2879, 17}:   "ucentric-ds",          // Missing description for ucentric-ds
	Port{2880, 6}:    "synapse",              // Synapse Transport
	Port{2880, 17}:   "synapse",              // Synapse Transport
	Port{2881, 6}:    "ndsp",                 // Missing description for ndsp
	Port{2881, 17}:   "ndsp",                 // NDSP
	Port{2882, 6}:    "ndtp",                 // Missing description for ndtp
	Port{2882, 17}:   "ndtp",                 // NDTP
	Port{2883, 6}:    "ndnp",                 // Missing description for ndnp
	Port{2883, 17}:   "ndnp",                 // NDNP
	Port{2884, 6}:    "flashmsg",             // Flash Msg
	Port{2884, 17}:   "flashmsg",             // Flash Msg
	Port{2885, 6}:    "topflow",              // Missing description for topflow
	Port{2885, 17}:   "topflow",              // TopFlow
	Port{2886, 6}:    "responselogic",        // Missing description for responselogic
	Port{2886, 17}:   "responselogic",        // RESPONSELOGIC
	Port{2887, 6}:    "aironetddp",           // aironet
	Port{2887, 17}:   "aironetddp",           // aironet
	Port{2888, 6}:    "spcsdlobby",           // Missing description for spcsdlobby
	Port{2888, 17}:   "spcsdlobby",           // SPCSDLOBBY
	Port{2889, 6}:    "rsom",                 // Missing description for rsom
	Port{2889, 17}:   "rsom",                 // RSOM
	Port{2890, 6}:    "cspclmulti",           // Missing description for cspclmulti
	Port{2890, 17}:   "cspclmulti",           // CSPCLMULTI
	Port{2891, 6}:    "cinegrfx-elmd",        // CINEGRFX-ELMD License Manager
	Port{2891, 17}:   "cinegrfx-elmd",        // CINEGRFX-ELMD License Manager
	Port{2892, 6}:    "snifferdata",          // Missing description for snifferdata
	Port{2892, 17}:   "snifferdata",          // SNIFFERDATA
	Port{2893, 6}:    "vseconnector",         // Missing description for vseconnector
	Port{2893, 17}:   "vseconnector",         // VSECONNECTOR
	Port{2894, 6}:    "abacus-remote",        // Missing description for abacus-remote
	Port{2894, 17}:   "abacus-remote",        // ABACUS-REMOTE
	Port{2895, 6}:    "natuslink",            // NATUS LINK
	Port{2895, 17}:   "natuslink",            // NATUS LINK
	Port{2896, 6}:    "ecovisiong6-1",        // Missing description for ecovisiong6-1
	Port{2896, 17}:   "ecovisiong6-1",        // ECOVISIONG6-1
	Port{2897, 6}:    "citrix-rtmp",          // Citrix RTMP
	Port{2897, 17}:   "citrix-rtmp",          // Citrix RTMP
	Port{2898, 6}:    "appliance-cfg",        // Missing description for appliance-cfg
	Port{2898, 17}:   "appliance-cfg",        // APPLIANCE-CFG
	Port{2899, 6}:    "powergemplus",         // Missing description for powergemplus
	Port{2899, 17}:   "powergemplus",         // POWERGEMPLUS
	Port{2900, 6}:    "quicksuite",           // Missing description for quicksuite
	Port{2900, 17}:   "quicksuite",           // QUICKSUITE
	Port{2901, 6}:    "allstorcns",           // Missing description for allstorcns
	Port{2901, 17}:   "allstorcns",           // ALLSTORCNS
	Port{2902, 6}:    "netaspi",              // NET ASPI
	Port{2902, 17}:   "netaspi",              // NET ASPI
	Port{2903, 6}:    "extensisportfolio",    // suitcase | Portfolio Server by Extensis Product Group | SUITCASE
	Port{2903, 17}:   "suitcase",             // SUITCASE
	Port{2904, 132}:  "m2ua",                 // SIGTRAN M2UA
	Port{2904, 6}:    "m2ua",                 // SIGTRAN M2UA
	Port{2904, 17}:   "m2ua",                 // SIGTRAN M2UA
	Port{2905, 132}:  "m3ua",                 // SIGTRAN M3UA
	Port{2905, 6}:    "m3ua",                 // SIGTRAN M3UA
	Port{2905, 17}:   "m3ua",                 // SIGTRAN M3UA
	Port{2906, 6}:    "caller9",              // Missing description for caller9
	Port{2906, 17}:   "caller9",              // CALLER9
	Port{2907, 6}:    "webmethods-b2b",       // WEBMETHODS B2B
	Port{2907, 17}:   "webmethods-b2b",       // WEBMETHODS B2B
	Port{2908, 6}:    "mao",                  // Missing description for mao
	Port{2908, 17}:   "mao",                  // Missing description for mao
	Port{2909, 6}:    "funk-dialout",         // Funk Dialout
	Port{2909, 17}:   "funk-dialout",         // Funk Dialout
	Port{2910, 6}:    "tdaccess",             // Missing description for tdaccess
	Port{2910, 17}:   "tdaccess",             // TDAccess
	Port{2911, 6}:    "blockade",             // Missing description for blockade
	Port{2911, 17}:   "blockade",             // Blockade
	Port{2912, 6}:    "epicon",               // Missing description for epicon
	Port{2912, 17}:   "epicon",               // Epicon
	Port{2913, 6}:    "boosterware",          // Booster Ware
	Port{2913, 17}:   "boosterware",          // Booster Ware
	Port{2914, 6}:    "gamelobby",            // Game Lobby
	Port{2914, 17}:   "gamelobby",            // Game Lobby
	Port{2915, 6}:    "tksocket",             // TK Socket
	Port{2915, 17}:   "tksocket",             // TK Socket
	Port{2916, 6}:    "elvin_server",         // elvin-server | Elvin Server
	Port{2916, 17}:   "elvin_server",         // Elvin Server
	Port{2917, 6}:    "elvin_client",         // elvin-client | Elvin Client
	Port{2917, 17}:   "elvin_client",         // Elvin Client
	Port{2918, 6}:    "kastenchasepad",       // Kasten Chase Pad
	Port{2918, 17}:   "kastenchasepad",       // Kasten Chase Pad
	Port{2919, 6}:    "roboer",               // Missing description for roboer
	Port{2919, 17}:   "roboer",               // roboER
	Port{2920, 6}:    "roboeda",              // Missing description for roboeda
	Port{2920, 17}:   "roboeda",              // roboEDA
	Port{2921, 6}:    "cesdcdman",            // CESD Contents Delivery Management
	Port{2921, 17}:   "cesdcdman",            // CESD Contents Delivery Management
	Port{2922, 6}:    "cesdcdtrn",            // CESD Contents Delivery Data Transfer
	Port{2922, 17}:   "cesdcdtrn",            // CESD Contents Delivery Data Transfer
	Port{2923, 6}:    "wta-wsp-wtp-s",        // Missing description for wta-wsp-wtp-s
	Port{2923, 17}:   "wta-wsp-wtp-s",        // WTA-WSP-WTP-S
	Port{2924, 6}:    "precise-vip",          // Missing description for precise-vip
	Port{2924, 17}:   "precise-vip",          // PRECISE-VIP
	Port{2926, 6}:    "mobile-file-dl",       // Missing description for mobile-file-dl
	Port{2926, 17}:   "mobile-file-dl",       // MOBILE-FILE-DL
	Port{2927, 6}:    "unimobilectrl",        // Missing description for unimobilectrl
	Port{2927, 17}:   "unimobilectrl",        // UNIMOBILECTRL
	Port{2928, 6}:    "redstone-cpss",        // Missing description for redstone-cpss
	Port{2928, 17}:   "redstone-cpss",        // REDSTONE-CPSS
	Port{2929, 6}:    "amx-webadmin",         // Missing description for amx-webadmin
	Port{2929, 17}:   "amx-webadmin",         // AMX-WEBADMIN
	Port{2930, 6}:    "amx-weblinx",          // Missing description for amx-weblinx
	Port{2930, 17}:   "amx-weblinx",          // AMX-WEBLINX
	Port{2931, 6}:    "circle-x",             // Missing description for circle-x
	Port{2931, 17}:   "circle-x",             // Circle-X
	Port{2932, 6}:    "incp",                 // Missing description for incp
	Port{2932, 17}:   "incp",                 // INCP
	Port{2933, 6}:    "4-tieropmgw",          // 4-TIER OPM GW
	Port{2933, 17}:   "4-tieropmgw",          // 4-TIER OPM GW
	Port{2934, 6}:    "4-tieropmcli",         // 4-TIER OPM CLI
	Port{2934, 17}:   "4-tieropmcli",         // 4-TIER OPM CLI
	Port{2935, 6}:    "qtp",                  // Missing description for qtp
	Port{2935, 17}:   "qtp",                  // QTP
	Port{2936, 6}:    "otpatch",              // Missing description for otpatch
	Port{2936, 17}:   "otpatch",              // OTPatch
	Port{2937, 6}:    "pnaconsult-lm",        // Missing description for pnaconsult-lm
	Port{2937, 17}:   "pnaconsult-lm",        // PNACONSULT-LM
	Port{2938, 6}:    "sm-pas-1",             // Missing description for sm-pas-1
	Port{2938, 17}:   "sm-pas-1",             // SM-PAS-1
	Port{2939, 6}:    "sm-pas-2",             // Missing description for sm-pas-2
	Port{2939, 17}:   "sm-pas-2",             // SM-PAS-2
	Port{2940, 6}:    "sm-pas-3",             // Missing description for sm-pas-3
	Port{2940, 17}:   "sm-pas-3",             // SM-PAS-3
	Port{2941, 6}:    "sm-pas-4",             // Missing description for sm-pas-4
	Port{2941, 17}:   "sm-pas-4",             // SM-PAS-4
	Port{2942, 6}:    "sm-pas-5",             // Missing description for sm-pas-5
	Port{2942, 17}:   "sm-pas-5",             // SM-PAS-5
	Port{2943, 6}:    "ttnrepository",        // Missing description for ttnrepository
	Port{2943, 17}:   "ttnrepository",        // TTNRepository
	Port{2944, 132}:  "megaco-h248",          // Megaco H-248 (Text) | Megaco H-248 | Megaco-H.248 text
	Port{2944, 6}:    "megaco-h248",          // Megaco H-248 (Text)
	Port{2944, 17}:   "megaco-h248",          // Megaco H-248 (Text)
	Port{2945, 132}:  "h248-binary",          // Megaco H-248 (Binary) | H248 Binary | Megaco H.248 binary
	Port{2945, 6}:    "h248-binary",          // Megaco H-248 (Binary)
	Port{2945, 17}:   "h248-binary",          // Megaco H-248 (Binary)
	Port{2946, 6}:    "fjsvmpor",             // Missing description for fjsvmpor
	Port{2946, 17}:   "fjsvmpor",             // FJSVmpor
	Port{2947, 6}:    "gpsd",                 // GPS Daemon request response protocol
	Port{2947, 17}:   "gpsd",                 // GPS Daemon request response protocol
	Port{2948, 6}:    "wap-push",             // WAP PUSH
	Port{2948, 17}:   "wap-push",             // Windows Mobile devices often have this
	Port{2949, 6}:    "wap-pushsecure",       // WAP PUSH SECURE
	Port{2949, 17}:   "wap-pushsecure",       // WAP PUSH SECURE
	Port{2950, 6}:    "esip",                 // Missing description for esip
	Port{2950, 17}:   "esip",                 // ESIP
	Port{2951, 6}:    "ottp",                 // Missing description for ottp
	Port{2951, 17}:   "ottp",                 // OTTP
	Port{2952, 6}:    "mpfwsas",              // Missing description for mpfwsas
	Port{2952, 17}:   "mpfwsas",              // MPFWSAS
	Port{2953, 6}:    "ovalarmsrv",           // Missing description for ovalarmsrv
	Port{2953, 17}:   "ovalarmsrv",           // OVALARMSRV
	Port{2954, 6}:    "ovalarmsrv-cmd",       // Missing description for ovalarmsrv-cmd
	Port{2954, 17}:   "ovalarmsrv-cmd",       // OVALARMSRV-CMD
	Port{2955, 6}:    "csnotify",             // Missing description for csnotify
	Port{2955, 17}:   "csnotify",             // CSNOTIFY
	Port{2956, 6}:    "ovrimosdbman",         // Missing description for ovrimosdbman
	Port{2956, 17}:   "ovrimosdbman",         // OVRIMOSDBMAN
	Port{2957, 6}:    "jmact5",               // JAMCT5
	Port{2957, 17}:   "jmact5",               // JAMCT5
	Port{2958, 6}:    "jmact6",               // JAMCT6
	Port{2958, 17}:   "jmact6",               // JAMCT6
	Port{2959, 6}:    "rmopagt",              // Missing description for rmopagt
	Port{2959, 17}:   "rmopagt",              // RMOPAGT
	Port{2960, 6}:    "dfoxserver",           // Missing description for dfoxserver
	Port{2960, 17}:   "dfoxserver",           // DFOXSERVER
	Port{2961, 6}:    "boldsoft-lm",          // Missing description for boldsoft-lm
	Port{2961, 17}:   "boldsoft-lm",          // BOLDSOFT-LM
	Port{2962, 6}:    "iph-policy-cli",       // Missing description for iph-policy-cli
	Port{2962, 17}:   "iph-policy-cli",       // IPH-POLICY-CLI
	Port{2963, 6}:    "iph-policy-adm",       // Missing description for iph-policy-adm
	Port{2963, 17}:   "iph-policy-adm",       // IPH-POLICY-ADM
	Port{2964, 6}:    "bullant-srap",         // BULLANT SRAP
	Port{2964, 17}:   "bullant-srap",         // BULLANT SRAP
	Port{2965, 6}:    "bullant-rap",          // BULLANT RAP
	Port{2965, 17}:   "bullant-rap",          // BULLANT RAP
	Port{2966, 6}:    "idp-infotrieve",       // Missing description for idp-infotrieve
	Port{2966, 17}:   "idp-infotrieve",       // IDP-INFOTRIEVE
	Port{2967, 6}:    "symantec-av",          // ssc-agent | Symantec AntiVirus (rtvscan.exe) | SSC-AGENT
	Port{2967, 17}:   "symantec-av",          // Symantec AntiVirus (rtvscan.exe)
	Port{2968, 6}:    "enpp",                 // Missing description for enpp
	Port{2968, 17}:   "enpp",                 // ENPP
	Port{2969, 6}:    "essp",                 // Missing description for essp
	Port{2969, 17}:   "essp",                 // ESSP
	Port{2970, 6}:    "index-net",            // Missing description for index-net
	Port{2970, 17}:   "index-net",            // INDEX-NET
	Port{2971, 6}:    "netclip",              // NetClip clipboard daemon
	Port{2971, 17}:   "netclip",              // NetClip clipboard daemon
	Port{2972, 6}:    "pmsm-webrctl",         // PMSM Webrctl
	Port{2972, 17}:   "pmsm-webrctl",         // PMSM Webrctl
	Port{2973, 6}:    "svnetworks",           // SV Networks
	Port{2973, 17}:   "svnetworks",           // SV Networks
	Port{2974, 6}:    "signal",               // Missing description for signal
	Port{2974, 17}:   "signal",               // Signal
	Port{2975, 6}:    "fjmpcm",               // Fujitsu Configuration Management Service
	Port{2975, 17}:   "fjmpcm",               // Fujitsu Configuration Management Service
	Port{2976, 6}:    "cns-srv-port",         // CNS Server Port
	Port{2976, 17}:   "cns-srv-port",         // CNS Server Port
	Port{2977, 6}:    "ttc-etap-ns",          // TTCs Enterprise Test Access Protocol - NS
	Port{2977, 17}:   "ttc-etap-ns",          // TTCs Enterprise Test Access Protocol - NS
	Port{2978, 6}:    "ttc-etap-ds",          // TTCs Enterprise Test Access Protocol - DS
	Port{2978, 17}:   "ttc-etap-ds",          // TTCs Enterprise Test Access Protocol - DS
	Port{2979, 6}:    "h263-video",           // H.263 Video Streaming
	Port{2979, 17}:   "h263-video",           // H.263 Video Streaming
	Port{2980, 6}:    "wimd",                 // Instant Messaging Service
	Port{2980, 17}:   "wimd",                 // Instant Messaging Service
	Port{2981, 6}:    "mylxamport",           // Missing description for mylxamport
	Port{2981, 17}:   "mylxamport",           // MYLXAMPORT
	Port{2982, 6}:    "iwb-whiteboard",       // Missing description for iwb-whiteboard
	Port{2982, 17}:   "iwb-whiteboard",       // IWB-WHITEBOARD
	Port{2983, 6}:    "netplan",              // Missing description for netplan
	Port{2983, 17}:   "netplan",              // NETPLAN
	Port{2984, 6}:    "hpidsadmin",           // Missing description for hpidsadmin
	Port{2984, 17}:   "hpidsadmin",           // HPIDSADMIN
	Port{2985, 6}:    "hpidsagent",           // Missing description for hpidsagent
	Port{2985, 17}:   "hpidsagent",           // HPIDSAGENT
	Port{2986, 6}:    "stonefalls",           // Missing description for stonefalls
	Port{2986, 17}:   "stonefalls",           // STONEFALLS
	Port{2987, 6}:    "identify",             // Missing description for identify
	Port{2987, 17}:   "identify",             // Missing description for identify
	Port{2988, 6}:    "hippad",               // HIPPA Reporting Protocol
	Port{2988, 17}:   "hippad",               // HIPPA Reporting Protocol
	Port{2989, 6}:    "zarkov",               // ZARKOV Intelligent Agent Communication
	Port{2989, 17}:   "zarkov",               // ZARKOV Intelligent Agent Communication
	Port{2990, 6}:    "boscap",               // Missing description for boscap
	Port{2990, 17}:   "boscap",               // BOSCAP
	Port{2991, 6}:    "wkstn-mon",            // Missing description for wkstn-mon
	Port{2991, 17}:   "wkstn-mon",            // WKSTN-MON
	Port{2992, 6}:    "avenyo",               // Avenyo Server
	Port{2992, 17}:   "avenyo",               // Avenyo Server
	Port{2993, 6}:    "veritas-vis1",         // VERITAS VIS1
	Port{2993, 17}:   "veritas-vis1",         // VERITAS VIS1
	Port{2994, 6}:    "veritas-vis2",         // VERITAS VIS2
	Port{2994, 17}:   "veritas-vis2",         // VERITAS VIS2
	Port{2995, 6}:    "idrs",                 // Missing description for idrs
	Port{2995, 17}:   "idrs",                 // IDRS
	Port{2996, 6}:    "vsixml",               // Missing description for vsixml
	Port{2996, 17}:   "vsixml",               // Missing description for vsixml
	Port{2997, 6}:    "rebol",                // Missing description for rebol
	Port{2997, 17}:   "rebol",                // REBOL
	Port{2998, 6}:    "iss-realsec",          // realsecure | ISS RealSecure IDS Remote Console Admin port | Real Secure
	Port{2998, 17}:   "realsecure",           // Real Secure
	Port{2999, 6}:    "remoteware-un",        // RemoteWare Unassigned
	Port{2999, 17}:   "remoteware-un",        // RemoteWare Unassigned
	Port{3000, 6}:    "ppp",                  // remoteware-cl | hbci | User-level ppp daemon, or chili!soft asp | HBCI | RemoteWare Client
	Port{3000, 17}:   "hbci",                 // HBCI
	Port{3001, 6}:    "nessus",               // origo-native | Nessus Security Scanner (www.nessus.org) Daemon or chili!soft asp | OrigoDB Server Native Interface
	Port{3002, 6}:    "exlm-agent",           // remoteware-srv | EXLM Agent | RemoteWare Server
	Port{3002, 17}:   "exlm-agent",           // EXLM Agent
	Port{3003, 6}:    "cgms",                 // Missing description for cgms
	Port{3003, 17}:   "cgms",                 // CGMS
	Port{3004, 6}:    "csoftragent",          // Csoft Agent
	Port{3004, 17}:   "csoftragent",          // Csoft Agent
	Port{3005, 6}:    "deslogin",             // geniuslm | encrypted symmetric telnet login | Genius License Manager
	Port{3005, 17}:   "geniuslm",             // Genius License Manager
	Port{3006, 6}:    "deslogind",            // ii-admin | Instant Internet Admin
	Port{3006, 17}:   "ii-admin",             // Instant Internet Admin
	Port{3007, 6}:    "lotusmtap",            // Lotus Mail Tracking Agent Protocol
	Port{3007, 17}:   "lotusmtap",            // Lotus Mail Tracking Agent Protocol
	Port{3008, 6}:    "midnight-tech",        // Midnight Technologies
	Port{3008, 17}:   "midnight-tech",        // Midnight Technologies
	Port{3009, 6}:    "pxc-ntfy",             // Missing description for pxc-ntfy
	Port{3009, 17}:   "pxc-ntfy",             // PXC-NTFY
	Port{3010, 6}:    "gw",                   // ping-pong | Telerate Workstation
	Port{3010, 17}:   "ping-pong",            // Telerate Workstation
	Port{3011, 6}:    "trusted-web",          // Trusted Web
	Port{3011, 17}:   "trusted-web",          // Trusted Web
	Port{3012, 6}:    "twsdss",               // Trusted Web Client
	Port{3012, 17}:   "twsdss",               // Trusted Web Client
	Port{3013, 6}:    "gilatskysurfer",       // Gilat Sky Surfer
	Port{3013, 17}:   "gilatskysurfer",       // Gilat Sky Surfer
	Port{3014, 6}:    "broker_service",       // broker-service | Broker Service
	Port{3014, 17}:   "broker_service",       // Broker Service
	Port{3015, 6}:    "nati-dstp",            // NATI DSTP
	Port{3015, 17}:   "nati-dstp",            // NATI DSTP
	Port{3016, 6}:    "notify_srvr",          // notify-srvr | Notify Server
	Port{3016, 17}:   "notify_srvr",          // Notify Server
	Port{3017, 6}:    "event_listener",       // event-listener | Event Listener
	Port{3017, 17}:   "event_listener",       // Event Listener
	Port{3018, 6}:    "srvc_registry",        // srvc-registry | Service Registry
	Port{3018, 17}:   "srvc_registry",        // Service Registry
	Port{3019, 6}:    "resource_mgr",         // resource-mgr | Resource Manager
	Port{3019, 17}:   "resource_mgr",         // Resource Manager
	Port{3020, 6}:    "cifs",                 // Missing description for cifs
	Port{3020, 17}:   "cifs",                 // CIFS
	Port{3021, 6}:    "agriserver",           // AGRI Server
	Port{3021, 17}:   "agriserver",           // AGRI Server
	Port{3022, 6}:    "csregagent",           // Missing description for csregagent
	Port{3022, 17}:   "csregagent",           // CSREGAGENT
	Port{3023, 6}:    "magicnotes",           // Missing description for magicnotes
	Port{3023, 17}:   "magicnotes",           // Missing description for magicnotes
	Port{3024, 6}:    "nds_sso",              // nds-sso
	Port{3024, 17}:   "nds_sso",              // NDS_SSO
	Port{3025, 6}:    "slnp",                 // arepa-raft | SLNP (Simple Library Network Protocol) by Sisis Informationssysteme GmbH | Arepa Raft
	Port{3025, 17}:   "arepa-raft",           // Arepa Raft
	Port{3026, 6}:    "agri-gateway",         // AGRI Gateway
	Port{3026, 17}:   "agri-gateway",         // AGRI Gateway
	Port{3027, 6}:    "LiebDevMgmt_C",        // LiebDevMgmt-C
	Port{3027, 17}:   "LiebDevMgmt_C",        // Missing description for LiebDevMgmt_C
	Port{3028, 6}:    "LiebDevMgmt_DM",       // LiebDevMgmt-DM
	Port{3028, 17}:   "LiebDevMgmt_DM",       // Missing description for LiebDevMgmt_DM
	Port{3029, 6}:    "LiebDevMgmt_A",        // LiebDevMgmt-A
	Port{3029, 17}:   "LiebDevMgmt_A",        // Missing description for LiebDevMgmt_A
	Port{3030, 6}:    "arepa-cas",            // Arepa Cas
	Port{3030, 17}:   "arepa-cas",            // Arepa Cas
	Port{3031, 6}:    "eppc",                 // Remote AppleEvents PPC Toolbox
	Port{3031, 17}:   "eppc",                 // Remote AppleEvents PPC Toolbox
	Port{3032, 6}:    "redwood-chat",         // Redwood Chat
	Port{3032, 17}:   "redwood-chat",         // Redwood Chat
	Port{3033, 6}:    "pdb",                  // Missing description for pdb
	Port{3033, 17}:   "pdb",                  // PDB
	Port{3034, 6}:    "osmosis-aeea",         // Osmosis   Helix (R) AEEA Port
	Port{3034, 17}:   "osmosis-aeea",         // Osmosis   Helix (R) AEEA Port
	Port{3035, 6}:    "fjsv-gssagt",          // FJSV gssagt
	Port{3035, 17}:   "fjsv-gssagt",          // FJSV gssagt
	Port{3036, 6}:    "hagel-dump",           // Hagel DUMP
	Port{3036, 17}:   "hagel-dump",           // Hagel DUMP
	Port{3037, 6}:    "hp-san-mgmt",          // HP SAN Mgmt
	Port{3037, 17}:   "hp-san-mgmt",          // HP SAN Mgmt
	Port{3038, 6}:    "santak-ups",           // Santak UPS
	Port{3038, 17}:   "santak-ups",           // Santak UPS
	Port{3039, 6}:    "cogitate",             // Cogitate, Inc.
	Port{3039, 17}:   "cogitate",             // Cogitate, Inc.
	Port{3040, 6}:    "tomato-springs",       // Tomato Springs
	Port{3040, 17}:   "tomato-springs",       // Tomato Springs
	Port{3041, 6}:    "di-traceware",         // Missing description for di-traceware
	Port{3041, 17}:   "di-traceware",         // Missing description for di-traceware
	Port{3042, 6}:    "journee",              // Missing description for journee
	Port{3042, 17}:   "journee",              // Missing description for journee
	Port{3043, 6}:    "brp",                  // Broadcast Routing Protocol
	Port{3043, 17}:   "brp",                  // Broadcast Routing Protocol
	Port{3044, 6}:    "epp",                  // EndPoint Protocol
	Port{3044, 17}:   "epp",                  // EndPoint Protocol
	Port{3045, 6}:    "slnp",                 // responsenet | SLNP (Simple Library Network Protocol) by Sisis Informationssysteme GmbH | ResponseNet
	Port{3045, 17}:   "responsenet",          // ResponseNet
	Port{3046, 6}:    "di-ase",               // Missing description for di-ase
	Port{3046, 17}:   "di-ase",               // Missing description for di-ase
	Port{3047, 6}:    "hlserver",             // Fast Security HL Server
	Port{3047, 17}:   "hlserver",             // Fast Security HL Server
	Port{3048, 6}:    "pctrader",             // Sierra Net PC Trader
	Port{3048, 17}:   "pctrader",             // Sierra Net PC Trader
	Port{3049, 6}:    "cfs",                  // nsws | cryptographic file system (nfs) (proposed) | NSWS
	Port{3049, 17}:   "cfs",                  // cryptographic file system (nfs)
	Port{3050, 6}:    "gds_db",               // gds-db
	Port{3050, 17}:   "gds_db",               // Missing description for gds_db
	Port{3051, 6}:    "galaxy-server",        // Galaxy Server
	Port{3051, 17}:   "galaxy-server",        // Galaxy Server
	Port{3052, 6}:    "powerchute",           // apc-3052 | APC 3052
	Port{3052, 17}:   "apc-3052",             // APC 3052
	Port{3053, 6}:    "dsom-server",          // Missing description for dsom-server
	Port{3053, 17}:   "dsom-server",          // Missing description for dsom-server
	Port{3054, 6}:    "amt-cnf-prot",         // AMT CNF PROT
	Port{3054, 17}:   "amt-cnf-prot",         // AMT CNF PROT
	Port{3055, 6}:    "policyserver",         // Policy Server
	Port{3055, 17}:   "policyserver",         // Policy Server
	Port{3056, 6}:    "cdl-server",           // CDL Server
	Port{3056, 17}:   "cdl-server",           // CDL Server
	Port{3057, 6}:    "goahead-fldup",        // GoAhead FldUp
	Port{3057, 17}:   "goahead-fldup",        // GoAhead FldUp
	Port{3058, 6}:    "videobeans",           // Missing description for videobeans
	Port{3058, 17}:   "videobeans",           // Missing description for videobeans
	Port{3059, 6}:    "qsoft",                // Missing description for qsoft
	Port{3059, 17}:   "qsoft",                // Missing description for qsoft
	Port{3060, 6}:    "interserver",          // Missing description for interserver
	Port{3060, 17}:   "interserver",          // Missing description for interserver
	Port{3061, 6}:    "cautcpd",              // Missing description for cautcpd
	Port{3061, 17}:   "cautcpd",              // Missing description for cautcpd
	Port{3062, 6}:    "ncacn-ip-tcp",         // Missing description for ncacn-ip-tcp
	Port{3062, 17}:   "ncacn-ip-tcp",         // Missing description for ncacn-ip-tcp
	Port{3063, 6}:    "ncadg-ip-udp",         // Missing description for ncadg-ip-udp
	Port{3063, 17}:   "ncadg-ip-udp",         // Missing description for ncadg-ip-udp
	Port{3064, 6}:    "dnet-tstproxy",        // rprt | distributed.net (a closed source crypto-cracking project) proxy test port | Remote Port Redirector
	Port{3064, 17}:   "rprt",                 // Remote Port Redirector
	Port{3065, 6}:    "slinterbase",          // Missing description for slinterbase
	Port{3065, 17}:   "slinterbase",          // Missing description for slinterbase
	Port{3066, 6}:    "netattachsdmp",        // Missing description for netattachsdmp
	Port{3066, 17}:   "netattachsdmp",        // NETATTACHSDMP
	Port{3067, 6}:    "fjhpjp",               // Missing description for fjhpjp
	Port{3067, 17}:   "fjhpjp",               // FJHPJP
	Port{3068, 6}:    "ls3bcast",             // ls3 Broadcast
	Port{3068, 17}:   "ls3bcast",             // ls3 Broadcast
	Port{3069, 6}:    "ls3",                  // Missing description for ls3
	Port{3069, 17}:   "ls3",                  // Missing description for ls3
	Port{3070, 6}:    "mgxswitch",            // Missing description for mgxswitch
	Port{3070, 17}:   "mgxswitch",            // MGXSWITCH
	Port{3071, 6}:    "csd-mgmt-port",        // xplat-replicate | ContinuStor Manager Port | Crossplatform replication protocol
	Port{3071, 17}:   "csd-mgmt-port",        // ContinuStor Manager Port
	Port{3072, 6}:    "csd-monitor",          // ContinuStor Monitor Port
	Port{3072, 17}:   "csd-monitor",          // ContinuStor Monitor Port
	Port{3073, 6}:    "vcrp",                 // Very simple chatroom prot
	Port{3073, 17}:   "vcrp",                 // Very simple chatroom prot
	Port{3074, 6}:    "xbox",                 // Xbox game port
	Port{3074, 17}:   "xbox",                 // Xbox game port
	Port{3075, 6}:    "orbix-locator",        // Orbix 2000 Locator
	Port{3075, 17}:   "orbix-locator",        // Orbix 2000 Locator
	Port{3076, 6}:    "orbix-config",         // Orbix 2000 Config
	Port{3076, 17}:   "orbix-config",         // Orbix 2000 Config
	Port{3077, 6}:    "orbix-loc-ssl",        // Orbix 2000 Locator SSL
	Port{3077, 17}:   "orbix-loc-ssl",        // Orbix 2000 Locator SSL
	Port{3078, 6}:    "orbix-cfg-ssl",        // Orbix 2000 Locator SSL
	Port{3078, 17}:   "orbix-cfg-ssl",        // Orbix 2000 Locator SSL
	Port{3079, 6}:    "lv-frontpanel",        // LV Front Panel
	Port{3079, 17}:   "lv-frontpanel",        // LV Front Panel
	Port{3080, 6}:    "stm_pproc",            // stm-pproc
	Port{3080, 17}:   "stm_pproc",            // Missing description for stm_pproc
	Port{3081, 6}:    "tl1-lv",               // Missing description for tl1-lv
	Port{3081, 17}:   "tl1-lv",               // TL1-LV
	Port{3082, 6}:    "tl1-raw",              // Missing description for tl1-raw
	Port{3082, 17}:   "tl1-raw",              // TL1-RAW
	Port{3083, 6}:    "tl1-telnet",           // Missing description for tl1-telnet
	Port{3083, 17}:   "tl1-telnet",           // TL1-TELNET
	Port{3084, 6}:    "itm-mccs",             // Missing description for itm-mccs
	Port{3084, 17}:   "itm-mccs",             // ITM-MCCS
	Port{3085, 6}:    "pcihreq",              // Missing description for pcihreq
	Port{3085, 17}:   "pcihreq",              // PCIHReq
	Port{3086, 6}:    "sj3",                  // jdl-dbkitchen | SJ3 (kanji input) | JDL-DBKitchen
	Port{3086, 17}:   "jdl-dbkitchen",        // JDL-DBKitchen
	Port{3087, 6}:    "asoki-sma",            // Asoki SMA
	Port{3087, 17}:   "asoki-sma",            // Asoki SMA
	Port{3088, 6}:    "xdtp",                 // eXtensible Data Transfer Protocol
	Port{3088, 17}:   "xdtp",                 // eXtensible Data Transfer Protocol
	Port{3089, 6}:    "ptk-alink",            // ParaTek Agent Linking
	Port{3089, 17}:   "ptk-alink",            // ParaTek Agent Linking
	Port{3090, 6}:    "stss",                 // Senforce Session Services
	Port{3090, 17}:   "stss",                 // Senforce Session Services
	Port{3091, 6}:    "1ci-smcs",             // 1Ci Server Management
	Port{3091, 17}:   "1ci-smcs",             // 1Ci Server Management
	Port{3093, 6}:    "rapidmq-center",       // Jiiva RapidMQ Center
	Port{3093, 17}:   "rapidmq-center",       // Jiiva RapidMQ Center
	Port{3094, 6}:    "rapidmq-reg",          // Jiiva RapidMQ Registry
	Port{3094, 17}:   "rapidmq-reg",          // Jiiva RapidMQ Registry
	Port{3095, 6}:    "panasas",              // Panasas rendevous port | Panasas rendezvous port
	Port{3095, 17}:   "panasas",              // Panasas rendevous port
	Port{3096, 6}:    "ndl-aps",              // Active Print Server Port
	Port{3096, 17}:   "ndl-aps",              // Active Print Server Port
	Port{3097, 132}:  "itu-bicc-stc",         // ITU-T Q.1902.1 Q.2150.3
	Port{3098, 6}:    "umm-port",             // Universal Message Manager
	Port{3098, 17}:   "umm-port",             // Universal Message Manager
	Port{3099, 6}:    "chmd",                 // CHIPSY Machine Daemon
	Port{3099, 17}:   "chmd",                 // CHIPSY Machine Daemon
	Port{3100, 6}:    "opcon-xps",            // OpCon xps
	Port{3100, 17}:   "opcon-xps",            // OpCon xps
	Port{3101, 6}:    "hp-pxpib",             // HP PolicyXpert PIB Server
	Port{3101, 17}:   "hp-pxpib",             // HP PolicyXpert PIB Server
	Port{3102, 6}:    "slslavemon",           // SoftlinK Slave Mon Port
	Port{3102, 17}:   "slslavemon",           // SoftlinK Slave Mon Port
	Port{3103, 6}:    "autocuesmi",           // Autocue SMI Protocol
	Port{3103, 17}:   "autocuesmi",           // Autocue SMI Protocol
	Port{3104, 6}:    "autocuelog",           // autocuetime | Autocue Logger Protocol | Autocue Time Service
	Port{3104, 17}:   "autocuetime",          // Autocue Time Service
	Port{3105, 6}:    "cardbox",              // Missing description for cardbox
	Port{3105, 17}:   "cardbox",              // Cardbox
	Port{3106, 6}:    "cardbox-http",         // Cardbox HTTP
	Port{3106, 17}:   "cardbox-http",         // Cardbox HTTP
	Port{3107, 6}:    "business",             // Business protocol
	Port{3107, 17}:   "business",             // Business protocol
	Port{3108, 6}:    "geolocate",            // Geolocate protocol
	Port{3108, 17}:   "geolocate",            // Geolocate protocol
	Port{3109, 6}:    "personnel",            // Personnel protocol
	Port{3109, 17}:   "personnel",            // Personnel protocol
	Port{3110, 6}:    "sim-control",          // simulator control port
	Port{3110, 17}:   "sim-control",          // simulator control port
	Port{3111, 6}:    "wsynch",               // Web Synchronous Services
	Port{3111, 17}:   "wsynch",               // Web Synchronous Services
	Port{3112, 6}:    "ksysguard",            // KDE System Guard
	Port{3112, 17}:   "ksysguard",            // KDE System Guard
	Port{3113, 6}:    "cs-auth-svr",          // CS-Authenticate Svr Port
	Port{3113, 17}:   "cs-auth-svr",          // CS-Authenticate Svr Port
	Port{3114, 6}:    "ccmad",                // CCM AutoDiscover
	Port{3114, 17}:   "ccmad",                // CCM AutoDiscover
	Port{3115, 6}:    "mctet-master",         // MCTET Master
	Port{3115, 17}:   "mctet-master",         // MCTET Master
	Port{3116, 6}:    "mctet-gateway",        // MCTET Gateway
	Port{3116, 17}:   "mctet-gateway",        // MCTET Gateway
	Port{3117, 6}:    "mctet-jserv",          // MCTET Jserv
	Port{3117, 17}:   "mctet-jserv",          // MCTET Jserv
	Port{3118, 6}:    "pkagent",              // Missing description for pkagent
	Port{3118, 17}:   "pkagent",              // PKAgent
	Port{3119, 6}:    "d2000kernel",          // D2000 Kernel Port
	Port{3119, 17}:   "d2000kernel",          // D2000 Kernel Port
	Port{3120, 6}:    "d2000webserver",       // D2000 Webserver Port
	Port{3120, 17}:   "d2000webserver",       // D2000 Webserver Port
	Port{3121, 6}:    "pcmk-remote",          // The pacemaker remote (pcmk-remote) service extends high availability functionality outside of the Linux cluster into remote nodes.
	Port{3122, 6}:    "vtr-emulator",         // MTI VTR Emulator port
	Port{3122, 17}:   "vtr-emulator",         // MTI VTR Emulator port
	Port{3123, 6}:    "edix",                 // EDI Translation Protocol
	Port{3123, 17}:   "edix",                 // EDI Translation Protocol
	Port{3124, 6}:    "beacon-port",          // Beacon Port
	Port{3124, 17}:   "beacon-port",          // Beacon Port
	Port{3125, 6}:    "a13-an",               // A13-AN Interface
	Port{3125, 17}:   "a13-an",               // A13-AN Interface
	Port{3127, 6}:    "ctx-bridge",           // CTX Bridge Port
	Port{3127, 17}:   "ctx-bridge",           // CTX Bridge Port
	Port{3128, 6}:    "squid-http",           // ndl-aas | Active API Server Port
	Port{3128, 17}:   "ndl-aas",              // Active API Server Port
	Port{3129, 6}:    "netport-id",           // NetPort Discovery Port
	Port{3129, 17}:   "netport-id",           // NetPort Discovery Port
	Port{3130, 6}:    "icpv2",                // Missing description for icpv2
	Port{3130, 17}:   "squid-ipc",            // Missing description for squid-ipc
	Port{3131, 6}:    "netbookmark",          // Net Book Mark
	Port{3131, 17}:   "netbookmark",          // Net Book Mark
	Port{3132, 6}:    "ms-rule-engine",       // Microsoft Business Rule Engine Update Service
	Port{3132, 17}:   "ms-rule-engine",       // Microsoft Business Rule Engine Update Service
	Port{3133, 6}:    "prism-deploy",         // Prism Deploy User Port
	Port{3133, 17}:   "prism-deploy",         // Prism Deploy User Port
	Port{3134, 6}:    "ecp",                  // Extensible Code Protocol
	Port{3134, 17}:   "ecp",                  // Extensible Code Protocol
	Port{3135, 6}:    "peerbook-port",        // PeerBook Port
	Port{3135, 17}:   "peerbook-port",        // PeerBook Port
	Port{3136, 6}:    "grubd",                // Grub Server Port
	Port{3136, 17}:   "grubd",                // Grub Server Port
	Port{3137, 6}:    "rtnt-1",               // rtnt-1 data packets
	Port{3137, 17}:   "rtnt-1",               // rtnt-1 data packets
	Port{3138, 6}:    "rtnt-2",               // rtnt-2 data packets
	Port{3138, 17}:   "rtnt-2",               // rtnt-2 data packets
	Port{3139, 6}:    "incognitorv",          // Incognito Rendez-Vous
	Port{3139, 17}:   "incognitorv",          // Incognito Rendez-Vous
	Port{3140, 6}:    "ariliamulti",          // Arilia Multiplexor
	Port{3140, 17}:   "ariliamulti",          // Arilia Multiplexor
	Port{3141, 6}:    "vmodem",               // Missing description for vmodem
	Port{3141, 17}:   "vmodem",               // Missing description for vmodem
	Port{3142, 6}:    "apt-cacher",           // rdc-wh-eos | A server which keeps a local cache of Debian Ubuntu package files | RDC WH EOS
	Port{3142, 17}:   "rdc-wh-eos",           // RDC WH EOS
	Port{3143, 6}:    "seaview",              // Sea View
	Port{3143, 17}:   "seaview",              // Sea View
	Port{3144, 6}:    "tarantella",           // Missing description for tarantella
	Port{3144, 17}:   "tarantella",           // Tarantella
	Port{3145, 6}:    "csi-lfap",             // Missing description for csi-lfap
	Port{3145, 17}:   "csi-lfap",             // CSI-LFAP
	Port{3146, 6}:    "bears-02",             // Missing description for bears-02
	Port{3146, 17}:   "bears-02",             // Missing description for bears-02
	Port{3147, 6}:    "rfio",                 // Missing description for rfio
	Port{3147, 17}:   "rfio",                 // RFIO
	Port{3148, 6}:    "nm-game-admin",        // NetMike Game Administrator
	Port{3148, 17}:   "nm-game-admin",        // NetMike Game Administrator
	Port{3149, 6}:    "nm-game-server",       // NetMike Game Server
	Port{3149, 17}:   "nm-game-server",       // NetMike Game Server
	Port{3150, 6}:    "nm-asses-admin",       // NetMike Assessor Administrator
	Port{3150, 17}:   "nm-asses-admin",       // NetMike Assessor Administrator
	Port{3151, 6}:    "nm-assessor",          // NetMike Assessor
	Port{3151, 17}:   "nm-assessor",          // NetMike Assessor
	Port{3152, 6}:    "feitianrockey",        // FeiTian Port
	Port{3152, 17}:   "feitianrockey",        // FeiTian Port
	Port{3153, 6}:    "s8-client-port",       // S8Cargo Client Port
	Port{3153, 17}:   "s8-client-port",       // S8Cargo Client Port
	Port{3154, 6}:    "ccmrmi",               // ON RMI Registry
	Port{3154, 17}:   "ccmrmi",               // ON RMI Registry
	Port{3155, 6}:    "jpegmpeg",             // JpegMpeg Port
	Port{3155, 17}:   "jpegmpeg",             // JpegMpeg Port
	Port{3156, 6}:    "indura",               // Indura Collector
	Port{3156, 17}:   "indura",               // Indura Collector
	Port{3157, 6}:    "e3consultants",        // CCC Listener Port
	Port{3157, 17}:   "e3consultants",        // CCC Listener Port
	Port{3158, 6}:    "stvp",                 // SmashTV Protocol
	Port{3158, 17}:   "stvp",                 // SmashTV Protocol
	Port{3159, 6}:    "navegaweb-port",       // NavegaWeb Tarification
	Port{3159, 17}:   "navegaweb-port",       // NavegaWeb Tarification
	Port{3160, 6}:    "tip-app-server",       // TIP Application Server
	Port{3160, 17}:   "tip-app-server",       // TIP Application Server
	Port{3161, 6}:    "doc1lm",               // DOC1 License Manager
	Port{3161, 17}:   "doc1lm",               // DOC1 License Manager
	Port{3162, 6}:    "sflm",                 // Missing description for sflm
	Port{3162, 17}:   "sflm",                 // SFLM
	Port{3163, 6}:    "res-sap",              // Missing description for res-sap
	Port{3163, 17}:   "res-sap",              // RES-SAP
	Port{3164, 6}:    "imprs",                // Missing description for imprs
	Port{3164, 17}:   "imprs",                // IMPRS
	Port{3165, 6}:    "newgenpay",            // Newgenpay Engine Service
	Port{3165, 17}:   "newgenpay",            // Newgenpay Engine Service
	Port{3166, 6}:    "sossecollector",       // Quest Spotlight Out-Of-Process Collector
	Port{3166, 17}:   "sossecollector",       // Quest Spotlight Out-Of-Process Collector
	Port{3167, 6}:    "nowcontact",           // Now Contact Public Server
	Port{3167, 17}:   "nowcontact",           // Now Contact Public Server
	Port{3168, 6}:    "poweronnud",           // Now Up-to-Date Public Server
	Port{3168, 17}:   "poweronnud",           // Now Up-to-Date Public Server
	Port{3169, 6}:    "serverview-as",        // Missing description for serverview-as
	Port{3169, 17}:   "serverview-as",        // SERVERVIEW-AS
	Port{3170, 6}:    "serverview-asn",       // Missing description for serverview-asn
	Port{3170, 17}:   "serverview-asn",       // SERVERVIEW-ASN
	Port{3171, 6}:    "serverview-gf",        // Missing description for serverview-gf
	Port{3171, 17}:   "serverview-gf",        // SERVERVIEW-GF
	Port{3172, 6}:    "serverview-rm",        // Missing description for serverview-rm
	Port{3172, 17}:   "serverview-rm",        // SERVERVIEW-RM
	Port{3173, 6}:    "serverview-icc",       // Missing description for serverview-icc
	Port{3173, 17}:   "serverview-icc",       // SERVERVIEW-ICC
	Port{3174, 6}:    "armi-server",          // ARMI Server
	Port{3174, 17}:   "armi-server",          // ARMI Server
	Port{3175, 6}:    "t1-e1-over-ip",        // T1_E1_Over_IP
	Port{3175, 17}:   "t1-e1-over-ip",        // T1_E1_Over_IP
	Port{3176, 6}:    "ars-master",           // ARS Master
	Port{3176, 17}:   "ars-master",           // ARS Master
	Port{3177, 6}:    "phonex-port",          // Phonex Protocol
	Port{3177, 17}:   "phonex-port",          // Phonex Protocol
	Port{3178, 6}:    "radclientport",        // Radiance UltraEdge Port
	Port{3178, 17}:   "radclientport",        // Radiance UltraEdge Port
	Port{3179, 6}:    "h2gf-w-2m",            // H2GF W.2m Handover prot.
	Port{3179, 17}:   "h2gf-w-2m",            // H2GF W.2m Handover prot.
	Port{3180, 6}:    "mc-brk-srv",           // Millicent Broker Server
	Port{3180, 17}:   "mc-brk-srv",           // Millicent Broker Server
	Port{3181, 6}:    "bmcpatrolagent",       // BMC Patrol Agent
	Port{3181, 17}:   "bmcpatrolagent",       // BMC Patrol Agent
	Port{3182, 6}:    "bmcpatrolrnvu",        // BMC Patrol Rendezvous
	Port{3182, 17}:   "bmcpatrolrnvu",        // BMC Patrol Rendezvous
	Port{3183, 6}:    "cops-tls",             // COPS TLS
	Port{3183, 17}:   "cops-tls",             // COPS TLS
	Port{3184, 6}:    "apogeex-port",         // ApogeeX Port
	Port{3184, 17}:   "apogeex-port",         // ApogeeX Port
	Port{3185, 6}:    "smpppd",               // SuSE Meta PPPD
	Port{3185, 17}:   "smpppd",               // SuSE Meta PPPD
	Port{3186, 6}:    "iiw-port",             // IIW Monitor User Port
	Port{3186, 17}:   "iiw-port",             // IIW Monitor User Port
	Port{3187, 6}:    "odi-port",             // Open Design Listen Port
	Port{3187, 17}:   "odi-port",             // Open Design Listen Port
	Port{3188, 6}:    "brcm-comm-port",       // Broadcom Port
	Port{3188, 17}:   "brcm-comm-port",       // Broadcom Port
	Port{3189, 6}:    "pcle-infex",           // Pinnacle Sys InfEx Port
	Port{3189, 17}:   "pcle-infex",           // Pinnacle Sys InfEx Port
	Port{3190, 6}:    "csvr-proxy",           // ConServR Proxy
	Port{3190, 17}:   "csvr-proxy",           // ConServR Proxy
	Port{3191, 6}:    "csvr-sslproxy",        // ConServR SSL Proxy
	Port{3191, 17}:   "csvr-sslproxy",        // ConServR SSL Proxy
	Port{3192, 6}:    "firemonrcc",           // FireMon Revision Control
	Port{3192, 17}:   "firemonrcc",           // FireMon Revision Control
	Port{3193, 6}:    "spandataport",         // Missing description for spandataport
	Port{3193, 17}:   "spandataport",         // SpanDataPort
	Port{3194, 6}:    "magbind",              // Rockstorm MAG protocol
	Port{3194, 17}:   "magbind",              // Rockstorm MAG protocol
	Port{3195, 6}:    "ncu-1",                // Network Control Unit
	Port{3195, 17}:   "ncu-1",                // Network Control Unit
	Port{3196, 6}:    "ncu-2",                // Network Control Unit
	Port{3196, 17}:   "ncu-2",                // Network Control Unit
	Port{3197, 6}:    "embrace-dp-s",         // Embrace Device Protocol Server
	Port{3197, 17}:   "embrace-dp-s",         // Embrace Device Protocol Server
	Port{3198, 6}:    "embrace-dp-c",         // Embrace Device Protocol Client
	Port{3198, 17}:   "embrace-dp-c",         // Embrace Device Protocol Client
	Port{3199, 6}:    "dmod-workspace",       // DMOD WorkSpace
	Port{3199, 17}:   "dmod-workspace",       // DMOD WorkSpace
	Port{3200, 6}:    "tick-port",            // Press-sense Tick Port
	Port{3200, 17}:   "tick-port",            // Press-sense Tick Port
	Port{3201, 6}:    "cpq-tasksmart",        // Missing description for cpq-tasksmart
	Port{3201, 17}:   "cpq-tasksmart",        // CPQ-TaskSmart
	Port{3202, 6}:    "intraintra",           // Missing description for intraintra
	Port{3202, 17}:   "intraintra",           // IntraIntra
	Port{3203, 6}:    "netwatcher-mon",       // Network Watcher Monitor
	Port{3203, 17}:   "netwatcher-mon",       // Network Watcher Monitor
	Port{3204, 6}:    "netwatcher-db",        // Network Watcher DB Access
	Port{3204, 17}:   "netwatcher-db",        // Network Watcher DB Access
	Port{3205, 6}:    "isns",                 // iSNS Server Port
	Port{3205, 17}:   "isns",                 // iSNS Server Port
	Port{3206, 6}:    "ironmail",             // IronMail POP Proxy
	Port{3206, 17}:   "ironmail",             // IronMail POP Proxy
	Port{3207, 6}:    "vx-auth-port",         // Veritas Authentication Port
	Port{3207, 17}:   "vx-auth-port",         // Veritas Authentication Port
	Port{3208, 6}:    "pfu-prcallback",       // PFU PR Callback
	Port{3208, 17}:   "pfu-prcallback",       // PFU PR Callback
	Port{3209, 6}:    "netwkpathengine",      // HP OpenView Network Path Engine Server
	Port{3209, 17}:   "netwkpathengine",      // HP OpenView Network Path Engine Server
	Port{3210, 6}:    "flamenco-proxy",       // Flamenco Networks Proxy
	Port{3210, 17}:   "flamenco-proxy",       // Flamenco Networks Proxy
	Port{3211, 6}:    "avsecuremgmt",         // Avocent Secure Management
	Port{3211, 17}:   "avsecuremgmt",         // Avocent Secure Management
	Port{3212, 6}:    "surveyinst",           // Survey Instrument
	Port{3212, 17}:   "surveyinst",           // Survey Instrument
	Port{3213, 6}:    "neon24x7",             // NEON 24X7 Mission Control
	Port{3213, 17}:   "neon24x7",             // NEON 24X7 Mission Control
	Port{3214, 6}:    "jmq-daemon-1",         // JMQ Daemon Port 1
	Port{3214, 17}:   "jmq-daemon-1",         // JMQ Daemon Port 1
	Port{3215, 6}:    "jmq-daemon-2",         // JMQ Daemon Port 2
	Port{3215, 17}:   "jmq-daemon-2",         // JMQ Daemon Port 2
	Port{3216, 6}:    "ferrari-foam",         // Ferrari electronic FOAM
	Port{3216, 17}:   "ferrari-foam",         // Ferrari electronic FOAM
	Port{3217, 6}:    "unite",                // Unified IP & Telecom Environment
	Port{3217, 17}:   "unite",                // Unified IP & Telecom Environment
	Port{3218, 6}:    "smartpackets",         // EMC SmartPackets
	Port{3218, 17}:   "smartpackets",         // EMC SmartPackets
	Port{3219, 6}:    "wms-messenger",        // WMS Messenger
	Port{3219, 17}:   "wms-messenger",        // WMS Messenger
	Port{3220, 6}:    "xnm-ssl",              // XML NM over SSL
	Port{3220, 17}:   "xnm-ssl",              // XML NM over SSL
	Port{3221, 6}:    "xnm-clear-text",       // XML NM over TCP
	Port{3221, 17}:   "xnm-clear-text",       // XML NM over TCP
	Port{3222, 6}:    "glbp",                 // Gateway Load Balancing Pr
	Port{3222, 17}:   "glbp",                 // Gateway Load Balancing Pr
	Port{3223, 6}:    "digivote",             // DIGIVOTE (R) Vote-Server
	Port{3223, 17}:   "digivote",             // DIGIVOTE (R) Vote-Server
	Port{3224, 6}:    "aes-discovery",        // AES Discovery Port
	Port{3224, 17}:   "aes-discovery",        // AES Discovery Port
	Port{3225, 6}:    "fcip-port",            // FCIP
	Port{3225, 17}:   "fcip-port",            // FCIP
	Port{3226, 6}:    "isi-irp",              // ISI Industry Software IRP
	Port{3226, 17}:   "isi-irp",              // ISI Industry Software IRP
	Port{3227, 6}:    "dwnmshttp",            // DiamondWave NMS Server
	Port{3227, 17}:   "dwnmshttp",            // DiamondWave NMS Server
	Port{3228, 6}:    "dwmsgserver",          // DiamondWave MSG Server
	Port{3228, 17}:   "dwmsgserver",          // DiamondWave MSG Server
	Port{3229, 6}:    "global-cd-port",       // Global CD Port
	Port{3229, 17}:   "global-cd-port",       // Global CD Port
	Port{3230, 6}:    "sftdst-port",          // Software Distributor Port
	Port{3230, 17}:   "sftdst-port",          // Software Distributor Port
	Port{3231, 6}:    "vidigo",               // VidiGo communication (previous was: Delta Solutions Direct)
	Port{3231, 17}:   "vidigo",               // VidiGo communication (previous was: Delta Solutions Direct)
	Port{3232, 6}:    "mdtp",                 // MDT port
	Port{3232, 17}:   "mdtp",                 // MDT port
	Port{3233, 6}:    "whisker",              // WhiskerControl main port
	Port{3233, 17}:   "whisker",              // WhiskerControl main port
	Port{3234, 6}:    "alchemy",              // Alchemy Server
	Port{3234, 17}:   "alchemy",              // Alchemy Server
	Port{3235, 6}:    "mdap-port",            // MDAP port | MDAP Port
	Port{3235, 17}:   "mdap-port",            // MDAP Port
	Port{3236, 6}:    "apparenet-ts",         // appareNet Test Server
	Port{3236, 17}:   "apparenet-ts",         // appareNet Test Server
	Port{3237, 6}:    "apparenet-tps",        // appareNet Test Packet Sequencer
	Port{3237, 17}:   "apparenet-tps",        // appareNet Test Packet Sequencer
	Port{3238, 6}:    "apparenet-as",         // appareNet Analysis Server
	Port{3238, 17}:   "apparenet-as",         // appareNet Analysis Server
	Port{3239, 6}:    "apparenet-ui",         // appareNet User Interface
	Port{3239, 17}:   "apparenet-ui",         // appareNet User Interface
	Port{3240, 6}:    "triomotion",           // Trio Motion Control Port
	Port{3240, 17}:   "triomotion",           // Trio Motion Control Port
	Port{3241, 6}:    "sysorb",               // SysOrb Monitoring Server
	Port{3241, 17}:   "sysorb",               // SysOrb Monitoring Server
	Port{3242, 6}:    "sdp-id-port",          // Session Description ID
	Port{3242, 17}:   "sdp-id-port",          // Session Description ID
	Port{3243, 6}:    "timelot",              // Timelot Port
	Port{3243, 17}:   "timelot",              // Timelot Port
	Port{3244, 6}:    "onesaf",               // Missing description for onesaf
	Port{3244, 17}:   "onesaf",               // OneSAF
	Port{3245, 6}:    "vieo-fe",              // VIEO Fabric Executive
	Port{3245, 17}:   "vieo-fe",              // VIEO Fabric Executive
	Port{3246, 6}:    "dvt-system",           // DVT SYSTEM PORT
	Port{3246, 17}:   "kademlia",             // Kademlia P2P (mlnet)
	Port{3247, 6}:    "dvt-data",             // DVT DATA LINK
	Port{3247, 17}:   "dvt-data",             // DVT DATA LINK
	Port{3248, 6}:    "procos-lm",            // PROCOS LM
	Port{3248, 17}:   "procos-lm",            // PROCOS LM
	Port{3249, 6}:    "ssp",                  // State Sync Protocol
	Port{3249, 17}:   "ssp",                  // State Sync Protocol
	Port{3250, 6}:    "hicp",                 // HMS hicp port
	Port{3250, 17}:   "hicp",                 // HMS hicp port
	Port{3251, 6}:    "sysscanner",           // Sys Scanner
	Port{3251, 17}:   "sysscanner",           // Sys Scanner
	Port{3252, 6}:    "dhe",                  // DHE port
	Port{3252, 17}:   "dhe",                  // DHE port
	Port{3253, 6}:    "pda-data",             // PDA Data
	Port{3253, 17}:   "pda-data",             // PDA Data
	Port{3254, 6}:    "pda-sys",              // PDA System
	Port{3254, 17}:   "pda-sys",              // PDA System
	Port{3255, 6}:    "semaphore",            // Semaphore Connection Port
	Port{3255, 17}:   "semaphore",            // Semaphore Connection Port
	Port{3256, 6}:    "cpqrpm-agent",         // Compaq RPM Agent Port
	Port{3256, 17}:   "cpqrpm-agent",         // Compaq RPM Agent Port
	Port{3257, 6}:    "cpqrpm-server",        // Compaq RPM Server Port
	Port{3257, 17}:   "cpqrpm-server",        // Compaq RPM Server Port
	Port{3258, 6}:    "ivecon-port",          // Ivecon Server Port
	Port{3258, 17}:   "ivecon-port",          // Ivecon Server Port
	Port{3259, 6}:    "epncdp2",              // Epson Network Common Devi
	Port{3259, 17}:   "epncdp2",              // Epson Network Common Devi
	Port{3260, 6}:    "iscsi",                // iscsi-target | iSCSI port
	Port{3260, 17}:   "iscsi",                // iSCSI
	Port{3261, 6}:    "winshadow",            // Missing description for winshadow
	Port{3261, 17}:   "winshadow",            // winShadow
	Port{3262, 6}:    "necp",                 // Missing description for necp
	Port{3262, 17}:   "necp",                 // NECP
	Port{3263, 6}:    "ecolor-imager",        // E-Color Enterprise Imager
	Port{3263, 17}:   "ecolor-imager",        // E-Color Enterprise Imager
	Port{3264, 6}:    "ccmail",               // cc:mail lotus
	Port{3264, 17}:   "ccmail",               // cc:mail lotus
	Port{3265, 6}:    "altav-tunnel",         // Altav Tunnel
	Port{3265, 17}:   "altav-tunnel",         // Altav Tunnel
	Port{3266, 6}:    "ns-cfg-server",        // NS CFG Server
	Port{3266, 17}:   "ns-cfg-server",        // NS CFG Server
	Port{3267, 6}:    "ibm-dial-out",         // IBM Dial Out
	Port{3267, 17}:   "ibm-dial-out",         // IBM Dial Out
	Port{3268, 6}:    "globalcatLDAP",        // msft-gc | Global Catalog LDAP | Microsoft Global Catalog
	Port{3268, 17}:   "msft-gc",              // Microsoft Global Catalog
	Port{3269, 6}:    "globalcatLDAPssl",     // msft-gc-ssl | Global Catalog LDAP over ssl | Microsoft Global Catalog with LDAP SSL
	Port{3269, 17}:   "msft-gc-ssl",          // Microsoft Global Catalog with LDAP SSL
	Port{3270, 6}:    "verismart",            // Missing description for verismart
	Port{3270, 17}:   "verismart",            // Verismart
	Port{3271, 6}:    "csoft-prev",           // CSoft Prev Port
	Port{3271, 17}:   "csoft-prev",           // CSoft Prev Port
	Port{3272, 6}:    "user-manager",         // Fujitsu User Manager
	Port{3272, 17}:   "user-manager",         // Fujitsu User Manager
	Port{3273, 6}:    "sxmp",                 // Simple Extensible Multiplexed Protocol
	Port{3273, 17}:   "sxmp",                 // Simple Extensible Multiplexed Protocol
	Port{3274, 6}:    "ordinox-server",       // Ordinox Server
	Port{3274, 17}:   "ordinox-server",       // Ordinox Server
	Port{3275, 6}:    "samd",                 // Missing description for samd
	Port{3275, 17}:   "samd",                 // SAMD
	Port{3276, 6}:    "maxim-asics",          // Maxim ASICs
	Port{3276, 17}:   "maxim-asics",          // Maxim ASICs
	Port{3277, 6}:    "awg-proxy",            // AWG Proxy
	Port{3277, 17}:   "awg-proxy",            // AWG Proxy
	Port{3278, 6}:    "lkcmserver",           // LKCM Server
	Port{3278, 17}:   "lkcmserver",           // LKCM Server
	Port{3279, 6}:    "admind",               // Missing description for admind
	Port{3279, 17}:   "admind",               // Missing description for admind
	Port{3280, 6}:    "vs-server",            // VS Server
	Port{3280, 17}:   "vs-server",            // VS Server
	Port{3281, 6}:    "sysopt",               // Missing description for sysopt
	Port{3281, 17}:   "sysopt",               // SYSOPT
	Port{3282, 6}:    "datusorb",             // Missing description for datusorb
	Port{3282, 17}:   "datusorb",             // Datusorb
	Port{3283, 6}:    "netassistant",         // #ERROR:Apple Remote Desktop (Net Assistant) | Apple Remote Desktop Net Assistant reporting feature | Net Assistant
	Port{3283, 17}:   "netassistant",         // Apple Remote Desktop Net Assistant reporting feature
	Port{3284, 6}:    "4talk",                // Missing description for 4talk
	Port{3284, 17}:   "4talk",                // 4Talk
	Port{3285, 6}:    "plato",                // Missing description for plato
	Port{3285, 17}:   "plato",                // Plato
	Port{3286, 6}:    "e-net",                // Missing description for e-net
	Port{3286, 17}:   "e-net",                // E-Net
	Port{3287, 6}:    "directvdata",          // Missing description for directvdata
	Port{3287, 17}:   "directvdata",          // DIRECTVDATA
	Port{3288, 6}:    "cops",                 // Missing description for cops
	Port{3288, 17}:   "cops",                 // COPS
	Port{3289, 6}:    "enpc",                 // Missing description for enpc
	Port{3289, 17}:   "enpc",                 // ENPC
	Port{3290, 6}:    "caps-lm",              // CAPS LOGISTICS TOOLKIT - LM
	Port{3290, 17}:   "caps-lm",              // CAPS LOGISTICS TOOLKIT - LM
	Port{3291, 6}:    "sah-lm",               // S A Holditch & Associates - LM
	Port{3291, 17}:   "sah-lm",               // S A Holditch & Associates - LM
	Port{3292, 6}:    "meetingmaker",         // cart-o-rama | Meeting maker time management software | Cart O Rama
	Port{3292, 17}:   "cart-o-rama",          // Cart O Rama
	Port{3293, 6}:    "fg-fps",               // Missing description for fg-fps
	Port{3293, 17}:   "fg-fps",               // Missing description for fg-fps
	Port{3294, 6}:    "fg-gip",               // Missing description for fg-gip
	Port{3294, 17}:   "fg-gip",               // Missing description for fg-gip
	Port{3295, 6}:    "dyniplookup",          // Dynamic IP Lookup
	Port{3295, 17}:   "dyniplookup",          // Dynamic IP Lookup
	Port{3296, 6}:    "rib-slm",              // Rib License Manager
	Port{3296, 17}:   "rib-slm",              // Rib License Manager
	Port{3297, 6}:    "cytel-lm",             // Cytel License Manager
	Port{3297, 17}:   "cytel-lm",             // Cytel License Manager
	Port{3298, 6}:    "deskview",             // Missing description for deskview
	Port{3298, 17}:   "deskview",             // DeskView
	Port{3299, 6}:    "saprouter",            // pdrncs
	Port{3299, 17}:   "pdrncs",               // Missing description for pdrncs
	Port{3300, 6}:    "ceph",                 // Ceph monitor
	Port{3302, 6}:    "mcs-fastmail",         // MCS Fastmail
	Port{3302, 17}:   "mcs-fastmail",         // MCS Fastmail
	Port{3303, 6}:    "opsession-clnt",       // OP Session Client
	Port{3303, 17}:   "opsession-clnt",       // OP Session Client
	Port{3304, 6}:    "opsession-srvr",       // OP Session Server
	Port{3304, 17}:   "opsession-srvr",       // OP Session Server
	Port{3305, 6}:    "odette-ftp",           // Missing description for odette-ftp
	Port{3305, 17}:   "odette-ftp",           // ODETTE-FTP
	Port{3306, 6}:    "mysql",                // Missing description for mysql
	Port{3306, 17}:   "mysql",                // MySQL
	Port{3307, 6}:    "opsession-prxy",       // OP Session Proxy
	Port{3307, 17}:   "opsession-prxy",       // OP Session Proxy
	Port{3308, 6}:    "tns-server",           // TNS Server
	Port{3308, 17}:   "tns-server",           // TNS Server
	Port{3309, 6}:    "tns-adv",              // TNS ADV
	Port{3309, 17}:   "tns-adv",              // TNS ADV
	Port{3310, 6}:    "dyna-access",          // Dyna Access
	Port{3310, 17}:   "dyna-access",          // Dyna Access
	Port{3311, 6}:    "mcns-tel-ret",         // MCNS Tel Ret
	Port{3311, 17}:   "mcns-tel-ret",         // MCNS Tel Ret
	Port{3312, 6}:    "appman-server",        // Application Management Server
	Port{3312, 17}:   "appman-server",        // Application Management Server
	Port{3313, 6}:    "uorb",                 // Unify Object Broker
	Port{3313, 17}:   "uorb",                 // Unify Object Broker
	Port{3314, 6}:    "uohost",               // Unify Object Host
	Port{3314, 17}:   "uohost",               // Unify Object Host
	Port{3315, 6}:    "cdid",                 // Missing description for cdid
	Port{3315, 17}:   "cdid",                 // CDID
	Port{3316, 6}:    "aicc-cmi",             // AICC CMI
	Port{3316, 17}:   "aicc-cmi",             // AICC CMI
	Port{3317, 6}:    "vsaiport",             // VSAI PORT
	Port{3317, 17}:   "vsaiport",             // VSAI PORT
	Port{3318, 6}:    "ssrip",                // Swith to Swith Routing Information Protocol
	Port{3318, 17}:   "ssrip",                // Swith to Swith Routing Information Protocol
	Port{3319, 6}:    "sdt-lmd",              // SDT License Manager
	Port{3319, 17}:   "sdt-lmd",              // SDT License Manager
	Port{3320, 6}:    "officelink2000",       // Office Link 2000
	Port{3320, 17}:   "officelink2000",       // Office Link 2000
	Port{3321, 6}:    "vnsstr",               // Missing description for vnsstr
	Port{3321, 17}:   "vnsstr",               // VNSSTR
	Port{3322, 6}:    "active-net",           // Active Networks
	Port{3322, 17}:   "active-net",           // Active Networks
	Port{3323, 6}:    "active-net",           // Active Networks
	Port{3323, 17}:   "active-net",           // Active Networks
	Port{3324, 6}:    "active-net",           // Active Networks
	Port{3324, 17}:   "active-net",           // Active Networks
	Port{3325, 6}:    "active-net",           // Active Networks
	Port{3325, 17}:   "active-net",           // Active Networks
	Port{3326, 6}:    "sftu",                 // Missing description for sftu
	Port{3326, 17}:   "sftu",                 // SFTU
	Port{3327, 6}:    "bbars",                // Missing description for bbars
	Port{3327, 17}:   "bbars",                // BBARS
	Port{3328, 6}:    "egptlm",               // Eaglepoint License Manager
	Port{3328, 17}:   "egptlm",               // Eaglepoint License Manager
	Port{3329, 6}:    "hp-device-disc",       // HP Device Disc
	Port{3329, 17}:   "hp-device-disc",       // HP Device Disc
	Port{3330, 6}:    "mcs-calypsoicf",       // MCS Calypso ICF
	Port{3330, 17}:   "mcs-calypsoicf",       // MCS Calypso ICF
	Port{3331, 6}:    "mcs-messaging",        // MCS Messaging
	Port{3331, 17}:   "mcs-messaging",        // MCS Messaging
	Port{3332, 6}:    "mcs-mailsvr",          // MCS Mail Server
	Port{3332, 17}:   "mcs-mailsvr",          // MCS Mail Server
	Port{3333, 6}:    "dec-notes",            // DEC Notes
	Port{3333, 17}:   "dec-notes",            // DEC Notes
	Port{3334, 6}:    "directv-web",          // Direct TV Webcasting
	Port{3334, 17}:   "directv-web",          // Direct TV Webcasting
	Port{3335, 6}:    "directv-soft",         // Direct TV Software Updates
	Port{3335, 17}:   "directv-soft",         // Direct TV Software Updates
	Port{3336, 6}:    "directv-tick",         // Direct TV Tickers
	Port{3336, 17}:   "directv-tick",         // Direct TV Tickers
	Port{3337, 6}:    "directv-catlg",        // Direct TV Data Catalog
	Port{3337, 17}:   "directv-catlg",        // Direct TV Data Catalog
	Port{3338, 6}:    "anet-b",               // OMF data b
	Port{3338, 17}:   "anet-b",               // OMF data b
	Port{3339, 6}:    "anet-l",               // OMF data l
	Port{3339, 17}:   "anet-l",               // OMF data l
	Port{3340, 6}:    "anet-m",               // OMF data m
	Port{3340, 17}:   "anet-m",               // OMF data m
	Port{3341, 6}:    "anet-h",               // OMF data h
	Port{3341, 17}:   "anet-h",               // OMF data h
	Port{3342, 6}:    "webtie",               // Missing description for webtie
	Port{3342, 17}:   "webtie",               // WebTIE
	Port{3343, 6}:    "ms-cluster-net",       // MS Cluster Net
	Port{3343, 17}:   "ms-cluster-net",       // MS Cluster Net
	Port{3344, 6}:    "bnt-manager",          // BNT Manager
	Port{3344, 17}:   "bnt-manager",          // BNT Manager
	Port{3345, 6}:    "influence",            // Missing description for influence
	Port{3345, 17}:   "influence",            // Influence
	Port{3346, 6}:    "trnsprntproxy",        // Trnsprnt Proxy
	Port{3346, 17}:   "trnsprntproxy",        // Trnsprnt Proxy
	Port{3347, 6}:    "phoenix-rpc",          // Phoenix RPC
	Port{3347, 17}:   "phoenix-rpc",          // Phoenix RPC
	Port{3348, 6}:    "pangolin-laser",       // Pangolin Laser
	Port{3348, 17}:   "pangolin-laser",       // Pangolin Laser
	Port{3349, 6}:    "chevinservices",       // Chevin Services
	Port{3349, 17}:   "chevinservices",       // Chevin Services
	Port{3350, 6}:    "findviatv",            // Missing description for findviatv
	Port{3350, 17}:   "findviatv",            // FINDVIATV
	Port{3351, 6}:    "btrieve",              // Btrieve port
	Port{3351, 17}:   "btrieve",              // Btrieve port
	Port{3352, 6}:    "ssql",                 // Scalable SQL
	Port{3352, 17}:   "ssql",                 // Scalable SQL
	Port{3353, 6}:    "fatpipe",              // Missing description for fatpipe
	Port{3353, 17}:   "fatpipe",              // FATPIPE
	Port{3354, 6}:    "suitjd",               // Missing description for suitjd
	Port{3354, 17}:   "suitjd",               // SUITJD
	Port{3355, 6}:    "ordinox-dbase",        // Ordinox Dbase
	Port{3355, 17}:   "ordinox-dbase",        // Ordinox Dbase
	Port{3356, 6}:    "upnotifyps",           // Missing description for upnotifyps
	Port{3356, 17}:   "upnotifyps",           // UPNOTIFYPS
	Port{3357, 6}:    "adtech-test",          // Adtech Test IP
	Port{3357, 17}:   "adtech-test",          // Adtech Test IP
	Port{3358, 6}:    "mpsysrmsvr",           // Mp Sys Rmsvr
	Port{3358, 17}:   "mpsysrmsvr",           // Mp Sys Rmsvr
	Port{3359, 6}:    "wg-netforce",          // WG NetForce
	Port{3359, 17}:   "wg-netforce",          // WG NetForce
	Port{3360, 6}:    "kv-server",            // KV Server
	Port{3360, 17}:   "kv-server",            // KV Server
	Port{3361, 6}:    "kv-agent",             // KV Agent
	Port{3361, 17}:   "kv-agent",             // KV Agent
	Port{3362, 6}:    "dj-ilm",               // DJ ILM
	Port{3362, 17}:   "dj-ilm",               // DJ ILM
	Port{3363, 6}:    "nati-vi-server",       // NATI Vi Server
	Port{3363, 17}:   "nati-vi-server",       // NATI Vi Server
	Port{3364, 6}:    "creativeserver",       // Creative Server
	Port{3364, 17}:   "creativeserver",       // Creative Server
	Port{3365, 6}:    "contentserver",        // Content Server
	Port{3365, 17}:   "contentserver",        // Content Server
	Port{3366, 6}:    "creativepartnr",       // Creative Partner
	Port{3366, 17}:   "creativepartnr",       // Creative Partner
	Port{3367, 6}:    "satvid-datalnk",       // Satellite Video Data Link
	Port{3367, 17}:   "satvid-datalnk",       // Satellite Video Data Link
	Port{3368, 6}:    "satvid-datalnk",       // Satellite Video Data Link
	Port{3368, 17}:   "satvid-datalnk",       // Satellite Video Data Link
	Port{3369, 6}:    "satvid-datalnk",       // Satellite Video Data Link
	Port{3369, 17}:   "satvid-datalnk",       // Satellite Video Data Link
	Port{3370, 6}:    "satvid-datalnk",       // Satellite Video Data Link
	Port{3370, 17}:   "satvid-datalnk",       // Satellite Video Data Link
	Port{3371, 6}:    "satvid-datalnk",       // Satellite Video Data Link
	Port{3371, 17}:   "satvid-datalnk",       // Satellite Video Data Link
	Port{3372, 6}:    "msdtc",                // tip2 | MS distributed transaction coordinator | TIP 2
	Port{3372, 17}:   "tip2",                 // TIP 2
	Port{3373, 6}:    "lavenir-lm",           // Lavenir License Manager
	Port{3373, 17}:   "lavenir-lm",           // Lavenir License Manager
	Port{3374, 6}:    "cluster-disc",         // Cluster Disc
	Port{3374, 17}:   "cluster-disc",         // Cluster Disc
	Port{3375, 6}:    "vsnm-agent",           // VSNM Agent
	Port{3375, 17}:   "vsnm-agent",           // VSNM Agent
	Port{3376, 6}:    "cdbroker",             // CD Broker
	Port{3376, 17}:   "cdbroker",             // CD Broker
	Port{3377, 6}:    "cogsys-lm",            // Cogsys Network License Manager
	Port{3377, 17}:   "cogsys-lm",            // Cogsys Network License Manager
	Port{3378, 6}:    "wsicopy",              // Missing description for wsicopy
	Port{3378, 17}:   "wsicopy",              // WSICOPY
	Port{3379, 6}:    "socorfs",              // Missing description for socorfs
	Port{3379, 17}:   "socorfs",              // SOCORFS
	Port{3380, 6}:    "sns-channels",         // SNS Channels
	Port{3380, 17}:   "sns-channels",         // SNS Channels
	Port{3381, 6}:    "geneous",              // Missing description for geneous
	Port{3381, 17}:   "geneous",              // Geneous
	Port{3382, 6}:    "fujitsu-neat",         // Fujitsu Network Enhanced Antitheft function
	Port{3382, 17}:   "fujitsu-neat",         // Fujitsu Network Enhanced Antitheft function
	Port{3383, 6}:    "esp-lm",               // Enterprise Software Products License Manager
	Port{3383, 17}:   "esp-lm",               // Enterprise Software Products License Manager
	Port{3384, 6}:    "hp-clic",              // Cluster Management Services | Hardware Management
	Port{3384, 17}:   "hp-clic",              // Hardware Management
	Port{3385, 6}:    "qnxnetman",            // Missing description for qnxnetman
	Port{3385, 17}:   "qnxnetman",            // Missing description for qnxnetman
	Port{3386, 6}:    "gprs-data",            // gprs-sig | GPRS Data | GPRS SIG
	Port{3386, 17}:   "gprs-sig",             // GPRS SIG
	Port{3387, 6}:    "backroomnet",          // Back Room Net
	Port{3387, 17}:   "backroomnet",          // Back Room Net
	Port{3388, 6}:    "cbserver",             // CB Server
	Port{3388, 17}:   "cbserver",             // CB Server
	Port{3389, 6}:    "ms-wbt-server",        // Microsoft Remote Display Protocol (aka ms-term-serv, microsoft-rdp) | MS WBT Server
	Port{3389, 17}:   "ms-wbt-server",        // Microsoft Remote Display Protocol (aka ms-term-serv, microsoft-rdp)
	Port{3390, 6}:    "dsc",                  // Distributed Service Coordinator
	Port{3390, 17}:   "dsc",                  // Distributed Service Coordinator
	Port{3391, 6}:    "savant",               // Missing description for savant
	Port{3391, 17}:   "savant",               // SAVANT
	Port{3392, 6}:    "efi-lm",               // EFI License Management
	Port{3392, 17}:   "efi-lm",               // EFI License Management
	Port{3393, 6}:    "d2k-tapestry1",        // D2K Tapestry Client to Server
	Port{3393, 17}:   "d2k-tapestry1",        // D2K Tapestry Client to Server
	Port{3394, 6}:    "d2k-tapestry2",        // D2K Tapestry Server to Server
	Port{3394, 17}:   "d2k-tapestry2",        // D2K Tapestry Server to Server
	Port{3395, 6}:    "dyna-lm",              // Dyna License Manager (Elam)
	Port{3395, 17}:   "dyna-lm",              // Dyna License Manager (Elam)
	Port{3396, 6}:    "printer_agent",        // printer-agent | Printer Agent
	Port{3396, 17}:   "printer_agent",        // Printer Agent
	Port{3397, 6}:    "saposs",               // cloanto-lm | SAP Oss | Cloanto License Manager
	Port{3397, 17}:   "cloanto-lm",           // Cloanto License Manager
	Port{3398, 6}:    "sapcomm",              // mercantile | Mercantile
	Port{3398, 17}:   "mercantile",           // Mercantile
	Port{3399, 6}:    "sapeps",               // csms | SAP EPS | CSMS
	Port{3399, 17}:   "csms",                 // CSMS
	Port{3400, 6}:    "csms2",                // Missing description for csms2
	Port{3400, 17}:   "csms2",                // CSMS2
	Port{3401, 6}:    "filecast",             // Missing description for filecast
	Port{3401, 17}:   "squid-snmp",           // Squid proxy SNMP port
	Port{3402, 6}:    "fxaengine-net",        // FXa Engine Network Port
	Port{3402, 17}:   "fxaengine-net",        // FXa Engine Network Port
	Port{3405, 6}:    "nokia-ann-ch1",        // Nokia Announcement ch 1
	Port{3405, 17}:   "nokia-ann-ch1",        // Nokia Announcement ch 1
	Port{3406, 6}:    "nokia-ann-ch2",        // Nokia Announcement ch 2
	Port{3406, 17}:   "nokia-ann-ch2",        // Nokia Announcement ch 2
	Port{3407, 6}:    "ldap-admin",           // LDAP admin server port
	Port{3407, 17}:   "ldap-admin",           // LDAP admin server port
	Port{3408, 6}:    "BESApi",               // BES Api Port
	Port{3408, 17}:   "BESApi",               // BES Api Port
	Port{3409, 6}:    "networklens",          // NetworkLens Event Port
	Port{3409, 17}:   "networklens",          // NetworkLens Event Port
	Port{3410, 6}:    "networklenss",         // NetworkLens SSL Event
	Port{3410, 17}:   "networklenss",         // NetworkLens SSL Event
	Port{3411, 6}:    "biolink-auth",         // BioLink Authenteon server
	Port{3411, 17}:   "biolink-auth",         // BioLink Authenteon server
	Port{3412, 6}:    "xmlblaster",           // Missing description for xmlblaster
	Port{3412, 17}:   "xmlblaster",           // xmlBlaster
	Port{3413, 6}:    "svnet",                // SpecView Networking
	Port{3413, 17}:   "svnet",                // SpecView Networking
	Port{3414, 6}:    "wip-port",             // BroadCloud WIP Port
	Port{3414, 17}:   "wip-port",             // BroadCloud WIP Port
	Port{3415, 6}:    "bcinameservice",       // BCI Name Service
	Port{3415, 17}:   "bcinameservice",       // BCI Name Service
	Port{3416, 6}:    "commandport",          // AirMobile IS Command Port
	Port{3416, 17}:   "commandport",          // AirMobile IS Command Port
	Port{3417, 6}:    "csvr",                 // ConServR file translation
	Port{3417, 17}:   "csvr",                 // ConServR file translation
	Port{3418, 6}:    "rnmap",                // Remote nmap
	Port{3418, 17}:   "rnmap",                // Remote nmap
	Port{3419, 6}:    "softaudit",            // Isogon SoftAudit | ISogon SoftAudit
	Port{3419, 17}:   "softaudit",            // ISogon SoftAudit
	Port{3420, 6}:    "ifcp-port",            // iFCP User Port
	Port{3420, 17}:   "ifcp-port",            // iFCP User Port
	Port{3421, 6}:    "bmap",                 // Bull Apprise portmapper
	Port{3421, 17}:   "bmap",                 // Bull Apprise portmapper
	Port{3422, 6}:    "rusb-sys-port",        // Remote USB System Port
	Port{3422, 17}:   "rusb-sys-port",        // Remote USB System Port
	Port{3423, 6}:    "xtrm",                 // xTrade Reliable Messaging
	Port{3423, 17}:   "xtrm",                 // xTrade Reliable Messaging
	Port{3424, 6}:    "xtrms",                // xTrade over TLS SSL
	Port{3424, 17}:   "xtrms",                // xTrade over TLS SSL
	Port{3425, 6}:    "agps-port",            // AGPS Access Port
	Port{3425, 17}:   "agps-port",            // AGPS Access Port
	Port{3426, 6}:    "arkivio",              // Arkivio Storage Protocol
	Port{3426, 17}:   "arkivio",              // Arkivio Storage Protocol
	Port{3427, 6}:    "websphere-snmp",       // WebSphere SNMP
	Port{3427, 17}:   "websphere-snmp",       // WebSphere SNMP
	Port{3428, 6}:    "twcss",                // 2Wire CSS
	Port{3428, 17}:   "twcss",                // 2Wire CSS
	Port{3429, 6}:    "gcsp",                 // GCSP user port
	Port{3429, 17}:   "gcsp",                 // GCSP user port
	Port{3430, 6}:    "ssdispatch",           // Scott Studios Dispatch
	Port{3430, 17}:   "ssdispatch",           // Scott Studios Dispatch
	Port{3431, 6}:    "ndl-als",              // Active License Server Port
	Port{3431, 17}:   "ndl-als",              // Active License Server Port
	Port{3432, 6}:    "osdcp",                // Secure Device Protocol
	Port{3432, 17}:   "osdcp",                // Secure Device Protocol
	Port{3433, 6}:    "alta-smp",             // opnet-smp | Altaworks Service Management Platform | OPNET Service Management Platform
	Port{3433, 17}:   "alta-smp",             // Altaworks Service Management Platform
	Port{3434, 6}:    "opencm",               // OpenCM Server
	Port{3434, 17}:   "opencm",               // OpenCM Server
	Port{3435, 6}:    "pacom",                // Pacom Security User Port
	Port{3435, 17}:   "pacom",                // Pacom Security User Port
	Port{3436, 6}:    "gc-config",            // GuardControl Exchange Protocol
	Port{3436, 17}:   "gc-config",            // GuardControl Exchange Protocol
	Port{3437, 6}:    "autocueds",            // Autocue Directory Service
	Port{3437, 17}:   "autocueds",            // Autocue Directory Service
	Port{3438, 6}:    "spiral-admin",         // Spiralcraft Admin
	Port{3438, 17}:   "spiral-admin",         // Spiralcraft Admin
	Port{3439, 6}:    "hri-port",             // HRI Interface Port
	Port{3439, 17}:   "hri-port",             // HRI Interface Port
	Port{3440, 6}:    "ans-console",          // Net Steward Mgmt Console
	Port{3440, 17}:   "ans-console",          // Net Steward Mgmt Console
	Port{3441, 6}:    "connect-client",       // OC Connect Client
	Port{3441, 17}:   "connect-client",       // OC Connect Client
	Port{3442, 6}:    "connect-server",       // OC Connect Server
	Port{3442, 17}:   "connect-server",       // OC Connect Server
	Port{3443, 6}:    "ov-nnm-websrv",        // OpenView Network Node Manager WEB Server
	Port{3443, 17}:   "ov-nnm-websrv",        // OpenView Network Node Manager WEB Server
	Port{3444, 6}:    "denali-server",        // Denali Server
	Port{3444, 17}:   "denali-server",        // Denali Server
	Port{3445, 6}:    "monp",                 // Media Object Network
	Port{3445, 17}:   "monp",                 // Media Object Network
	Port{3446, 6}:    "3comfaxrpc",           // 3Com FAX RPC port
	Port{3446, 17}:   "3comfaxrpc",           // 3Com FAX RPC port
	Port{3447, 6}:    "directnet",            // DirectNet IM System
	Port{3447, 17}:   "directnet",            // DirectNet IM System
	Port{3448, 6}:    "dnc-port",             // Discovery and Net Config
	Port{3448, 17}:   "dnc-port",             // Discovery and Net Config
	Port{3449, 6}:    "hotu-chat",            // HotU Chat
	Port{3449, 17}:   "hotu-chat",            // HotU Chat
	Port{3450, 6}:    "castorproxy",          // Missing description for castorproxy
	Port{3450, 17}:   "castorproxy",          // CAStorProxy
	Port{3451, 6}:    "asam",                 // ASAM Services
	Port{3451, 17}:   "asam",                 // ASAM Services
	Port{3452, 6}:    "sabp-signal",          // SABP-Signalling Protocol
	Port{3452, 17}:   "sabp-signal",          // SABP-Signalling Protocol
	Port{3453, 6}:    "pscupd",               // PSC Update Port | PSC Update
	Port{3453, 17}:   "pscupd",               // PSC Update Port
	Port{3454, 6}:    "mira",                 // Apple Remote Access Protocol
	Port{3455, 6}:    "prsvp",                // RSVP Port
	Port{3455, 17}:   "prsvp",                // RSVP Port
	Port{3456, 6}:    "vat",                  // VAT default data
	Port{3456, 17}:   "IISrpc-or-vat",        // also VAT default data
	Port{3457, 6}:    "vat-control",          // VAT default control
	Port{3457, 17}:   "vat-control",          // VAT default control
	Port{3458, 6}:    "d3winosfi",            // Missing description for d3winosfi
	Port{3458, 17}:   "d3winosfi",            // D3WinOSFI
	Port{3459, 6}:    "integral",             // TIP Integral
	Port{3459, 17}:   "integral",             // TIP Integral
	Port{3460, 6}:    "edm-manager",          // EDM Manger
	Port{3460, 17}:   "edm-manager",          // EDM Manger
	Port{3461, 6}:    "edm-stager",           // EDM Stager
	Port{3461, 17}:   "edm-stager",           // EDM Stager
	Port{3462, 6}:    "track",                // edm-std-notify | software distribution | EDM STD Notify
	Port{3462, 17}:   "edm-std-notify",       // EDM STD Notify
	Port{3463, 6}:    "edm-adm-notify",       // EDM ADM Notify
	Port{3463, 17}:   "edm-adm-notify",       // EDM ADM Notify
	Port{3464, 6}:    "edm-mgr-sync",         // EDM MGR Sync
	Port{3464, 17}:   "edm-mgr-sync",         // EDM MGR Sync
	Port{3465, 6}:    "edm-mgr-cntrl",        // EDM MGR Cntrl
	Port{3465, 17}:   "edm-mgr-cntrl",        // EDM MGR Cntrl
	Port{3466, 6}:    "workflow",             // Missing description for workflow
	Port{3466, 17}:   "workflow",             // WORKFLOW
	Port{3467, 6}:    "rcst",                 // Missing description for rcst
	Port{3467, 17}:   "rcst",                 // RCST
	Port{3468, 6}:    "ttcmremotectrl",       // TTCM Remote Controll
	Port{3468, 17}:   "ttcmremotectrl",       // TTCM Remote Controll
	Port{3469, 6}:    "pluribus",             // Missing description for pluribus
	Port{3469, 17}:   "pluribus",             // Pluribus
	Port{3470, 6}:    "jt400",                // Missing description for jt400
	Port{3470, 17}:   "jt400",                // Missing description for jt400
	Port{3471, 6}:    "jt400-ssl",            // Missing description for jt400-ssl
	Port{3471, 17}:   "jt400-ssl",            // Missing description for jt400-ssl
	Port{3472, 6}:    "jaugsremotec-1",       // JAUGS N-G Remotec 1
	Port{3472, 17}:   "jaugsremotec-1",       // JAUGS N-G Remotec 1
	Port{3473, 6}:    "jaugsremotec-2",       // JAUGS N-G Remotec 2
	Port{3473, 17}:   "jaugsremotec-2",       // JAUGS N-G Remotec 2
	Port{3474, 6}:    "ttntspauto",           // TSP Automation
	Port{3474, 17}:   "ttntspauto",           // TSP Automation
	Port{3475, 6}:    "genisar-port",         // Genisar Comm Port
	Port{3475, 17}:   "genisar-port",         // Genisar Comm Port
	Port{3476, 6}:    "nppmp",                // NVIDIA Mgmt Protocol
	Port{3476, 17}:   "nppmp",                // NVIDIA Mgmt Protocol
	Port{3477, 6}:    "ecomm",                // eComm link port
	Port{3477, 17}:   "ecomm",                // eComm link port
	Port{3478, 6}:    "stun",                 // stun-behavior | turn | Session Traversal Utilities for NAT (STUN) port | TURN over TCP | TURN over UDP | STUN Behavior Discovery over TCP | STUN Behavior Discovery over UDP
	Port{3478, 17}:   "stun",                 // Session Traversal Utilities for NAT (STUN) port
	Port{3479, 6}:    "twrpc",                // 2Wire RPC
	Port{3479, 17}:   "twrpc",                // 2Wire RPC
	Port{3480, 6}:    "plethora",             // Secure Virtual Workspace
	Port{3480, 17}:   "plethora",             // Secure Virtual Workspace
	Port{3481, 6}:    "cleanerliverc",        // CleanerLive remote ctrl
	Port{3481, 17}:   "cleanerliverc",        // CleanerLive remote ctrl
	Port{3482, 6}:    "vulture",              // Vulture Monitoring System
	Port{3482, 17}:   "vulture",              // Vulture Monitoring System
	Port{3483, 6}:    "slim-devices",         // Slim Devices Protocol
	Port{3483, 17}:   "slim-devices",         // Slim Devices Protocol
	Port{3484, 6}:    "gbs-stp",              // GBS SnapTalk Protocol
	Port{3484, 17}:   "gbs-stp",              // GBS SnapTalk Protocol
	Port{3485, 6}:    "celatalk",             // Missing description for celatalk
	Port{3485, 17}:   "celatalk",             // CelaTalk
	Port{3486, 6}:    "ifsf-hb-port",         // IFSF Heartbeat Port
	Port{3486, 17}:   "ifsf-hb-port",         // IFSF Heartbeat Port
	Port{3487, 6}:    "ltctcp",               // ltcudp | LISA TCP Transfer Channel | LISA UDP Transfer Channel
	Port{3487, 17}:   "ltcudp",               // LISA UDP Transfer Channel
	Port{3488, 6}:    "fs-rh-srv",            // FS Remote Host Server
	Port{3488, 17}:   "fs-rh-srv",            // FS Remote Host Server
	Port{3489, 6}:    "dtp-dia",              // DTP DIA
	Port{3489, 17}:   "dtp-dia",              // DTP DIA
	Port{3490, 6}:    "colubris",             // Colubris Management Port
	Port{3490, 17}:   "colubris",             // Colubris Management Port
	Port{3491, 6}:    "swr-port",             // SWR Port
	Port{3491, 17}:   "swr-port",             // SWR Port
	Port{3492, 6}:    "tvdumtray-port",       // TVDUM Tray Port
	Port{3492, 17}:   "tvdumtray-port",       // TVDUM Tray Port
	Port{3493, 6}:    "nut",                  // Network UPS Tools
	Port{3493, 17}:   "nut",                  // Network UPS Tools
	Port{3494, 6}:    "ibm3494",              // IBM 3494
	Port{3494, 17}:   "ibm3494",              // IBM 3494
	Port{3495, 6}:    "seclayer-tcp",         // securitylayer over tcp
	Port{3495, 17}:   "seclayer-tcp",         // securitylayer over tcp
	Port{3496, 6}:    "seclayer-tls",         // securitylayer over tls
	Port{3496, 17}:   "seclayer-tls",         // securitylayer over tls
	Port{3497, 6}:    "ipether232port",       // Missing description for ipether232port
	Port{3497, 17}:   "ipether232port",       // ipEther232Port
	Port{3498, 6}:    "dashpas-port",         // DASHPAS user port
	Port{3498, 17}:   "dashpas-port",         // DASHPAS user port
	Port{3499, 6}:    "sccip-media",          // SccIP Media
	Port{3499, 17}:   "sccip-media",          // SccIP Media
	Port{3500, 6}:    "rtmp-port",            // RTMP Port
	Port{3500, 17}:   "rtmp-port",            // RTMP Port
	Port{3501, 6}:    "isoft-p2p",            // Missing description for isoft-p2p
	Port{3501, 17}:   "isoft-p2p",            // iSoft-P2P
	Port{3502, 6}:    "avinstalldisc",        // Avocent Install Discovery
	Port{3502, 17}:   "avinstalldisc",        // Avocent Install Discovery
	Port{3503, 6}:    "lsp-ping",             // MPLS LSP-echo Port
	Port{3503, 17}:   "lsp-ping",             // MPLS LSP-echo Port
	Port{3504, 6}:    "ironstorm",            // IronStorm game server
	Port{3504, 17}:   "ironstorm",            // IronStorm game server
	Port{3505, 6}:    "ccmcomm",              // CCM communications port
	Port{3505, 17}:   "ccmcomm",              // CCM communications port
	Port{3506, 6}:    "apc-3506",             // APC 3506
	Port{3506, 17}:   "apc-3506",             // APC 3506
	Port{3507, 6}:    "nesh-broker",          // Nesh Broker Port
	Port{3507, 17}:   "nesh-broker",          // Nesh Broker Port
	Port{3508, 6}:    "interactionweb",       // Interaction Web
	Port{3508, 17}:   "interactionweb",       // Interaction Web
	Port{3509, 6}:    "vt-ssl",               // Virtual Token SSL Port
	Port{3509, 17}:   "vt-ssl",               // Virtual Token SSL Port
	Port{3510, 6}:    "xss-port",             // XSS Port
	Port{3510, 17}:   "xss-port",             // XSS Port
	Port{3511, 6}:    "webmail-2",            // WebMail 2
	Port{3511, 17}:   "webmail-2",            // WebMail 2
	Port{3512, 6}:    "aztec",                // Aztec Distribution Port
	Port{3512, 17}:   "aztec",                // Aztec Distribution Port
	Port{3513, 6}:    "arcpd",                // Adaptec Remote Protocol
	Port{3513, 17}:   "arcpd",                // Adaptec Remote Protocol
	Port{3514, 6}:    "must-p2p",             // MUST Peer to Peer
	Port{3514, 17}:   "must-p2p",             // MUST Peer to Peer
	Port{3515, 6}:    "must-backplane",       // MUST Backplane
	Port{3515, 17}:   "must-backplane",       // MUST Backplane
	Port{3516, 6}:    "smartcard-port",       // Smartcard Port
	Port{3516, 17}:   "smartcard-port",       // Smartcard Port
	Port{3517, 6}:    "802-11-iapp",          // IEEE 802.11 WLANs WG IAPP
	Port{3517, 17}:   "802-11-iapp",          // IEEE 802.11 WLANs WG IAPP
	Port{3518, 6}:    "artifact-msg",         // Artifact Message Server
	Port{3518, 17}:   "artifact-msg",         // Artifact Message Server
	Port{3519, 6}:    "nvmsgd",               // galileo | Netvion Messenger Port | Netvion Galileo Port
	Port{3519, 17}:   "galileo",              // Netvion Galileo Port
	Port{3520, 6}:    "galileolog",           // Netvion Galileo Log Port
	Port{3520, 17}:   "galileolog",           // Netvion Galileo Log Port
	Port{3521, 6}:    "mc3ss",                // Telequip Labs MC3SS
	Port{3521, 17}:   "mc3ss",                // Telequip Labs MC3SS
	Port{3522, 6}:    "nssocketport",         // DO over NSSocketPort
	Port{3522, 17}:   "nssocketport",         // DO over NSSocketPort
	Port{3523, 6}:    "odeumservlink",        // Odeum Serverlink
	Port{3523, 17}:   "odeumservlink",        // Odeum Serverlink
	Port{3524, 6}:    "ecmport",              // ECM Server port
	Port{3524, 17}:   "ecmport",              // ECM Server port
	Port{3525, 6}:    "eisport",              // EIS Server port
	Port{3525, 17}:   "eisport",              // EIS Server port
	Port{3526, 6}:    "starquiz-port",        // starQuiz Port
	Port{3526, 17}:   "starquiz-port",        // starQuiz Port
	Port{3527, 6}:    "beserver-msg-q",       // VERITAS Backup Exec Server
	Port{3527, 17}:   "beserver-msg-q",       // VERITAS Backup Exec Server
	Port{3528, 6}:    "jboss-iiop",           // JBoss IIOP
	Port{3528, 17}:   "jboss-iiop",           // JBoss IIOP
	Port{3529, 6}:    "jboss-iiop-ssl",       // JBoss IIOP SSL
	Port{3529, 17}:   "jboss-iiop-ssl",       // JBoss IIOP SSL
	Port{3530, 6}:    "gf",                   // Grid Friendly
	Port{3530, 17}:   "gf",                   // Grid Friendly
	Port{3531, 6}:    "peerenabler",          // joltid | P2PNetworking PeerEnabler protocol | Joltid
	Port{3531, 17}:   "peerenabler",          // P2PNetworking PeerEnabler protocol
	Port{3532, 6}:    "raven-rmp",            // Raven Remote Management Control
	Port{3532, 17}:   "raven-rmp",            // Raven Remote Management Control
	Port{3533, 6}:    "raven-rdp",            // Raven Remote Management Data
	Port{3533, 17}:   "raven-rdp",            // Raven Remote Management Data
	Port{3534, 6}:    "urld-port",            // URL Daemon Port
	Port{3534, 17}:   "urld-port",            // URL Daemon Port
	Port{3535, 6}:    "ms-la",                // Missing description for ms-la
	Port{3535, 17}:   "ms-la",                // MS-LA
	Port{3536, 6}:    "snac",                 // Missing description for snac
	Port{3536, 17}:   "snac",                 // SNAC
	Port{3537, 6}:    "ni-visa-remote",       // Remote NI-VISA port
	Port{3537, 17}:   "ni-visa-remote",       // Remote NI-VISA port
	Port{3538, 6}:    "ibm-diradm",           // IBM Directory Server
	Port{3538, 17}:   "ibm-diradm",           // IBM Directory Server
	Port{3539, 6}:    "ibm-diradm-ssl",       // IBM Directory Server SSL
	Port{3539, 17}:   "ibm-diradm-ssl",       // IBM Directory Server SSL
	Port{3540, 6}:    "pnrp-port",            // PNRP User Port
	Port{3540, 17}:   "pnrp-port",            // PNRP User Port
	Port{3541, 6}:    "voispeed-port",        // VoiSpeed Port
	Port{3541, 17}:   "voispeed-port",        // VoiSpeed Port
	Port{3542, 6}:    "hacl-monitor",         // HA cluster monitor
	Port{3542, 17}:   "hacl-monitor",         // HA cluster monitor
	Port{3543, 6}:    "qftest-lookup",        // qftest Lookup Port
	Port{3543, 17}:   "qftest-lookup",        // qftest Lookup Port
	Port{3544, 6}:    "teredo",               // Teredo Port
	Port{3544, 17}:   "teredo",               // Teredo Port
	Port{3545, 6}:    "camac",                // CAMAC equipment
	Port{3545, 17}:   "camac",                // CAMAC equipment
	Port{3547, 6}:    "symantec-sim",         // Symantec SIM
	Port{3547, 17}:   "symantec-sim",         // Symantec SIM
	Port{3548, 6}:    "interworld",           // Missing description for interworld
	Port{3548, 17}:   "interworld",           // Interworld
	Port{3549, 6}:    "tellumat-nms",         // Tellumat MDR NMS
	Port{3549, 17}:   "tellumat-nms",         // Tellumat MDR NMS
	Port{3550, 6}:    "ssmpp",                // Secure SMPP
	Port{3550, 17}:   "ssmpp",                // Secure SMPP
	Port{3551, 6}:    "apcupsd",              // Apcupsd Information Port
	Port{3551, 17}:   "apcupsd",              // Apcupsd Information Port
	Port{3552, 6}:    "taserver",             // TeamAgenda Server Port
	Port{3552, 17}:   "taserver",             // TeamAgenda Server Port
	Port{3553, 6}:    "rbr-discovery",        // Red Box Recorder ADP
	Port{3553, 17}:   "rbr-discovery",        // Red Box Recorder ADP
	Port{3554, 6}:    "questnotify",          // Quest Notification Server
	Port{3554, 17}:   "questnotify",          // Quest Notification Server
	Port{3555, 6}:    "razor",                // Vipul's Razor
	Port{3555, 17}:   "razor",                // Vipul's Razor
	Port{3556, 6}:    "sky-transport",        // Sky Transport Protocol
	Port{3556, 17}:   "sky-transport",        // Sky Transport Protocol
	Port{3557, 6}:    "personalos-001",       // PersonalOS Comm Port
	Port{3557, 17}:   "personalos-001",       // PersonalOS Comm Port
	Port{3558, 6}:    "mcp-port",             // MCP user port
	Port{3558, 17}:   "mcp-port",             // MCP user port
	Port{3559, 6}:    "cctv-port",            // CCTV control port
	Port{3559, 17}:   "cctv-port",            // CCTV control port
	Port{3560, 6}:    "iniserve-port",        // INIServe port
	Port{3560, 17}:   "iniserve-port",        // INIServe port
	Port{3561, 6}:    "bmc-onekey",           // Missing description for bmc-onekey
	Port{3561, 17}:   "bmc-onekey",           // BMC-OneKey
	Port{3562, 6}:    "sdbproxy",             // Missing description for sdbproxy
	Port{3562, 17}:   "sdbproxy",             // SDBProxy
	Port{3563, 6}:    "watcomdebug",          // Watcom Debug
	Port{3563, 17}:   "watcomdebug",          // Watcom Debug
	Port{3564, 6}:    "esimport",             // Electromed SIM port
	Port{3564, 17}:   "esimport",             // Electromed SIM port
	Port{3565, 132}:  "m2pa",                 // Missing description for m2pa
	Port{3565, 6}:    "m2pa",                 // M2PA
	Port{3566, 6}:    "quest-data-hub",       // Quest Data Hub
	Port{3567, 6}:    "oap",                  // dof-eps | Object Access Protocol | DOF Protocol Stack
	Port{3567, 17}:   "oap",                  // Object Access Protocol
	Port{3568, 6}:    "oap-s",                // dof-tunnel-sec | Object Access Protocol over SSL | DOF Secure Tunnel
	Port{3568, 17}:   "oap-s",                // Object Access Protocol over SSL
	Port{3569, 6}:    "mbg-ctrl",             // Meinberg Control Service
	Port{3569, 17}:   "mbg-ctrl",             // Meinberg Control Service
	Port{3570, 6}:    "mccwebsvr-port",       // MCC Web Server Port
	Port{3570, 17}:   "mccwebsvr-port",       // MCC Web Server Port
	Port{3571, 6}:    "megardsvr-port",       // MegaRAID Server Port
	Port{3571, 17}:   "megardsvr-port",       // MegaRAID Server Port
	Port{3572, 6}:    "megaregsvrport",       // Registration Server Port
	Port{3572, 17}:   "megaregsvrport",       // Registration Server Port
	Port{3573, 6}:    "tag-ups-1",            // Advantage Group UPS Suite
	Port{3573, 17}:   "tag-ups-1",            // Advantage Group UPS Suite
	Port{3574, 6}:    "dmaf-server",          // dmaf-caster | DMAF Server | DMAF Caster
	Port{3574, 17}:   "dmaf-caster",          // DMAF Caster
	Port{3575, 6}:    "ccm-port",             // Coalsere CCM Port
	Port{3575, 17}:   "ccm-port",             // Coalsere CCM Port
	Port{3576, 6}:    "cmc-port",             // Coalsere CMC Port
	Port{3576, 17}:   "cmc-port",             // Coalsere CMC Port
	Port{3577, 6}:    "config-port",          // Configuration Port
	Port{3577, 17}:   "config-port",          // Configuration Port
	Port{3578, 6}:    "data-port",            // Data Port
	Port{3578, 17}:   "data-port",            // Data Port
	Port{3579, 6}:    "ttat3lb",              // Tarantella Load Balancing
	Port{3579, 17}:   "ttat3lb",              // Tarantella Load Balancing
	Port{3580, 6}:    "nati-svrloc",          // NATI-ServiceLocator
	Port{3580, 17}:   "nati-svrloc",          // NATI-ServiceLocator
	Port{3581, 6}:    "kfxaclicensing",       // Ascent Capture Licensing
	Port{3581, 17}:   "kfxaclicensing",       // Ascent Capture Licensing
	Port{3582, 6}:    "press",                // PEG PRESS Server
	Port{3582, 17}:   "press",                // PEG PRESS Server
	Port{3583, 6}:    "canex-watch",          // CANEX Watch System
	Port{3583, 17}:   "canex-watch",          // CANEX Watch System
	Port{3584, 6}:    "u-dbap",               // U-DBase Access Protocol
	Port{3584, 17}:   "u-dbap",               // U-DBase Access Protocol
	Port{3585, 6}:    "emprise-lls",          // Emprise License Server
	Port{3585, 17}:   "emprise-lls",          // Emprise License Server
	Port{3586, 6}:    "emprise-lsc",          // License Server Console
	Port{3586, 17}:   "emprise-lsc",          // License Server Console
	Port{3587, 6}:    "p2pgroup",             // Peer to Peer Grouping
	Port{3587, 17}:   "p2pgroup",             // Peer to Peer Grouping
	Port{3588, 6}:    "sentinel",             // Sentinel Server
	Port{3588, 17}:   "sentinel",             // Sentinel Server
	Port{3589, 6}:    "isomair",              // Missing description for isomair
	Port{3589, 17}:   "isomair",              // Missing description for isomair
	Port{3590, 6}:    "wv-csp-sms",           // WV CSP SMS Binding
	Port{3590, 17}:   "wv-csp-sms",           // WV CSP SMS Binding
	Port{3591, 6}:    "gtrack-server",        // LOCANIS G-TRACK Server
	Port{3591, 17}:   "gtrack-server",        // LOCANIS G-TRACK Server
	Port{3592, 6}:    "gtrack-ne",            // LOCANIS G-TRACK NE Port
	Port{3592, 17}:   "gtrack-ne",            // LOCANIS G-TRACK NE Port
	Port{3593, 6}:    "bpmd",                 // BP Model Debugger
	Port{3593, 17}:   "bpmd",                 // BP Model Debugger
	Port{3594, 6}:    "mediaspace",           // Missing description for mediaspace
	Port{3594, 17}:   "mediaspace",           // MediaSpace
	Port{3595, 6}:    "shareapp",             // Missing description for shareapp
	Port{3595, 17}:   "shareapp",             // ShareApp
	Port{3596, 6}:    "iw-mmogame",           // Illusion Wireless MMOG
	Port{3596, 17}:   "iw-mmogame",           // Illusion Wireless MMOG
	Port{3597, 6}:    "a14",                  // A14 (AN-to-SC MM)
	Port{3597, 17}:   "a14",                  // A14 (AN-to-SC MM)
	Port{3598, 6}:    "a15",                  // A15 (AN-to-AN)
	Port{3598, 17}:   "a15",                  // A15 (AN-to-AN)
	Port{3599, 6}:    "quasar-server",        // Quasar Accounting Server
	Port{3599, 17}:   "quasar-server",        // Quasar Accounting Server
	Port{3600, 6}:    "trap-daemon",          // text relay-answer
	Port{3600, 17}:   "trap-daemon",          // text relay-answer
	Port{3601, 6}:    "visinet-gui",          // Visinet Gui
	Port{3601, 17}:   "visinet-gui",          // Visinet Gui
	Port{3602, 6}:    "infiniswitchcl",       // InfiniSwitch Mgr Client
	Port{3602, 17}:   "infiniswitchcl",       // InfiniSwitch Mgr Client
	Port{3603, 6}:    "int-rcv-cntrl",        // Integrated Rcvr Control
	Port{3603, 17}:   "int-rcv-cntrl",        // Integrated Rcvr Control
	Port{3604, 6}:    "bmc-jmx-port",         // BMC JMX Port
	Port{3604, 17}:   "bmc-jmx-port",         // BMC JMX Port
	Port{3605, 6}:    "comcam-io",            // ComCam IO Port
	Port{3605, 17}:   "comcam-io",            // ComCam IO Port
	Port{3606, 6}:    "splitlock",            // Splitlock Server
	Port{3606, 17}:   "splitlock",            // Splitlock Server
	Port{3607, 6}:    "precise-i3",           // Precise I3
	Port{3607, 17}:   "precise-i3",           // Precise I3
	Port{3608, 6}:    "trendchip-dcp",        // Trendchip control protocol
	Port{3608, 17}:   "trendchip-dcp",        // Trendchip control protocol
	Port{3609, 6}:    "cpdi-pidas-cm",        // CPDI PIDAS Connection Mon
	Port{3609, 17}:   "cpdi-pidas-cm",        // CPDI PIDAS Connection Mon
	Port{3610, 6}:    "echonet",              // Missing description for echonet
	Port{3610, 17}:   "echonet",              // ECHONET
	Port{3611, 6}:    "six-degrees",          // Six Degrees Port
	Port{3611, 17}:   "six-degrees",          // Six Degrees Port
	Port{3612, 6}:    "hp-dataprotect",       // HP Data Protector
	Port{3612, 17}:   "hp-dataprotect",       // HP Data Protector
	Port{3613, 6}:    "alaris-disc",          // Alaris Device Discovery
	Port{3613, 17}:   "alaris-disc",          // Alaris Device Discovery
	Port{3614, 6}:    "sigma-port",           // Invensys Sigma Port | Satchwell Sigma
	Port{3614, 17}:   "sigma-port",           // Invensys Sigma Port
	Port{3615, 6}:    "start-network",        // Start Messaging Network
	Port{3615, 17}:   "start-network",        // Start Messaging Network
	Port{3616, 6}:    "cd3o-protocol",        // cd3o Control Protocol
	Port{3616, 17}:   "cd3o-protocol",        // cd3o Control Protocol
	Port{3617, 6}:    "sharp-server",         // ATI SHARP Logic Engine
	Port{3617, 17}:   "sharp-server",         // ATI SHARP Logic Engine
	Port{3618, 6}:    "aairnet-1",            // AAIR-Network 1
	Port{3618, 17}:   "aairnet-1",            // AAIR-Network 1
	Port{3619, 6}:    "aairnet-2",            // AAIR-Network 2
	Port{3619, 17}:   "aairnet-2",            // AAIR-Network 2
	Port{3620, 6}:    "ep-pcp",               // EPSON Projector Control Port
	Port{3620, 17}:   "ep-pcp",               // EPSON Projector Control Port
	Port{3621, 6}:    "ep-nsp",               // EPSON Network Screen Port
	Port{3621, 17}:   "ep-nsp",               // EPSON Network Screen Port
	Port{3622, 6}:    "ff-lr-port",           // FF LAN Redundancy Port
	Port{3622, 17}:   "ff-lr-port",           // FF LAN Redundancy Port
	Port{3623, 6}:    "haipe-discover",       // HAIPIS Dynamic Discovery
	Port{3623, 17}:   "haipe-discover",       // HAIPIS Dynamic Discovery
	Port{3624, 6}:    "dist-upgrade",         // Distributed Upgrade Port
	Port{3624, 17}:   "dist-upgrade",         // Distributed Upgrade Port
	Port{3625, 6}:    "volley",               // Missing description for volley
	Port{3625, 17}:   "volley",               // Volley
	Port{3626, 6}:    "bvcdaemon-port",       // bvControl Daemon
	Port{3626, 17}:   "bvcdaemon-port",       // bvControl Daemon
	Port{3627, 6}:    "jamserverport",        // Jam Server Port
	Port{3627, 17}:   "jamserverport",        // Jam Server Port
	Port{3628, 6}:    "ept-machine",          // EPT Machine Interface
	Port{3628, 17}:   "ept-machine",          // EPT Machine Interface
	Port{3629, 6}:    "escvpnet",             // ESC VP.net
	Port{3629, 17}:   "escvpnet",             // ESC VP.net
	Port{3630, 6}:    "cs-remote-db",         // C&S Remote Database Port
	Port{3630, 17}:   "cs-remote-db",         // C&S Remote Database Port
	Port{3631, 6}:    "cs-services",          // C&S Web Services Port
	Port{3631, 17}:   "cs-services",          // C&S Web Services Port
	Port{3632, 6}:    "distccd",              // distcc | Distributed compiler daemon | distributed compiler
	Port{3632, 17}:   "distcc",               // distributed compiler
	Port{3633, 6}:    "wacp",                 // Wyrnix AIS port
	Port{3633, 17}:   "wacp",                 // Wyrnix AIS port
	Port{3634, 6}:    "hlibmgr",              // hNTSP Library Manager
	Port{3634, 17}:   "hlibmgr",              // hNTSP Library Manager
	Port{3635, 6}:    "sdo",                  // Simple Distributed Objects
	Port{3635, 17}:   "sdo",                  // Simple Distributed Objects
	Port{3636, 6}:    "servistaitsm",         // Missing description for servistaitsm
	Port{3636, 17}:   "servistaitsm",         // SerVistaITSM
	Port{3637, 6}:    "scservp",              // Customer Service Port
	Port{3637, 17}:   "scservp",              // Customer Service Port
	Port{3638, 6}:    "ehp-backup",           // EHP Backup Protocol
	Port{3638, 17}:   "ehp-backup",           // EHP Backup Protocol
	Port{3639, 6}:    "xap-ha",               // Extensible Automation
	Port{3639, 17}:   "xap-ha",               // Extensible Automation
	Port{3640, 6}:    "netplay-port1",        // Netplay Port 1
	Port{3640, 17}:   "netplay-port1",        // Netplay Port 1
	Port{3641, 6}:    "netplay-port2",        // Netplay Port 2
	Port{3641, 17}:   "netplay-port2",        // Netplay Port 2
	Port{3642, 6}:    "juxml-port",           // Juxml Replication port
	Port{3642, 17}:   "juxml-port",           // Juxml Replication port
	Port{3643, 6}:    "audiojuggler",         // Missing description for audiojuggler
	Port{3643, 17}:   "audiojuggler",         // AudioJuggler
	Port{3644, 6}:    "ssowatch",             // Missing description for ssowatch
	Port{3644, 17}:   "ssowatch",             // Missing description for ssowatch
	Port{3645, 6}:    "cyc",                  // Missing description for cyc
	Port{3645, 17}:   "cyc",                  // Cyc
	Port{3646, 6}:    "xss-srv-port",         // XSS Server Port
	Port{3646, 17}:   "xss-srv-port",         // XSS Server Port
	Port{3647, 6}:    "splitlock-gw",         // Splitlock Gateway
	Port{3647, 17}:   "splitlock-gw",         // Splitlock Gateway
	Port{3648, 6}:    "fjcp",                 // Fujitsu Cooperation Port
	Port{3648, 17}:   "fjcp",                 // Fujitsu Cooperation Port
	Port{3649, 6}:    "nmmp",                 // Nishioka Miyuki Msg Protocol
	Port{3649, 17}:   "nmmp",                 // Nishioka Miyuki Msg Protocol
	Port{3650, 6}:    "prismiq-plugin",       // PRISMIQ VOD plug-in
	Port{3650, 17}:   "prismiq-plugin",       // PRISMIQ VOD plug-in
	Port{3651, 6}:    "xrpc-registry",        // XRPC Registry
	Port{3651, 17}:   "xrpc-registry",        // XRPC Registry
	Port{3652, 6}:    "vxcrnbuport",          // VxCR NBU Default Port
	Port{3652, 17}:   "vxcrnbuport",          // VxCR NBU Default Port
	Port{3653, 6}:    "tsp",                  // Tunnel Setup Protocol
	Port{3653, 17}:   "tsp",                  // Tunnel Setup Protocol
	Port{3654, 6}:    "vaprtm",               // VAP RealTime Messenger
	Port{3654, 17}:   "vaprtm",               // VAP RealTime Messenger
	Port{3655, 6}:    "abatemgr",             // ActiveBatch Exec Agent
	Port{3655, 17}:   "abatemgr",             // ActiveBatch Exec Agent
	Port{3656, 6}:    "abatjss",              // ActiveBatch Job Scheduler
	Port{3656, 17}:   "abatjss",              // ActiveBatch Job Scheduler
	Port{3657, 6}:    "immedianet-bcn",       // ImmediaNet Beacon
	Port{3657, 17}:   "immedianet-bcn",       // ImmediaNet Beacon
	Port{3658, 6}:    "ps-ams",               // PlayStation AMS (Secure)
	Port{3658, 17}:   "ps-ams",               // PlayStation AMS (Secure)
	Port{3659, 6}:    "apple-sasl",           // Apple SASL
	Port{3659, 17}:   "apple-sasl",           // Apple SASL
	Port{3660, 6}:    "can-nds-ssl",          // IBM Tivoli Directory Service using SSL
	Port{3660, 17}:   "can-nds-ssl",          // IBM Tivoli Directory Service using SSL
	Port{3661, 6}:    "can-ferret-ssl",       // IBM Tivoli Directory Service using SSL
	Port{3661, 17}:   "can-ferret-ssl",       // IBM Tivoli Directory Service using SSL
	Port{3662, 6}:    "pserver",              // Missing description for pserver
	Port{3662, 17}:   "pserver",              // Missing description for pserver
	Port{3663, 6}:    "dtp",                  // DIRECWAY Tunnel Protocol
	Port{3663, 17}:   "dtp",                  // DIRECWAY Tunnel Protocol
	Port{3664, 6}:    "ups-engine",           // UPS Engine Port
	Port{3664, 17}:   "ups-engine",           // UPS Engine Port
	Port{3665, 6}:    "ent-engine",           // Enterprise Engine Port
	Port{3665, 17}:   "ent-engine",           // Enterprise Engine Port
	Port{3666, 6}:    "eserver-pap",          // IBM eServer PAP | IBM EServer PAP
	Port{3666, 17}:   "eserver-pap",          // IBM EServer PAP
	Port{3667, 6}:    "infoexch",             // IBM Information Exchange
	Port{3667, 17}:   "infoexch",             // IBM Information Exchange
	Port{3668, 6}:    "dell-rm-port",         // Dell Remote Management
	Port{3668, 17}:   "dell-rm-port",         // Dell Remote Management
	Port{3669, 6}:    "casanswmgmt",          // CA SAN Switch Management
	Port{3669, 17}:   "casanswmgmt",          // CA SAN Switch Management
	Port{3670, 6}:    "smile",                // SMILE TCP UDP Interface
	Port{3670, 17}:   "smile",                // SMILE TCP UDP Interface
	Port{3671, 6}:    "efcp",                 // e Field Control (EIBnet)
	Port{3671, 17}:   "efcp",                 // e Field Control (EIBnet)
	Port{3672, 6}:    "lispworks-orb",        // LispWorks ORB
	Port{3672, 17}:   "lispworks-orb",        // LispWorks ORB
	Port{3673, 6}:    "mediavault-gui",       // Openview Media Vault GUI
	Port{3673, 17}:   "mediavault-gui",       // Openview Media Vault GUI
	Port{3674, 6}:    "wininstall-ipc",       // WinINSTALL IPC Port
	Port{3674, 17}:   "wininstall-ipc",       // WinINSTALL IPC Port
	Port{3675, 6}:    "calltrax",             // CallTrax Data Port
	Port{3675, 17}:   "calltrax",             // CallTrax Data Port
	Port{3676, 6}:    "va-pacbase",           // VisualAge Pacbase server
	Port{3676, 17}:   "va-pacbase",           // VisualAge Pacbase server
	Port{3677, 6}:    "roverlog",             // RoverLog IPC
	Port{3677, 17}:   "roverlog",             // RoverLog IPC
	Port{3678, 6}:    "ipr-dglt",             // DataGuardianLT
	Port{3678, 17}:   "ipr-dglt",             // DataGuardianLT
	Port{3679, 6}:    "newton-dock",          // #ERROR:Escale (Newton Dock) | Newton Dock
	Port{3679, 17}:   "newton-dock",          // Newton Dock
	Port{3680, 6}:    "npds-tracker",         // NPDS Tracker
	Port{3680, 17}:   "npds-tracker",         // NPDS Tracker
	Port{3681, 6}:    "bts-x73",              // BTS X73 Port
	Port{3681, 17}:   "bts-x73",              // BTS X73 Port
	Port{3682, 6}:    "cas-mapi",             // EMC SmartPackets-MAPI
	Port{3682, 17}:   "cas-mapi",             // EMC SmartPackets-MAPI
	Port{3683, 6}:    "bmc-ea",               // BMC EDV EA
	Port{3683, 17}:   "bmc-ea",               // BMC EDV EA
	Port{3684, 6}:    "faxstfx-port",         // FAXstfX
	Port{3684, 17}:   "faxstfx-port",         // FAXstfX
	Port{3685, 6}:    "dsx-agent",            // DS Expert Agent
	Port{3685, 17}:   "dsx-agent",            // DS Expert Agent
	Port{3686, 6}:    "tnmpv2",               // Trivial Network Management
	Port{3686, 17}:   "tnmpv2",               // Trivial Network Management
	Port{3687, 6}:    "simple-push",          // Missing description for simple-push
	Port{3687, 17}:   "simple-push",          // Missing description for simple-push
	Port{3688, 6}:    "simple-push-s",        // simple-push Secure
	Port{3688, 17}:   "simple-push-s",        // simple-push Secure
	Port{3689, 6}:    "rendezvous",           // daap | Rendezvous Zeroconf (used by Apple iTunes) | Digital Audio Access Protocol (iTunes)
	Port{3689, 17}:   "daap",                 // Digital Audio Access Protocol
	Port{3690, 6}:    "svn",                  // Subversion
	Port{3690, 17}:   "svn",                  // Subversion
	Port{3691, 6}:    "magaya-network",       // Magaya Network Port
	Port{3691, 17}:   "magaya-network",       // Magaya Network Port
	Port{3692, 6}:    "intelsync",            // Brimstone IntelSync
	Port{3692, 17}:   "intelsync",            // Brimstone IntelSync
	Port{3693, 6}:    "easl",                 // Emergency Automatic Structure Lockdown System
	Port{3695, 6}:    "bmc-data-coll",        // BMC Data Collection
	Port{3695, 17}:   "bmc-data-coll",        // BMC Data Collection
	Port{3696, 6}:    "telnetcpcd",           // Telnet Com Port Control
	Port{3696, 17}:   "telnetcpcd",           // Telnet Com Port Control
	Port{3697, 6}:    "nw-license",           // NavisWorks License System | NavisWorks Licnese System
	Port{3697, 17}:   "nw-license",           // NavisWorks Licnese System
	Port{3698, 6}:    "sagectlpanel",         // Missing description for sagectlpanel
	Port{3698, 17}:   "sagectlpanel",         // SAGECTLPANEL
	Port{3699, 6}:    "kpn-icw",              // Internet Call Waiting
	Port{3699, 17}:   "kpn-icw",              // Internet Call Waiting
	Port{3700, 6}:    "lrs-paging",           // LRS NetPage
	Port{3700, 17}:   "lrs-paging",           // LRS NetPage
	Port{3701, 6}:    "netcelera",            // Missing description for netcelera
	Port{3701, 17}:   "netcelera",            // NetCelera
	Port{3702, 6}:    "ws-discovery",         // Web Service Discovery
	Port{3702, 17}:   "ws-discovery",         // Web Service Discovery
	Port{3703, 6}:    "adobeserver-3",        // Adobe Server 3
	Port{3703, 17}:   "adobeserver-3",        // Adobe Server 3
	Port{3704, 6}:    "adobeserver-4",        // Adobe Server 4
	Port{3704, 17}:   "adobeserver-4",        // Adobe Server 4
	Port{3705, 6}:    "adobeserver-5",        // Adobe Server 5
	Port{3705, 17}:   "adobeserver-5",        // Adobe Server 5
	Port{3706, 6}:    "rt-event",             // Real-Time Event Port
	Port{3706, 17}:   "rt-event",             // Real-Time Event Port
	Port{3707, 6}:    "rt-event-s",           // Real-Time Event Secure Port
	Port{3707, 17}:   "rt-event-s",           // Real-Time Event Secure Port
	Port{3708, 6}:    "sun-as-iiops",         // Sun App Svr - Naming
	Port{3708, 17}:   "sun-as-iiops",         // Sun App Svr - Naming
	Port{3709, 6}:    "ca-idms",              // CA-IDMS Server
	Port{3709, 17}:   "ca-idms",              // CA-IDMS Server
	Port{3710, 6}:    "portgate-auth",        // PortGate Authentication
	Port{3710, 17}:   "portgate-auth",        // PortGate Authentication
	Port{3711, 6}:    "edb-server2",          // EBD Server 2
	Port{3711, 17}:   "edb-server2",          // EBD Server 2
	Port{3712, 6}:    "sentinel-ent",         // Sentinel Enterprise
	Port{3712, 17}:   "sentinel-ent",         // Sentinel Enterprise
	Port{3713, 6}:    "tftps",                // TFTP over TLS
	Port{3713, 17}:   "tftps",                // TFTP over TLS
	Port{3714, 6}:    "delos-dms",            // DELOS Direct Messaging
	Port{3714, 17}:   "delos-dms",            // DELOS Direct Messaging
	Port{3715, 6}:    "anoto-rendezv",        // Anoto Rendezvous Port
	Port{3715, 17}:   "anoto-rendezv",        // Anoto Rendezvous Port
	Port{3716, 6}:    "wv-csp-sms-cir",       // WV CSP SMS CIR Channel
	Port{3716, 17}:   "wv-csp-sms-cir",       // WV CSP SMS CIR Channel
	Port{3717, 6}:    "wv-csp-udp-cir",       // WV CSP UDP IP CIR Channel
	Port{3717, 17}:   "wv-csp-udp-cir",       // WV CSP UDP IP CIR Channel
	Port{3718, 6}:    "opus-services",        // OPUS Server Port
	Port{3718, 17}:   "opus-services",        // OPUS Server Port
	Port{3719, 6}:    "itelserverport",       // iTel Server Port
	Port{3719, 17}:   "itelserverport",       // iTel Server Port
	Port{3720, 6}:    "ufastro-instr",        // UF Astro. Instr. Services
	Port{3720, 17}:   "ufastro-instr",        // UF Astro. Instr. Services
	Port{3721, 6}:    "xsync",                // Missing description for xsync
	Port{3721, 17}:   "xsync",                // Xsync
	Port{3722, 6}:    "xserveraid",           // Xserve RAID
	Port{3722, 17}:   "xserveraid",           // Xserve RAID
	Port{3723, 6}:    "sychrond",             // Sychron Service Daemon
	Port{3723, 17}:   "sychrond",             // Sychron Service Daemon
	Port{3724, 6}:    "blizwow",              // World of Warcraft
	Port{3724, 17}:   "blizwow",              // World of Warcraft
	Port{3725, 6}:    "na-er-tip",            // Netia NA-ER Port
	Port{3725, 17}:   "na-er-tip",            // Netia NA-ER Port
	Port{3726, 6}:    "array-manager",        // Xyratex Array Manager | Xyartex Array Manager
	Port{3726, 17}:   "array-manager",        // Xyartex Array Manager
	Port{3727, 6}:    "e-mdu",                // Ericsson Mobile Data Unit
	Port{3727, 17}:   "e-mdu",                // Ericsson Mobile Data Unit
	Port{3728, 6}:    "e-woa",                // Ericsson Web on Air
	Port{3728, 17}:   "e-woa",                // Ericsson Web on Air
	Port{3729, 6}:    "fksp-audit",           // Fireking Audit Port
	Port{3729, 17}:   "fksp-audit",           // Fireking Audit Port
	Port{3730, 6}:    "client-ctrl",          // Client Control
	Port{3730, 17}:   "client-ctrl",          // Client Control
	Port{3731, 6}:    "smap",                 // Service Manager
	Port{3731, 17}:   "smap",                 // Service Manager
	Port{3732, 6}:    "m-wnn",                // Mobile Wnn
	Port{3732, 17}:   "m-wnn",                // Mobile Wnn
	Port{3733, 6}:    "multip-msg",           // Multipuesto Msg Port
	Port{3733, 17}:   "multip-msg",           // Multipuesto Msg Port
	Port{3734, 6}:    "synel-data",           // Synel Data Collection Port
	Port{3734, 17}:   "synel-data",           // Synel Data Collection Port
	Port{3735, 6}:    "pwdis",                // Password Distribution
	Port{3735, 17}:   "pwdis",                // Password Distribution
	Port{3736, 6}:    "rs-rmi",               // RealSpace RMI
	Port{3736, 17}:   "rs-rmi",               // RealSpace RMI
	Port{3737, 6}:    "xpanel",               // XPanel Daemon
	Port{3738, 6}:    "versatalk",            // versaTalk Server Port
	Port{3738, 17}:   "versatalk",            // versaTalk Server Port
	Port{3739, 6}:    "launchbird-lm",        // Launchbird LicenseManager
	Port{3739, 17}:   "launchbird-lm",        // Launchbird LicenseManager
	Port{3740, 6}:    "heartbeat",            // Heartbeat Protocol
	Port{3740, 17}:   "heartbeat",            // Heartbeat Protocol
	Port{3741, 6}:    "wysdma",               // WysDM Agent
	Port{3741, 17}:   "wysdma",               // WysDM Agent
	Port{3742, 6}:    "cst-port",             // CST - Configuration & Service Tracker
	Port{3742, 17}:   "cst-port",             // CST - Configuration & Service Tracker
	Port{3743, 6}:    "ipcs-command",         // IP Control Systems Ltd.
	Port{3743, 17}:   "ipcs-command",         // IP Control Systems Ltd.
	Port{3744, 6}:    "sasg",                 // Missing description for sasg
	Port{3744, 17}:   "sasg",                 // SASG
	Port{3745, 6}:    "gw-call-port",         // GWRTC Call Port
	Port{3745, 17}:   "gw-call-port",         // GWRTC Call Port
	Port{3746, 6}:    "linktest",             // LXPRO.COM LinkTest
	Port{3746, 17}:   "linktest",             // LXPRO.COM LinkTest
	Port{3747, 6}:    "linktest-s",           // LXPRO.COM LinkTest SSL
	Port{3747, 17}:   "linktest-s",           // LXPRO.COM LinkTest SSL
	Port{3748, 6}:    "webdata",              // Missing description for webdata
	Port{3748, 17}:   "webdata",              // webData
	Port{3749, 6}:    "cimtrak",              // Missing description for cimtrak
	Port{3749, 17}:   "cimtrak",              // CimTrak
	Port{3750, 6}:    "cbos-ip-port",         // CBOS IP ncapsalation port | CBOS IP ncapsalatoin port
	Port{3750, 17}:   "cbos-ip-port",         // CBOS IP ncapsalatoin port
	Port{3751, 6}:    "gprs-cube",            // CommLinx GPRS Cube
	Port{3751, 17}:   "gprs-cube",            // CommLinx GPRS Cube
	Port{3752, 6}:    "vipremoteagent",       // Vigil-IP RemoteAgent
	Port{3752, 17}:   "vipremoteagent",       // Vigil-IP RemoteAgent
	Port{3753, 6}:    "nattyserver",          // NattyServer Port
	Port{3753, 17}:   "nattyserver",          // NattyServer Port
	Port{3754, 6}:    "timestenbroker",       // TimesTen Broker Port
	Port{3754, 17}:   "timestenbroker",       // TimesTen Broker Port
	Port{3755, 6}:    "sas-remote-hlp",       // SAS Remote Help Server
	Port{3755, 17}:   "sas-remote-hlp",       // SAS Remote Help Server
	Port{3756, 6}:    "canon-capt",           // Canon CAPT Port
	Port{3756, 17}:   "canon-capt",           // Canon CAPT Port
	Port{3757, 6}:    "grf-port",             // GRF Server Port
	Port{3757, 17}:   "grf-port",             // GRF Server Port
	Port{3758, 6}:    "apw-registry",         // apw RMI registry
	Port{3758, 17}:   "apw-registry",         // apw RMI registry
	Port{3759, 6}:    "exapt-lmgr",           // Exapt License Manager
	Port{3759, 17}:   "exapt-lmgr",           // Exapt License Manager
	Port{3760, 6}:    "adtempusclient",       // adTempus Client | adTEmpus Client
	Port{3760, 17}:   "adtempusclient",       // adTEmpus Client
	Port{3761, 6}:    "gsakmp",               // gsakmp port
	Port{3761, 17}:   "gsakmp",               // gsakmp port
	Port{3762, 6}:    "gbs-smp",              // GBS SnapMail Protocol
	Port{3762, 17}:   "gbs-smp",              // GBS SnapMail Protocol
	Port{3763, 6}:    "xo-wave",              // XO Wave Control Port
	Port{3763, 17}:   "xo-wave",              // XO Wave Control Port
	Port{3764, 6}:    "mni-prot-rout",        // MNI Protected Routing
	Port{3764, 17}:   "mni-prot-rout",        // MNI Protected Routing
	Port{3765, 6}:    "rtraceroute",          // Remote Traceroute
	Port{3765, 17}:   "rtraceroute",          // Remote Traceroute
	Port{3766, 6}:    "sitewatch-s",          // SSL e-watch sitewatch server
	Port{3767, 6}:    "listmgr-port",         // ListMGR Port
	Port{3767, 17}:   "listmgr-port",         // ListMGR Port
	Port{3768, 6}:    "rblcheckd",            // rblcheckd server daemon
	Port{3768, 17}:   "rblcheckd",            // rblcheckd server daemon
	Port{3769, 6}:    "haipe-otnk",           // HAIPE Network Keying
	Port{3769, 17}:   "haipe-otnk",           // HAIPE Network Keying
	Port{3770, 6}:    "cindycollab",          // Cinderella Collaboration
	Port{3770, 17}:   "cindycollab",          // Cinderella Collaboration
	Port{3771, 6}:    "paging-port",          // RTP Paging Port
	Port{3771, 17}:   "paging-port",          // RTP Paging Port
	Port{3772, 6}:    "ctp",                  // Chantry Tunnel Protocol
	Port{3772, 17}:   "ctp",                  // Chantry Tunnel Protocol
	Port{3773, 6}:    "ctdhercules",          // Missing description for ctdhercules
	Port{3773, 17}:   "ctdhercules",          // Missing description for ctdhercules
	Port{3774, 6}:    "zicom",                // Missing description for zicom
	Port{3774, 17}:   "zicom",                // ZICOM
	Port{3775, 6}:    "ispmmgr",              // ISPM Manager Port
	Port{3775, 17}:   "ispmmgr",              // ISPM Manager Port
	Port{3776, 6}:    "dvcprov-port",         // Device Provisioning Port
	Port{3776, 17}:   "dvcprov-port",         // Device Provisioning Port
	Port{3777, 6}:    "jibe-eb",              // Jibe EdgeBurst
	Port{3777, 17}:   "jibe-eb",              // Jibe EdgeBurst
	Port{3778, 6}:    "c-h-it-port",          // Cutler-Hammer IT Port
	Port{3778, 17}:   "c-h-it-port",          // Cutler-Hammer IT Port
	Port{3779, 6}:    "cognima",              // Cognima Replication
	Port{3779, 17}:   "cognima",              // Cognima Replication
	Port{3780, 6}:    "nnp",                  // Nuzzler Network Protocol
	Port{3780, 17}:   "nnp",                  // Nuzzler Network Protocol
	Port{3781, 6}:    "abcvoice-port",        // ABCvoice server port
	Port{3781, 17}:   "abcvoice-port",        // ABCvoice server port
	Port{3782, 6}:    "iso-tp0s",             // Secure ISO TP0 port
	Port{3782, 17}:   "iso-tp0s",             // Secure ISO TP0 port
	Port{3783, 6}:    "bim-pem",              // Impact Mgr. PEM Gateway
	Port{3783, 17}:   "bim-pem",              // Impact Mgr. PEM Gateway
	Port{3784, 6}:    "bfd-control",          // BFD Control Protocol
	Port{3784, 17}:   "bfd-control",          // BFD Control Protocol
	Port{3785, 6}:    "bfd-echo",             // BFD Echo Protocol
	Port{3785, 17}:   "bfd-echo",             // BFD Echo Protocol
	Port{3786, 6}:    "upstriggervsw",        // VSW Upstrigger port
	Port{3786, 17}:   "upstriggervsw",        // VSW Upstrigger port
	Port{3787, 6}:    "fintrx",               // Missing description for fintrx
	Port{3787, 17}:   "fintrx",               // Fintrx
	Port{3788, 6}:    "isrp-port",            // SPACEWAY Routing port
	Port{3788, 17}:   "isrp-port",            // SPACEWAY Routing port
	Port{3789, 6}:    "remotedeploy",         // RemoteDeploy Administration Port [July 2003]
	Port{3789, 17}:   "remotedeploy",         // RemoteDeploy Administration Port [July 2003]
	Port{3790, 6}:    "quickbooksrds",        // QuickBooks RDS
	Port{3790, 17}:   "quickbooksrds",        // QuickBooks RDS
	Port{3791, 6}:    "tvnetworkvideo",       // TV NetworkVideo Data port
	Port{3791, 17}:   "tvnetworkvideo",       // TV NetworkVideo Data port
	Port{3792, 6}:    "sitewatch",            // e-Watch Corporation SiteWatch
	Port{3792, 17}:   "sitewatch",            // e-Watch Corporation SiteWatch
	Port{3793, 6}:    "dcsoftware",           // DataCore Software
	Port{3793, 17}:   "dcsoftware",           // DataCore Software
	Port{3794, 6}:    "jaus",                 // JAUS Robots
	Port{3794, 17}:   "jaus",                 // JAUS Robots
	Port{3795, 6}:    "myblast",              // myBLAST Mekentosj port
	Port{3795, 17}:   "myblast",              // myBLAST Mekentosj port
	Port{3796, 6}:    "spw-dialer",           // Spaceway Dialer
	Port{3796, 17}:   "spw-dialer",           // Spaceway Dialer
	Port{3797, 6}:    "idps",                 // Missing description for idps
	Port{3797, 17}:   "idps",                 // Missing description for idps
	Port{3798, 6}:    "minilock",             // Missing description for minilock
	Port{3798, 17}:   "minilock",             // Minilock
	Port{3799, 6}:    "radius-dynauth",       // RADIUS Dynamic Authorization
	Port{3799, 17}:   "radius-dynauth",       // RADIUS Dynamic Authorization
	Port{3800, 6}:    "pwgpsi",               // Print Services Interface
	Port{3800, 17}:   "pwgpsi",               // Print Services Interface
	Port{3801, 6}:    "ibm-mgr",              // ibm manager service
	Port{3801, 17}:   "ibm-mgr",              // ibm manager service
	Port{3802, 6}:    "vhd",                  // Missing description for vhd
	Port{3802, 17}:   "vhd",                  // VHD
	Port{3803, 6}:    "soniqsync",            // Missing description for soniqsync
	Port{3803, 17}:   "soniqsync",            // SoniqSync
	Port{3804, 6}:    "iqnet-port",           // Harman IQNet Port
	Port{3804, 17}:   "iqnet-port",           // Harman IQNet Port
	Port{3805, 6}:    "tcpdataserver",        // ThorGuard Server Port
	Port{3805, 17}:   "tcpdataserver",        // ThorGuard Server Port
	Port{3806, 6}:    "wsmlb",                // Remote System Manager
	Port{3806, 17}:   "wsmlb",                // Remote System Manager
	Port{3807, 6}:    "spugna",               // SpuGNA Communication Port
	Port{3807, 17}:   "spugna",               // SpuGNA Communication Port
	Port{3808, 6}:    "sun-as-iiops-ca",      // Sun App Svr-IIOPClntAuth
	Port{3808, 17}:   "sun-as-iiops-ca",      // Sun App Svr-IIOPClntAuth
	Port{3809, 6}:    "apocd",                // Java Desktop System Configuration Agent
	Port{3809, 17}:   "apocd",                // Java Desktop System Configuration Agent
	Port{3810, 6}:    "wlanauth",             // WLAN AS server
	Port{3810, 17}:   "wlanauth",             // WLAN AS server
	Port{3811, 6}:    "amp",                  // Missing description for amp
	Port{3811, 17}:   "amp",                  // AMP
	Port{3812, 6}:    "neto-wol-server",      // netO WOL Server
	Port{3812, 17}:   "neto-wol-server",      // netO WOL Server
	Port{3813, 6}:    "rap-ip",               // Rhapsody Interface Protocol
	Port{3813, 17}:   "rap-ip",               // Rhapsody Interface Protocol
	Port{3814, 6}:    "neto-dcs",             // netO DCS
	Port{3814, 17}:   "neto-dcs",             // netO DCS
	Port{3815, 6}:    "lansurveyorxml",       // LANsurveyor XML
	Port{3815, 17}:   "lansurveyorxml",       // LANsurveyor XML
	Port{3816, 6}:    "sunlps-http",          // Sun Local Patch Server
	Port{3816, 17}:   "sunlps-http",          // Sun Local Patch Server
	Port{3817, 6}:    "tapeware",             // Yosemite Tech Tapeware
	Port{3817, 17}:   "tapeware",             // Yosemite Tech Tapeware
	Port{3818, 6}:    "crinis-hb",            // Crinis Heartbeat
	Port{3818, 17}:   "crinis-hb",            // Crinis Heartbeat
	Port{3819, 6}:    "epl-slp",              // EPL Sequ Layer Protocol
	Port{3819, 17}:   "epl-slp",              // EPL Sequ Layer Protocol
	Port{3820, 6}:    "scp",                  // Siemens AuD SCP
	Port{3820, 17}:   "scp",                  // Siemens AuD SCP
	Port{3821, 6}:    "pmcp",                 // ATSC PMCP Standard
	Port{3821, 17}:   "pmcp",                 // ATSC PMCP Standard
	Port{3822, 6}:    "acp-discovery",        // Compute Pool Discovery
	Port{3822, 17}:   "acp-discovery",        // Compute Pool Discovery
	Port{3823, 6}:    "acp-conduit",          // Compute Pool Conduit
	Port{3823, 17}:   "acp-conduit",          // Compute Pool Conduit
	Port{3824, 6}:    "acp-policy",           // Compute Pool Policy
	Port{3824, 17}:   "acp-policy",           // Compute Pool Policy
	Port{3825, 6}:    "ffserver",             // Antera FlowFusion Process Simulation
	Port{3825, 17}:   "ffserver",             // Antera FlowFusion Process Simulation
	Port{3826, 6}:    "wormux",               // warmux | Wormux server | WarMUX game server
	Port{3826, 17}:   "wormux",               // Wormux server
	Port{3827, 6}:    "netmpi",               // Netadmin Systems MPI service
	Port{3827, 17}:   "netmpi",               // Netadmin Systems MPI service
	Port{3828, 6}:    "neteh",                // Netadmin Systems Event Handler
	Port{3828, 17}:   "neteh",                // Netadmin Systems Event Handler
	Port{3829, 6}:    "neteh-ext",            // Netadmin Systems Event Handler External
	Port{3829, 17}:   "neteh-ext",            // Netadmin Systems Event Handler External
	Port{3830, 6}:    "cernsysmgmtagt",       // Cerner System Management Agent
	Port{3830, 17}:   "cernsysmgmtagt",       // Cerner System Management Agent
	Port{3831, 6}:    "dvapps",               // Docsvault Application Service
	Port{3831, 17}:   "dvapps",               // Docsvault Application Service
	Port{3832, 6}:    "xxnetserver",          // Missing description for xxnetserver
	Port{3832, 17}:   "xxnetserver",          // xxNETserver
	Port{3833, 6}:    "aipn-auth",            // AIPN LS Authentication
	Port{3833, 17}:   "aipn-auth",            // AIPN LS Authentication
	Port{3834, 6}:    "spectardata",          // Spectar Data Stream Service
	Port{3834, 17}:   "spectardata",          // Spectar Data Stream Service
	Port{3835, 6}:    "spectardb",            // Spectar Database Rights Service
	Port{3835, 17}:   "spectardb",            // Spectar Database Rights Service
	Port{3836, 6}:    "markem-dcp",           // MARKEM NEXTGEN DCP
	Port{3836, 17}:   "markem-dcp",           // MARKEM NEXTGEN DCP
	Port{3837, 6}:    "mkm-discovery",        // MARKEM Auto-Discovery
	Port{3837, 17}:   "mkm-discovery",        // MARKEM Auto-Discovery
	Port{3838, 6}:    "sos",                  // Scito Object Server
	Port{3838, 17}:   "sos",                  // Scito Object Server
	Port{3839, 6}:    "amx-rms",              // AMX Resource Management Suite
	Port{3839, 17}:   "amx-rms",              // AMX Resource Management Suite
	Port{3840, 6}:    "flirtmitmir",          // www.FlirtMitMir.de
	Port{3840, 17}:   "flirtmitmir",          // www.FlirtMitMir.de
	Port{3841, 6}:    "zfirm-shiprush3",      // shiprush-db-svr | Z-Firm ShipRush v3 | ShipRush Database Server
	Port{3841, 17}:   "zfirm-shiprush3",      // Z-Firm ShipRush v3
	Port{3842, 6}:    "nhci",                 // NHCI status port
	Port{3842, 17}:   "nhci",                 // NHCI status port
	Port{3843, 6}:    "quest-agent",          // Quest Common Agent
	Port{3843, 17}:   "quest-agent",          // Quest Common Agent
	Port{3844, 6}:    "rnm",                  // Missing description for rnm
	Port{3844, 17}:   "rnm",                  // RNM
	Port{3845, 6}:    "v-one-spp",            // V-ONE Single Port Proxy
	Port{3845, 17}:   "v-one-spp",            // V-ONE Single Port Proxy
	Port{3846, 6}:    "an-pcp",               // Astare Network PCP
	Port{3846, 17}:   "an-pcp",               // Astare Network PCP
	Port{3847, 6}:    "msfw-control",         // MS Firewall Control
	Port{3847, 17}:   "msfw-control",         // MS Firewall Control
	Port{3848, 6}:    "item",                 // IT Environmental Monitor
	Port{3848, 17}:   "item",                 // IT Environmental Monitor
	Port{3849, 6}:    "spw-dnspreload",       // SPACEWAY DNS Preload | SPACEWAY DNS Prelaod
	Port{3849, 17}:   "spw-dnspreload",       // SPACEWAY DNS Prelaod
	Port{3850, 6}:    "qtms-bootstrap",       // QTMS Bootstrap Protocol
	Port{3850, 17}:   "qtms-bootstrap",       // QTMS Bootstrap Protocol
	Port{3851, 6}:    "spectraport",          // SpectraTalk Port
	Port{3851, 17}:   "spectraport",          // SpectraTalk Port
	Port{3852, 6}:    "sse-app-config",       // SSE App Configuration
	Port{3852, 17}:   "sse-app-config",       // SSE App Configuration
	Port{3853, 6}:    "sscan",                // SONY scanning protocol
	Port{3853, 17}:   "sscan",                // SONY scanning protocol
	Port{3854, 6}:    "stryker-com",          // Stryker Comm Port
	Port{3854, 17}:   "stryker-com",          // Stryker Comm Port
	Port{3855, 6}:    "opentrac",             // Missing description for opentrac
	Port{3855, 17}:   "opentrac",             // OpenTRAC
	Port{3856, 6}:    "informer",             // Missing description for informer
	Port{3856, 17}:   "informer",             // INFORMER
	Port{3857, 6}:    "trap-port",            // Trap Port
	Port{3857, 17}:   "trap-port",            // Trap Port
	Port{3858, 6}:    "trap-port-mom",        // Trap Port MOM
	Port{3858, 17}:   "trap-port-mom",        // Trap Port MOM
	Port{3859, 6}:    "nav-port",             // Navini Port
	Port{3859, 17}:   "nav-port",             // Navini Port
	Port{3860, 6}:    "sasp",                 // Server Application State Protocol (SASP)
	Port{3860, 17}:   "sasp",                 // Server Application State Protocol (SASP)
	Port{3861, 6}:    "winshadow-hd",         // winShadow Host Discovery
	Port{3861, 17}:   "winshadow-hd",         // winShadow Host Discovery
	Port{3862, 6}:    "giga-pocket",          // Missing description for giga-pocket
	Port{3862, 17}:   "giga-pocket",          // GIGA-POCKET
	Port{3863, 132}:  "asap-sctp",            // asap-udp | asap-tcp | RSerPool ASAP (SCTP) | asap tcp port | asap udp port | asap sctp
	Port{3863, 6}:    "asap-tcp",             // RSerPool ASAP (TCP)
	Port{3863, 17}:   "asap-tcp",             // RSerPool ASAP (UDP)
	Port{3864, 132}:  "asap-sctp-tls",        // asap-tcp-tls | RSerPool ASAP TLS (SCTP) | asap tls tcp port | asap-sctp tls
	Port{3864, 6}:    "asap-tcp-tls",         // RSerPool ASAP TLS (TCP)
	Port{3865, 6}:    "xpl",                  // xpl automation protocol
	Port{3865, 17}:   "xpl",                  // xpl automation protocol
	Port{3866, 6}:    "dzdaemon",             // Sun SDViz DZDAEMON Port
	Port{3866, 17}:   "dzdaemon",             // Sun SDViz DZDAEMON Port
	Port{3867, 6}:    "dzoglserver",          // Sun SDViz DZOGLSERVER Port
	Port{3867, 17}:   "dzoglserver",          // Sun SDViz DZOGLSERVER Port
	Port{3868, 132}:  "diameter",             // Missing description for diameter
	Port{3868, 6}:    "diameter",             // DIAMETER
	Port{3869, 6}:    "ovsam-mgmt",           // hp OVSAM MgmtServer Disco
	Port{3869, 17}:   "ovsam-mgmt",           // hp OVSAM MgmtServer Disco
	Port{3870, 6}:    "ovsam-d-agent",        // hp OVSAM HostAgent Disco
	Port{3870, 17}:   "ovsam-d-agent",        // hp OVSAM HostAgent Disco
	Port{3871, 6}:    "avocent-adsap",        // Avocent DS Authorization
	Port{3871, 17}:   "avocent-adsap",        // Avocent DS Authorization
	Port{3872, 6}:    "oem-agent",            // OEM Agent
	Port{3872, 17}:   "oem-agent",            // OEM Agent
	Port{3873, 6}:    "fagordnc",             // Missing description for fagordnc
	Port{3873, 17}:   "fagordnc",             // Missing description for fagordnc
	Port{3874, 6}:    "sixxsconfig",          // SixXS Configuration
	Port{3874, 17}:   "sixxsconfig",          // SixXS Configuration
	Port{3875, 6}:    "pnbscada",             // Missing description for pnbscada
	Port{3875, 17}:   "pnbscada",             // PNBSCADA
	Port{3876, 6}:    "dl_agent",             // dl-agent | DirectoryLockdown Agent
	Port{3876, 17}:   "dl_agent",             // DirectoryLockdown Agent
	Port{3877, 6}:    "xmpcr-interface",      // XMPCR Interface Port
	Port{3877, 17}:   "xmpcr-interface",      // XMPCR Interface Port
	Port{3878, 6}:    "fotogcad",             // FotoG CAD interface
	Port{3878, 17}:   "fotogcad",             // FotoG CAD interface
	Port{3879, 6}:    "appss-lm",             // appss license manager
	Port{3879, 17}:   "appss-lm",             // appss license manager
	Port{3880, 6}:    "igrs",                 // Missing description for igrs
	Port{3880, 17}:   "igrs",                 // IGRS
	Port{3881, 6}:    "idac",                 // Data Acquisition and Control
	Port{3881, 17}:   "idac",                 // Data Acquisition and Control
	Port{3882, 6}:    "msdts1",               // DTS Service Port
	Port{3882, 17}:   "msdts1",               // DTS Service Port
	Port{3883, 6}:    "vrpn",                 // VR Peripheral Network
	Port{3883, 17}:   "vrpn",                 // VR Peripheral Network
	Port{3884, 6}:    "softrack-meter",       // SofTrack Metering
	Port{3884, 17}:   "softrack-meter",       // SofTrack Metering
	Port{3885, 6}:    "topflow-ssl",          // TopFlow SSL
	Port{3885, 17}:   "topflow-ssl",          // TopFlow SSL
	Port{3886, 6}:    "nei-management",       // NEI management port
	Port{3886, 17}:   "nei-management",       // NEI management port
	Port{3887, 6}:    "ciphire-data",         // Ciphire Data Transport
	Port{3887, 17}:   "ciphire-data",         // Ciphire Data Transport
	Port{3888, 6}:    "ciphire-serv",         // Ciphire Services
	Port{3888, 17}:   "ciphire-serv",         // Ciphire Services
	Port{3889, 6}:    "dandv-tester",         // D and V Tester Control Port
	Port{3889, 17}:   "dandv-tester",         // D and V Tester Control Port
	Port{3890, 6}:    "ndsconnect",           // Niche Data Server Connect
	Port{3890, 17}:   "ndsconnect",           // Niche Data Server Connect
	Port{3891, 6}:    "rtc-pm-port",          // Oracle RTC-PM port
	Port{3891, 17}:   "rtc-pm-port",          // Oracle RTC-PM port
	Port{3892, 6}:    "pcc-image-port",       // Missing description for pcc-image-port
	Port{3892, 17}:   "pcc-image-port",       // PCC-image-port
	Port{3893, 6}:    "cgi-starapi",          // CGI StarAPI Server
	Port{3893, 17}:   "cgi-starapi",          // CGI StarAPI Server
	Port{3894, 6}:    "syam-agent",           // SyAM Agent Port
	Port{3894, 17}:   "syam-agent",           // SyAM Agent Port
	Port{3895, 6}:    "syam-smc",             // SyAm SMC Service Port
	Port{3895, 17}:   "syam-smc",             // SyAm SMC Service Port
	Port{3896, 6}:    "sdo-tls",              // Simple Distributed Objects over TLS
	Port{3896, 17}:   "sdo-tls",              // Simple Distributed Objects over TLS
	Port{3897, 6}:    "sdo-ssh",              // Simple Distributed Objects over SSH
	Port{3897, 17}:   "sdo-ssh",              // Simple Distributed Objects over SSH
	Port{3898, 6}:    "senip",                // IAS, Inc. SmartEye NET Internet Protocol
	Port{3898, 17}:   "senip",                // IAS, Inc. SmartEye NET Internet Protocol
	Port{3899, 6}:    "itv-control",          // ITV Port
	Port{3899, 17}:   "itv-control",          // ITV Port
	Port{3900, 6}:    "udt_os",               // udt-os | Unidata UDT OS
	Port{3900, 17}:   "udt_os",               // Unidata UDT OS
	Port{3901, 6}:    "nimsh",                // NIM Service Handler
	Port{3901, 17}:   "nimsh",                // NIM Service Handler
	Port{3902, 6}:    "nimaux",               // NIMsh Auxiliary Port
	Port{3902, 17}:   "nimaux",               // NIMsh Auxiliary Port
	Port{3903, 6}:    "charsetmgr",           // Missing description for charsetmgr
	Port{3903, 17}:   "charsetmgr",           // CharsetMGR
	Port{3904, 6}:    "omnilink-port",        // Arnet Omnilink Port
	Port{3904, 17}:   "omnilink-port",        // Arnet Omnilink Port
	Port{3905, 6}:    "mupdate",              // Mailbox Update (MUPDATE) protocol
	Port{3905, 17}:   "mupdate",              // Mailbox Update (MUPDATE) protocol
	Port{3906, 6}:    "topovista-data",       // TopoVista elevation data
	Port{3906, 17}:   "topovista-data",       // TopoVista elevation data
	Port{3907, 6}:    "imoguia-port",         // Imoguia Port
	Port{3907, 17}:   "imoguia-port",         // Imoguia Port
	Port{3908, 6}:    "hppronetman",          // HP Procurve NetManagement
	Port{3908, 17}:   "hppronetman",          // HP Procurve NetManagement
	Port{3909, 6}:    "surfcontrolcpa",       // SurfControl CPA
	Port{3909, 17}:   "surfcontrolcpa",       // SurfControl CPA
	Port{3910, 6}:    "prnrequest",           // Printer Request Port
	Port{3910, 17}:   "prnrequest",           // Printer Request Port
	Port{3911, 6}:    "prnstatus",            // Printer Status Port
	Port{3911, 17}:   "prnstatus",            // Printer Status Port
	Port{3912, 6}:    "gbmt-stars",           // Global Maintech Stars
	Port{3912, 17}:   "gbmt-stars",           // Global Maintech Stars
	Port{3913, 6}:    "listcrt-port",         // ListCREATOR Port
	Port{3913, 17}:   "listcrt-port",         // ListCREATOR Port
	Port{3914, 6}:    "listcrt-port-2",       // ListCREATOR Port 2
	Port{3914, 17}:   "listcrt-port-2",       // ListCREATOR Port 2
	Port{3915, 6}:    "agcat",                // Auto-Graphics Cataloging
	Port{3915, 17}:   "agcat",                // Auto-Graphics Cataloging
	Port{3916, 6}:    "wysdmc",               // WysDM Controller
	Port{3916, 17}:   "wysdmc",               // WysDM Controller
	Port{3917, 6}:    "aftmux",               // AFT multiplex port | AFT multiples port
	Port{3917, 17}:   "aftmux",               // AFT multiples port
	Port{3918, 6}:    "pktcablemmcops",       // PacketCableMultimediaCOPS
	Port{3918, 17}:   "pktcablemmcops",       // PacketCableMultimediaCOPS
	Port{3919, 6}:    "hyperip",              // Missing description for hyperip
	Port{3919, 17}:   "hyperip",              // HyperIP
	Port{3920, 6}:    "exasoftport1",         // Exasoft IP Port
	Port{3920, 17}:   "exasoftport1",         // Exasoft IP Port
	Port{3921, 6}:    "herodotus-net",        // Herodotus Net
	Port{3921, 17}:   "herodotus-net",        // Herodotus Net
	Port{3922, 6}:    "sor-update",           // Soronti Update Port
	Port{3922, 17}:   "sor-update",           // Soronti Update Port
	Port{3923, 6}:    "symb-sb-port",         // Symbian Service Broker
	Port{3923, 17}:   "symb-sb-port",         // Symbian Service Broker
	Port{3924, 6}:    "mpl-gprs-port",        // MPL_GPRS_PORT | MPL_GPRS_Port
	Port{3924, 17}:   "mpl-gprs-port",        // MPL_GPRS_Port
	Port{3925, 6}:    "zmp",                  // Zoran Media Port
	Port{3925, 17}:   "zmp",                  // Zoran Media Port
	Port{3926, 6}:    "winport",              // Missing description for winport
	Port{3926, 17}:   "winport",              // WINPort
	Port{3927, 6}:    "natdataservice",       // ScsTsr
	Port{3927, 17}:   "natdataservice",       // ScsTsr
	Port{3928, 6}:    "netboot-pxe",          // PXE NetBoot Manager
	Port{3928, 17}:   "netboot-pxe",          // PXE NetBoot Manager
	Port{3929, 6}:    "smauth-port",          // AMS Port
	Port{3929, 17}:   "smauth-port",          // AMS Port
	Port{3930, 6}:    "syam-webserver",       // Syam Web Server Port
	Port{3930, 17}:   "syam-webserver",       // Syam Web Server Port
	Port{3931, 6}:    "msr-plugin-port",      // MSR Plugin Port
	Port{3931, 17}:   "msr-plugin-port",      // MSR Plugin Port
	Port{3932, 6}:    "dyn-site",             // Dynamic Site System
	Port{3932, 17}:   "dyn-site",             // Dynamic Site System
	Port{3933, 6}:    "plbserve-port",        // PL B App Server User Port
	Port{3933, 17}:   "plbserve-port",        // PL B App Server User Port
	Port{3934, 6}:    "sunfm-port",           // PL B File Manager Port
	Port{3934, 17}:   "sunfm-port",           // PL B File Manager Port
	Port{3935, 6}:    "sdp-portmapper",       // SDP Port Mapper Protocol
	Port{3935, 17}:   "sdp-portmapper",       // SDP Port Mapper Protocol
	Port{3936, 6}:    "mailprox",             // Missing description for mailprox
	Port{3936, 17}:   "mailprox",             // Mailprox
	Port{3937, 6}:    "dvbservdsc",           // DVB Service Discovery
	Port{3937, 17}:   "dvbservdsc",           // DVB Service Discovery
	Port{3938, 6}:    "dbcontrol_agent",      // dbcontrol-agent | Oracle dbControl Agent po | Oracel dbControl Agent po
	Port{3938, 17}:   "dbcontrol_agent",      // Oracel dbControl Agent po
	Port{3939, 6}:    "aamp",                 // Anti-virus Application Management Port
	Port{3939, 17}:   "aamp",                 // Anti-virus Application Management Port
	Port{3940, 6}:    "xecp-node",            // XeCP Node Service
	Port{3940, 17}:   "xecp-node",            // XeCP Node Service
	Port{3941, 6}:    "homeportal-web",       // Home Portal Web Server
	Port{3941, 17}:   "homeportal-web",       // Home Portal Web Server
	Port{3942, 6}:    "srdp",                 // satellite distribution
	Port{3942, 17}:   "srdp",                 // satellite distribution
	Port{3943, 6}:    "tig",                  // TetraNode Ip Gateway
	Port{3943, 17}:   "tig",                  // TetraNode Ip Gateway
	Port{3944, 6}:    "sops",                 // S-Ops Management
	Port{3944, 17}:   "sops",                 // S-Ops Management
	Port{3945, 6}:    "emcads",               // EMCADS Server Port
	Port{3945, 17}:   "emcads",               // EMCADS Server Port
	Port{3946, 6}:    "backupedge",           // BackupEDGE Server
	Port{3946, 17}:   "backupedge",           // BackupEDGE Server
	Port{3947, 6}:    "ccp",                  // Connect and Control Protocol for Consumer, Commercial, and Industrial Electronic Devices
	Port{3947, 17}:   "ccp",                  // Connect and Control Protocol for Consumer, Commercial, and Industrial Electronic Devices
	Port{3948, 6}:    "apdap",                // Anton Paar Device Administration Protocol
	Port{3948, 17}:   "apdap",                // Anton Paar Device Administration Protocol
	Port{3949, 6}:    "drip",                 // Dynamic Routing Information Protocol
	Port{3949, 17}:   "drip",                 // Dynamic Routing Information Protocol
	Port{3950, 6}:    "namemunge",            // Name Munging
	Port{3950, 17}:   "namemunge",            // Name Munging
	Port{3951, 6}:    "pwgippfax",            // PWG IPP Facsimile
	Port{3951, 17}:   "pwgippfax",            // PWG IPP Facsimile
	Port{3952, 6}:    "i3-sessionmgr",        // I3 Session Manager
	Port{3952, 17}:   "i3-sessionmgr",        // I3 Session Manager
	Port{3953, 6}:    "xmlink-connect",       // Eydeas XMLink Connect
	Port{3953, 17}:   "xmlink-connect",       // Eydeas XMLink Connect
	Port{3954, 6}:    "adrep",                // AD Replication RPC
	Port{3954, 17}:   "adrep",                // AD Replication RPC
	Port{3955, 6}:    "p2pcommunity",         // Missing description for p2pcommunity
	Port{3955, 17}:   "p2pcommunity",         // p2pCommunity
	Port{3956, 6}:    "gvcp",                 // GigE Vision Control
	Port{3956, 17}:   "gvcp",                 // GigE Vision Control
	Port{3957, 6}:    "mqe-broker",           // MQEnterprise Broker
	Port{3957, 17}:   "mqe-broker",           // MQEnterprise Broker
	Port{3958, 6}:    "mqe-agent",            // MQEnterprise Agent
	Port{3958, 17}:   "mqe-agent",            // MQEnterprise Agent
	Port{3959, 6}:    "treehopper",           // Tree Hopper Networking
	Port{3959, 17}:   "treehopper",           // Tree Hopper Networking
	Port{3960, 6}:    "bess",                 // Bess Peer Assessment
	Port{3960, 17}:   "bess",                 // Bess Peer Assessment
	Port{3961, 6}:    "proaxess",             // ProAxess Server
	Port{3961, 17}:   "proaxess",             // ProAxess Server
	Port{3962, 6}:    "sbi-agent",            // SBI Agent Protocol
	Port{3962, 17}:   "sbi-agent",            // SBI Agent Protocol
	Port{3963, 6}:    "thrp",                 // Teran Hybrid Routing Protocol
	Port{3963, 17}:   "thrp",                 // Teran Hybrid Routing Protocol
	Port{3964, 6}:    "sasggprs",             // SASG GPRS
	Port{3964, 17}:   "sasggprs",             // SASG GPRS
	Port{3965, 6}:    "ati-ip-to-ncpe",       // Avanti IP to NCPE API
	Port{3965, 17}:   "ati-ip-to-ncpe",       // Avanti IP to NCPE API
	Port{3966, 6}:    "bflckmgr",             // BuildForge Lock Manager
	Port{3966, 17}:   "bflckmgr",             // BuildForge Lock Manager
	Port{3967, 6}:    "ppsms",                // PPS Message Service
	Port{3967, 17}:   "ppsms",                // PPS Message Service
	Port{3968, 6}:    "ianywhere-dbns",       // iAnywhere DBNS
	Port{3968, 17}:   "ianywhere-dbns",       // iAnywhere DBNS
	Port{3969, 6}:    "landmarks",            // Landmark Messages
	Port{3969, 17}:   "landmarks",            // Landmark Messages
	Port{3970, 6}:    "lanrevagent",          // LANrev Agent
	Port{3970, 17}:   "lanrevagent",          // LANrev Agent
	Port{3971, 6}:    "lanrevserver",         // LANrev Server
	Port{3971, 17}:   "lanrevserver",         // LANrev Server
	Port{3972, 6}:    "iconp",                // ict-control Protocol
	Port{3972, 17}:   "iconp",                // ict-control Protocol
	Port{3973, 6}:    "progistics",           // ConnectShip Progistics
	Port{3973, 17}:   "progistics",           // ConnectShip Progistics
	Port{3974, 6}:    "citysearch",           // Remote Applicant Tracking Service
	Port{3974, 17}:   "citysearch",           // Remote Applicant Tracking Service
	Port{3975, 6}:    "airshot",              // Air Shot
	Port{3975, 17}:   "airshot",              // Air Shot
	Port{3976, 6}:    "opswagent",            // Opsware Agent | Server Automation Agent
	Port{3976, 17}:   "opswagent",            // Opsware Agent
	Port{3977, 6}:    "opswmanager",          // Opsware Manager
	Port{3977, 17}:   "opswmanager",          // Opsware Manager
	Port{3978, 6}:    "secure-cfg-svr",       // Secured Configuration Server
	Port{3978, 17}:   "secure-cfg-svr",       // Secured Configuration Server
	Port{3979, 6}:    "smwan",                // Smith Micro Wide Area Network Service
	Port{3979, 17}:   "smwan",                // Smith Micro Wide Area Network Service
	Port{3980, 6}:    "acms",                 // Aircraft Cabin Management System
	Port{3980, 17}:   "acms",                 // Aircraft Cabin Management System
	Port{3981, 6}:    "starfish",             // Starfish System Admin
	Port{3981, 17}:   "starfish",             // Starfish System Admin
	Port{3982, 6}:    "eis",                  // ESRI Image Server
	Port{3982, 17}:   "eis",                  // ESRI Image Server
	Port{3983, 6}:    "eisp",                 // ESRI Image Service
	Port{3983, 17}:   "eisp",                 // ESRI Image Service
	Port{3984, 6}:    "mapper-nodemgr",       // MAPPER network node manager
	Port{3984, 17}:   "mapper-nodemgr",       // MAPPER network node manager
	Port{3985, 6}:    "mapper-mapethd",       // MAPPER TCP IP server
	Port{3985, 17}:   "mapper-mapethd",       // MAPPER TCP IP server
	Port{3986, 6}:    "mapper-ws_ethd",       // mapper-ws-ethd | MAPPER workstation server
	Port{3986, 17}:   "mapper-ws_ethd",       // MAPPER workstation server
	Port{3987, 6}:    "centerline",           // Missing description for centerline
	Port{3987, 17}:   "centerline",           // Centerline
	Port{3988, 6}:    "dcs-config",           // DCS Configuration Port
	Port{3988, 17}:   "dcs-config",           // DCS Configuration Port
	Port{3989, 6}:    "bv-queryengine",       // BindView-Query Engine
	Port{3989, 17}:   "bv-queryengine",       // BindView-Query Engine
	Port{3990, 6}:    "bv-is",                // BindView-IS
	Port{3990, 17}:   "bv-is",                // BindView-IS
	Port{3991, 6}:    "bv-smcsrv",            // BindView-SMCServer
	Port{3991, 17}:   "bv-smcsrv",            // BindView-SMCServer
	Port{3992, 6}:    "bv-ds",                // BindView-DirectoryServer
	Port{3992, 17}:   "bv-ds",                // BindView-DirectoryServer
	Port{3993, 6}:    "bv-agent",             // BindView-Agent
	Port{3993, 17}:   "bv-agent",             // BindView-Agent
	Port{3995, 6}:    "iss-mgmt-ssl",         // ISS Management Svcs SSL
	Port{3995, 17}:   "iss-mgmt-ssl",         // ISS Management Svcs SSL
	Port{3996, 6}:    "abcsoftware",          // abcsoftware-01
	Port{3996, 17}:   "remoteanything",       // neoworx remote-anything slave daemon
	Port{3997, 6}:    "agentsease-db",        // aes_db
	Port{3997, 17}:   "remoteanything",       // neoworx remote-anything master daemon
	Port{3998, 6}:    "dnx",                  // Distributed Nagios Executor Service
	Port{3998, 17}:   "remoteanything",       // neoworx remote-anything reserved
	Port{3999, 6}:    "remoteanything",       // nvcnet | neoworx remote-anything slave file browser | Norman distributes scanning service
	Port{3999, 17}:   "nvcnet",               // Norman distributes scanning service
	Port{4000, 6}:    "remoteanything",       // terabase | neoworx remote-anything slave remote control | Terabase
	Port{4000, 17}:   "icq",                  // AOL ICQ instant messaging clent-server communication
	Port{4001, 6}:    "newoak",               // Missing description for newoak
	Port{4001, 17}:   "newoak",               // NewOak
	Port{4002, 6}:    "mlchat-proxy",         // mlnet - MLChat P2P chat proxy | pxc-spvr-ft
	Port{4002, 17}:   "pxc-spvr-ft",          // Missing description for pxc-spvr-ft
	Port{4003, 6}:    "pxc-splr-ft",          // Missing description for pxc-splr-ft
	Port{4003, 17}:   "pxc-splr-ft",          // Missing description for pxc-splr-ft
	Port{4004, 6}:    "pxc-roid",             // Missing description for pxc-roid
	Port{4004, 17}:   "pxc-roid",             // Missing description for pxc-roid
	Port{4005, 6}:    "pxc-pin",              // Missing description for pxc-pin
	Port{4005, 17}:   "pxc-pin",              // Missing description for pxc-pin
	Port{4006, 6}:    "pxc-spvr",             // Missing description for pxc-spvr
	Port{4006, 17}:   "pxc-spvr",             // Missing description for pxc-spvr
	Port{4007, 6}:    "pxc-splr",             // Missing description for pxc-splr
	Port{4007, 17}:   "pxc-splr",             // Missing description for pxc-splr
	Port{4008, 6}:    "netcheque",            // NetCheque accounting
	Port{4008, 17}:   "netcheque",            // NetCheque accounting
	Port{4009, 6}:    "chimera-hwm",          // Chimera HWM
	Port{4009, 17}:   "chimera-hwm",          // Chimera HWM
	Port{4010, 6}:    "samsung-unidex",       // Samsung Unidex
	Port{4010, 17}:   "samsung-unidex",       // Samsung Unidex
	Port{4011, 6}:    "altserviceboot",       // Alternate Service Boot
	Port{4011, 17}:   "altserviceboot",       // Alternate Service Boot
	Port{4012, 6}:    "pda-gate",             // PDA Gate
	Port{4012, 17}:   "pda-gate",             // PDA Gate
	Port{4013, 6}:    "acl-manager",          // ACL Manager
	Port{4013, 17}:   "acl-manager",          // ACL Manager
	Port{4014, 6}:    "taiclock",             // Missing description for taiclock
	Port{4014, 17}:   "taiclock",             // TAICLOCK
	Port{4015, 6}:    "talarian-mcast1",      // Talarian Mcast
	Port{4015, 17}:   "talarian-mcast1",      // Talarian Mcast
	Port{4016, 6}:    "talarian-mcast2",      // Talarian Mcast
	Port{4016, 17}:   "talarian-mcast2",      // Talarian Mcast
	Port{4017, 6}:    "talarian-mcast3",      // Talarian Mcast
	Port{4017, 17}:   "talarian-mcast3",      // Talarian Mcast
	Port{4018, 6}:    "talarian-mcast4",      // Talarian Mcast
	Port{4018, 17}:   "talarian-mcast4",      // Talarian Mcast
	Port{4019, 6}:    "talarian-mcast5",      // Talarian Mcast
	Port{4019, 17}:   "talarian-mcast5",      // Talarian Mcast
	Port{4020, 6}:    "trap",                 // TRAP Port
	Port{4020, 17}:   "trap",                 // TRAP Port
	Port{4021, 6}:    "nexus-portal",         // Nexus Portal
	Port{4021, 17}:   "nexus-portal",         // Nexus Portal
	Port{4022, 6}:    "dnox",                 // Missing description for dnox
	Port{4022, 17}:   "dnox",                 // DNOX
	Port{4023, 6}:    "esnm-zoning",          // ESNM Zoning Port
	Port{4023, 17}:   "esnm-zoning",          // ESNM Zoning Port
	Port{4024, 6}:    "tnp1-port",            // TNP1 User Port
	Port{4024, 17}:   "tnp1-port",            // TNP1 User Port
	Port{4025, 6}:    "partimage",            // Partition Image Port
	Port{4025, 17}:   "partimage",            // Partition Image Port
	Port{4026, 6}:    "as-debug",             // Graphical Debug Server
	Port{4026, 17}:   "as-debug",             // Graphical Debug Server
	Port{4027, 6}:    "bxp",                  // bitxpress
	Port{4027, 17}:   "bxp",                  // bitxpress
	Port{4028, 6}:    "dtserver-port",        // DTServer Port
	Port{4028, 17}:   "dtserver-port",        // DTServer Port
	Port{4029, 6}:    "ip-qsig",              // IP Q signaling protocol
	Port{4029, 17}:   "ip-qsig",              // IP Q signaling protocol
	Port{4030, 6}:    "jdmn-port",            // Accell JSP Daemon Port
	Port{4030, 17}:   "jdmn-port",            // Accell JSP Daemon Port
	Port{4031, 6}:    "suucp",                // UUCP over SSL
	Port{4031, 17}:   "suucp",                // UUCP over SSL
	Port{4032, 6}:    "vrts-auth-port",       // VERITAS Authorization Service
	Port{4032, 17}:   "vrts-auth-port",       // VERITAS Authorization Service
	Port{4033, 6}:    "sanavigator",          // SANavigator Peer Port
	Port{4033, 17}:   "sanavigator",          // SANavigator Peer Port
	Port{4034, 6}:    "ubxd",                 // Ubiquinox Daemon
	Port{4034, 17}:   "ubxd",                 // Ubiquinox Daemon
	Port{4035, 6}:    "wap-push-http",        // WAP Push OTA-HTTP port
	Port{4035, 17}:   "wap-push-http",        // WAP Push OTA-HTTP port
	Port{4036, 6}:    "wap-push-https",       // WAP Push OTA-HTTP secure
	Port{4036, 17}:   "wap-push-https",       // WAP Push OTA-HTTP secure
	Port{4037, 6}:    "ravehd",               // RaveHD network control
	Port{4037, 17}:   "ravehd",               // RaveHD network control
	Port{4038, 6}:    "fazzt-ptp",            // Fazzt Point-To-Point
	Port{4038, 17}:   "fazzt-ptp",            // Fazzt Point-To-Point
	Port{4039, 6}:    "fazzt-admin",          // Fazzt Administration
	Port{4039, 17}:   "fazzt-admin",          // Fazzt Administration
	Port{4040, 6}:    "yo-main",              // Yo.net main service
	Port{4040, 17}:   "yo-main",              // Yo.net main service
	Port{4041, 6}:    "houston",              // Rocketeer-Houston
	Port{4041, 17}:   "houston",              // Rocketeer-Houston
	Port{4042, 6}:    "ldxp",                 // Missing description for ldxp
	Port{4042, 17}:   "ldxp",                 // LDXP
	Port{4043, 6}:    "nirp",                 // Neighbour Identity Resolution
	Port{4043, 17}:   "nirp",                 // Neighbour Identity Resolution
	Port{4044, 6}:    "ltp",                  // Location Tracking Protocol
	Port{4044, 17}:   "ltp",                  // Location Tracking Protocol
	Port{4045, 6}:    "lockd",                // npp | Network Paging Protocol
	Port{4045, 17}:   "lockd",                // NFS lock daemon manager
	Port{4046, 6}:    "acp-proto",            // Accounting Protocol
	Port{4046, 17}:   "acp-proto",            // Accounting Protocol
	Port{4047, 6}:    "ctp-state",            // Context Transfer Protocol
	Port{4047, 17}:   "ctp-state",            // Context Transfer Protocol
	Port{4049, 6}:    "wafs",                 // Wide Area File Services
	Port{4049, 17}:   "wafs",                 // Wide Area File Services
	Port{4050, 6}:    "cisco-wafs",           // Wide Area File Services
	Port{4050, 17}:   "cisco-wafs",           // Wide Area File Services
	Port{4051, 6}:    "cppdp",                // Cisco Peer to Peer Distribution Protocol
	Port{4051, 17}:   "cppdp",                // Cisco Peer to Peer Distribution Protocol
	Port{4052, 6}:    "interact",             // VoiceConnect Interact
	Port{4052, 17}:   "interact",             // VoiceConnect Interact
	Port{4053, 6}:    "ccu-comm-1",           // CosmoCall Universe Communications Port 1
	Port{4053, 17}:   "ccu-comm-1",           // CosmoCall Universe Communications Port 1
	Port{4054, 6}:    "ccu-comm-2",           // CosmoCall Universe Communications Port 2
	Port{4054, 17}:   "ccu-comm-2",           // CosmoCall Universe Communications Port 2
	Port{4055, 6}:    "ccu-comm-3",           // CosmoCall Universe Communications Port 3
	Port{4055, 17}:   "ccu-comm-3",           // CosmoCall Universe Communications Port 3
	Port{4056, 6}:    "lms",                  // Location Message Service
	Port{4056, 17}:   "lms",                  // Location Message Service
	Port{4057, 6}:    "wfm",                  // Servigistics WFM server
	Port{4057, 17}:   "wfm",                  // Servigistics WFM server
	Port{4058, 6}:    "kingfisher",           // Kingfisher protocol
	Port{4058, 17}:   "kingfisher",           // Kingfisher protocol
	Port{4059, 6}:    "dlms-cosem",           // DLMS COSEM
	Port{4059, 17}:   "dlms-cosem",           // DLMS COSEM
	Port{4060, 6}:    "dsmeter_iatc",         // dsmeter-iatc | DSMETER Inter-Agent Transfer Channel
	Port{4060, 17}:   "dsmeter_iatc",         // DSMETER Inter-Agent Transfer Channel
	Port{4061, 6}:    "ice-location",         // Ice Location Service (TCP)
	Port{4061, 17}:   "ice-location",         // Ice Location Service (TCP)
	Port{4062, 6}:    "ice-slocation",        // Ice Location Service (SSL)
	Port{4062, 17}:   "ice-slocation",        // Ice Location Service (SSL)
	Port{4063, 6}:    "ice-router",           // Ice Firewall Traversal Service (TCP)
	Port{4063, 17}:   "ice-router",           // Ice Firewall Traversal Service (TCP)
	Port{4064, 6}:    "ice-srouter",          // Ice Firewall Traversal Service (SSL)
	Port{4064, 17}:   "ice-srouter",          // Ice Firewall Traversal Service (SSL)
	Port{4065, 6}:    "avanti_cdp",           // avanti-cdp | Avanti Common Data
	Port{4065, 17}:   "avanti_cdp",           // Avanti Common Data
	Port{4066, 6}:    "pmas",                 // Performance Measurement and Analysis
	Port{4066, 17}:   "pmas",                 // Performance Measurement and Analysis
	Port{4067, 6}:    "idp",                  // Information Distribution Protocol
	Port{4067, 17}:   "idp",                  // Information Distribution Protocol
	Port{4068, 6}:    "ipfltbcst",            // IP Fleet Broadcast
	Port{4068, 17}:   "ipfltbcst",            // IP Fleet Broadcast
	Port{4069, 6}:    "minger",               // Minger Email Address Validation Service
	Port{4069, 17}:   "minger",               // Minger Email Address Validation Service
	Port{4070, 6}:    "tripe",                // Trivial IP Encryption (TrIPE)
	Port{4070, 17}:   "tripe",                // Trivial IP Encryption (TrIPE)
	Port{4071, 6}:    "aibkup",               // Automatically Incremental Backup
	Port{4071, 17}:   "aibkup",               // Automatically Incremental Backup
	Port{4072, 6}:    "zieto-sock",           // Zieto Socket Communications
	Port{4072, 17}:   "zieto-sock",           // Zieto Socket Communications
	Port{4073, 6}:    "iRAPP",                // iRAPP Server Protocol | Interactive Remote Application Pairing Protocol
	Port{4073, 17}:   "iRAPP",                // iRAPP Server Protocol
	Port{4074, 6}:    "cequint-cityid",       // Cequint City ID UI trigger
	Port{4074, 17}:   "cequint-cityid",       // Cequint City ID UI trigger
	Port{4075, 6}:    "perimlan",             // ISC Alarm Message Service
	Port{4075, 17}:   "perimlan",             // ISC Alarm Message Service
	Port{4076, 6}:    "seraph",               // Seraph DCS
	Port{4076, 17}:   "seraph",               // Seraph DCS
	Port{4077, 6}:    "ascomalarm",           // Ascom IP Alarming
	Port{4077, 17}:   "ascomalarm",           // Ascom IP Alarming
	Port{4078, 6}:    "cssp",                 // Coordinated Security Service Protocol
	Port{4079, 6}:    "santools",             // SANtools Diagnostic Server
	Port{4079, 17}:   "santools",             // SANtools Diagnostic Server
	Port{4080, 6}:    "lorica-in",            // Lorica inside facing
	Port{4080, 17}:   "lorica-in",            // Lorica inside facing
	Port{4081, 6}:    "lorica-in-sec",        // Lorica inside facing (SSL)
	Port{4081, 17}:   "lorica-in-sec",        // Lorica inside facing (SSL)
	Port{4082, 6}:    "lorica-out",           // Lorica outside facing
	Port{4082, 17}:   "lorica-out",           // Lorica outside facing
	Port{4083, 6}:    "lorica-out-sec",       // Lorica outside facing (SSL)
	Port{4083, 17}:   "lorica-out-sec",       // Lorica outside facing (SSL)
	Port{4084, 6}:    "fortisphere-vm",       // Fortisphere VM Service
	Port{4084, 17}:   "fortisphere-vm",       // Fortisphere VM Service
	Port{4085, 6}:    "ezmessagesrv",         // EZNews Newsroom Message Service
	Port{4086, 6}:    "ftsync",               // Firewall NAT state table synchronization
	Port{4086, 17}:   "ftsync",               // Firewall NAT state table synchronization
	Port{4087, 6}:    "applusservice",        // APplus Service
	Port{4088, 6}:    "npsp",                 // Noah Printing Service Protocol
	Port{4089, 6}:    "opencore",             // OpenCORE Remote Control Service
	Port{4089, 17}:   "opencore",             // OpenCORE Remote Control Service
	Port{4090, 6}:    "omasgport",            // OMA BCAST Service Guide
	Port{4090, 17}:   "omasgport",            // OMA BCAST Service Guide
	Port{4091, 6}:    "ewinstaller",          // EminentWare Installer
	Port{4091, 17}:   "ewinstaller",          // EminentWare Installer
	Port{4092, 6}:    "ewdgs",                // EminentWare DGS
	Port{4092, 17}:   "ewdgs",                // EminentWare DGS
	Port{4093, 6}:    "pvxpluscs",            // Pvx Plus CS Host
	Port{4093, 17}:   "pvxpluscs",            // Pvx Plus CS Host
	Port{4094, 6}:    "sysrqd",               // sysrq daemon
	Port{4094, 17}:   "sysrqd",               // sysrq daemon
	Port{4095, 6}:    "xtgui",                // xtgui information service
	Port{4095, 17}:   "xtgui",                // xtgui information service
	Port{4096, 6}:    "bre",                  // BRE (Bridge Relay Element)
	Port{4096, 17}:   "bre",                  // BRE (Bridge Relay Element)
	Port{4097, 6}:    "patrolview",           // Patrol View
	Port{4097, 17}:   "patrolview",           // Patrol View
	Port{4098, 6}:    "drmsfsd",              // Missing description for drmsfsd
	Port{4098, 17}:   "drmsfsd",              // Missing description for drmsfsd
	Port{4099, 6}:    "dpcp",                 // Missing description for dpcp
	Port{4099, 17}:   "dpcp",                 // DPCP
	Port{4100, 6}:    "igo-incognito",        // IGo Incognito Data Port
	Port{4100, 17}:   "igo-incognito",        // IGo Incognito Data Port
	Port{4101, 6}:    "brlp-0",               // Braille protocol
	Port{4101, 17}:   "brlp-0",               // Braille protocol
	Port{4102, 6}:    "brlp-1",               // Braille protocol
	Port{4102, 17}:   "brlp-1",               // Braille protocol
	Port{4103, 6}:    "brlp-2",               // Braille protocol
	Port{4103, 17}:   "brlp-2",               // Braille protocol
	Port{4104, 6}:    "brlp-3",               // Braille protocol
	Port{4104, 17}:   "brlp-3",               // Braille protocol
	Port{4105, 6}:    "shofarplayer",         // shofar | Shofar
	Port{4105, 17}:   "shofarplayer",         // ShofarPlayer
	Port{4106, 6}:    "synchronite",          // Missing description for synchronite
	Port{4106, 17}:   "synchronite",          // Synchronite
	Port{4107, 6}:    "j-ac",                 // JDL Accounting LAN Service
	Port{4107, 17}:   "j-ac",                 // JDL Accounting LAN Service
	Port{4108, 6}:    "accel",                // Missing description for accel
	Port{4108, 17}:   "accel",                // ACCEL
	Port{4109, 6}:    "izm",                  // Instantiated Zero-control Messaging
	Port{4109, 17}:   "izm",                  // Instantiated Zero-control Messaging
	Port{4110, 6}:    "g2tag",                // G2 RFID Tag Telemetry Data
	Port{4110, 17}:   "g2tag",                // G2 RFID Tag Telemetry Data
	Port{4111, 6}:    "xgrid",                // Missing description for xgrid
	Port{4111, 17}:   "xgrid",                // Xgrid
	Port{4112, 6}:    "apple-vpns-rp",        // Apple VPN Server Reporting Protocol
	Port{4112, 17}:   "apple-vpns-rp",        // Apple VPN Server Reporting Protocol
	Port{4113, 6}:    "aipn-reg",             // AIPN LS Registration
	Port{4113, 17}:   "aipn-reg",             // AIPN LS Registration
	Port{4114, 6}:    "jomamqmonitor",        // Missing description for jomamqmonitor
	Port{4114, 17}:   "jomamqmonitor",        // JomaMQMonitor
	Port{4115, 6}:    "cds",                  // CDS Transfer Agent
	Port{4115, 17}:   "cds",                  // CDS Transfer Agent
	Port{4116, 6}:    "smartcard-tls",        // Missing description for smartcard-tls
	Port{4116, 17}:   "smartcard-tls",        // smartcard-TLS
	Port{4117, 6}:    "hillrserv",            // Hillr Connection Manager
	Port{4117, 17}:   "hillrserv",            // Hillr Connection Manager
	Port{4118, 6}:    "netscript",            // Netadmin Systems NETscript service
	Port{4118, 17}:   "netscript",            // Netadmin Systems NETscript service
	Port{4119, 6}:    "assuria-slm",          // Assuria Log Manager
	Port{4119, 17}:   "assuria-slm",          // Assuria Log Manager
	Port{4120, 6}:    "minirem",              // MiniRem Remote Telemetry and Control
	Port{4121, 6}:    "e-builder",            // e-Builder Application Communication
	Port{4121, 17}:   "e-builder",            // e-Builder Application Communication
	Port{4122, 6}:    "fprams",               // Fiber Patrol Alarm Service
	Port{4122, 17}:   "fprams",               // Fiber Patrol Alarm Service
	Port{4123, 6}:    "z-wave",               // Zensys Z-Wave Control Protocol | Z-Wave Protocol
	Port{4123, 17}:   "z-wave",               // Zensys Z-Wave Control Protocol
	Port{4124, 6}:    "tigv2",                // Rohill TetraNode Ip Gateway v2
	Port{4124, 17}:   "tigv2",                // Rohill TetraNode Ip Gateway v2
	Port{4125, 6}:    "rww",                  // opsview-envoy | Microsoft Remote Web Workplace on Small Business Server | Opsview Envoy
	Port{4125, 17}:   "opsview-envoy",        // Opsview Envoy
	Port{4126, 6}:    "ddrepl",               // Data Domain Replication Service
	Port{4126, 17}:   "ddrepl",               // Data Domain Replication Service
	Port{4127, 6}:    "unikeypro",            // NetUniKeyServer
	Port{4127, 17}:   "unikeypro",            // NetUniKeyServer
	Port{4128, 6}:    "nufw",                 // NuFW decision delegation protocol
	Port{4128, 17}:   "nufw",                 // NuFW decision delegation protocol
	Port{4129, 6}:    "nuauth",               // NuFW authentication protocol
	Port{4129, 17}:   "nuauth",               // NuFW authentication protocol
	Port{4130, 6}:    "fronet",               // FRONET message protocol
	Port{4130, 17}:   "fronet",               // FRONET message protocol
	Port{4131, 6}:    "stars",                // Global Maintech Stars
	Port{4131, 17}:   "stars",                // Global Maintech Stars
	Port{4132, 6}:    "nuts_dem",             // nuts-dem | NUTS Daemon
	Port{4132, 17}:   "nuts_dem",             // NUTS Daemon
	Port{4133, 6}:    "nuts_bootp",           // nuts-bootp | NUTS Bootp Server
	Port{4133, 17}:   "nuts_bootp",           // NUTS Bootp Server
	Port{4134, 6}:    "nifty-hmi",            // NIFTY-Serve HMI protocol
	Port{4134, 17}:   "nifty-hmi",            // NIFTY-Serve HMI protocol
	Port{4135, 6}:    "cl-db-attach",         // Classic Line Database Server Attach
	Port{4135, 17}:   "cl-db-attach",         // Classic Line Database Server Attach
	Port{4136, 6}:    "cl-db-request",        // Classic Line Database Server Request
	Port{4136, 17}:   "cl-db-request",        // Classic Line Database Server Request
	Port{4137, 6}:    "cl-db-remote",         // Classic Line Database Server Remote
	Port{4137, 17}:   "cl-db-remote",         // Classic Line Database Server Remote
	Port{4138, 6}:    "nettest",              // Missing description for nettest
	Port{4138, 17}:   "nettest",              // Missing description for nettest
	Port{4139, 6}:    "thrtx",                // Imperfect Networks Server
	Port{4139, 17}:   "thrtx",                // Imperfect Networks Server
	Port{4140, 6}:    "cedros_fds",           // cedros-fds | Cedros Fraud Detection System
	Port{4140, 17}:   "cedros_fds",           // Cedros Fraud Detection System
	Port{4141, 6}:    "oirtgsvc",             // Workflow Server
	Port{4141, 17}:   "oirtgsvc",             // Workflow Server
	Port{4142, 6}:    "oidocsvc",             // Document Server
	Port{4142, 17}:   "oidocsvc",             // Document Server
	Port{4143, 6}:    "oidsr",                // Document Replication
	Port{4143, 17}:   "oidsr",                // Document Replication
	Port{4144, 6}:    "wincim",               // pc windows compuserve.com protocol
	Port{4145, 6}:    "vvr-control",          // VVR Control
	Port{4145, 17}:   "vvr-control",          // VVR Control
	Port{4146, 6}:    "tgcconnect",           // TGCConnect Beacon
	Port{4146, 17}:   "tgcconnect",           // TGCConnect Beacon
	Port{4147, 6}:    "vrxpservman",          // Multum Service Manager
	Port{4147, 17}:   "vrxpservman",          // Multum Service Manager
	Port{4148, 6}:    "hhb-handheld",         // HHB Handheld Client
	Port{4148, 17}:   "hhb-handheld",         // HHB Handheld Client
	Port{4149, 6}:    "agslb",                // A10 GSLB Service
	Port{4149, 17}:   "agslb",                // A10 GSLB Service
	Port{4150, 6}:    "PowerAlert-nsa",       // PowerAlert Network Shutdown Agent
	Port{4150, 17}:   "PowerAlert-nsa",       // PowerAlert Network Shutdown Agent
	Port{4151, 6}:    "menandmice_noh",       // menandmice-noh | Men & Mice Remote Control
	Port{4151, 17}:   "menandmice_noh",       // Men & Mice Remote Control
	Port{4152, 6}:    "idig_mux",             // idig-mux | iDigTech Multiplex
	Port{4152, 17}:   "idig_mux",             // iDigTech Multiplex
	Port{4153, 6}:    "mbl-battd",            // MBL Remote Battery Monitoring
	Port{4153, 17}:   "mbl-battd",            // MBL Remote Battery Monitoring
	Port{4154, 6}:    "atlinks",              // atlinks device discovery
	Port{4154, 17}:   "atlinks",              // atlinks device discovery
	Port{4155, 6}:    "bzr",                  // Bazaar version control system
	Port{4155, 17}:   "bzr",                  // Bazaar version control system
	Port{4156, 6}:    "stat-results",         // STAT Results
	Port{4156, 17}:   "stat-results",         // STAT Results
	Port{4157, 6}:    "stat-scanner",         // STAT Scanner Control
	Port{4157, 17}:   "stat-scanner",         // STAT Scanner Control
	Port{4158, 6}:    "stat-cc",              // STAT Command Center
	Port{4158, 17}:   "stat-cc",              // STAT Command Center
	Port{4159, 6}:    "nss",                  // Network Security Service
	Port{4159, 17}:   "nss",                  // Network Security Service
	Port{4160, 6}:    "jini-discovery",       // Jini Discovery
	Port{4160, 17}:   "jini-discovery",       // Jini Discovery
	Port{4161, 6}:    "omscontact",           // OMS Contact
	Port{4161, 17}:   "omscontact",           // OMS Contact
	Port{4162, 6}:    "omstopology",          // OMS Topology
	Port{4162, 17}:   "omstopology",          // OMS Topology
	Port{4163, 6}:    "silverpeakpeer",       // Silver Peak Peer Protocol
	Port{4163, 17}:   "silverpeakpeer",       // Silver Peak Peer Protocol
	Port{4164, 6}:    "silverpeakcomm",       // Silver Peak Communication Protocol
	Port{4164, 17}:   "silverpeakcomm",       // Silver Peak Communication Protocol
	Port{4165, 6}:    "altcp",                // ArcLink over Ethernet
	Port{4165, 17}:   "altcp",                // ArcLink over Ethernet
	Port{4166, 6}:    "joost",                // Joost Peer to Peer Protocol
	Port{4166, 17}:   "joost",                // Joost Peer to Peer Protocol
	Port{4167, 6}:    "ddgn",                 // DeskDirect Global Network
	Port{4167, 17}:   "ddgn",                 // DeskDirect Global Network
	Port{4168, 6}:    "pslicser",             // PrintSoft License Server
	Port{4168, 17}:   "pslicser",             // PrintSoft License Server
	Port{4169, 6}:    "iadt",                 // iadt-disc | Automation Drive Interface Transport | Internet ADT Discovery Protocol
	Port{4169, 17}:   "iadt-disc",            // Internet ADT Discovery Protocol
	Port{4170, 6}:    "d-cinema-csp",         // SMPTE Content Synchonization Protocol
	Port{4171, 6}:    "ml-svnet",             // Maxlogic Supervisor Communication
	Port{4172, 6}:    "pcoip",                // PC over IP
	Port{4172, 17}:   "pcoip",                // PC over IP
	Port{4173, 6}:    "mma-discovery",        // MMA Device Discovery
	Port{4174, 6}:    "smcluster",            // sm-disc | StorMagic Cluster Services | StorMagic Discovery
	Port{4175, 6}:    "bccp",                 // Brocade Cluster Communication Protocol
	Port{4176, 6}:    "tl-ipcproxy",          // Translattice Cluster IPC Proxy
	Port{4177, 6}:    "wello",                // Wello P2P pubsub service
	Port{4177, 17}:   "wello",                // Wello P2P pubsub service
	Port{4178, 6}:    "storman",              // Missing description for storman
	Port{4178, 17}:   "storman",              // StorMan
	Port{4179, 6}:    "MaxumSP",              // Maxum Services
	Port{4179, 17}:   "MaxumSP",              // Maxum Services
	Port{4180, 6}:    "httpx",                // Missing description for httpx
	Port{4180, 17}:   "httpx",                // HTTPX
	Port{4181, 6}:    "macbak",               // Missing description for macbak
	Port{4181, 17}:   "macbak",               // MacBak
	Port{4182, 6}:    "pcptcpservice",        // Production Company Pro TCP Service
	Port{4182, 17}:   "pcptcpservice",        // Production Company Pro TCP Service
	Port{4183, 6}:    "gmmp",                 // cyborgnet | General Metaverse Messaging Protocol | CyborgNet communications protocol
	Port{4183, 17}:   "gmmp",                 // General Metaverse Messaging Protocol
	Port{4184, 6}:    "universe_suite",       // universe-suite | UNIVERSE SUITE MESSAGE SERVICE
	Port{4184, 17}:   "universe_suite",       // UNIVERSE SUITE MESSAGE SERVICE
	Port{4185, 6}:    "wcpp",                 // Woven Control Plane Protocol
	Port{4185, 17}:   "wcpp",                 // Woven Control Plane Protocol
	Port{4186, 6}:    "boxbackupstore",       // Box Backup Store Service
	Port{4187, 6}:    "csc_proxy",            // csc-proxy | Cascade Proxy
	Port{4188, 6}:    "vatata",               // Vatata Peer to Peer Protocol
	Port{4188, 17}:   "vatata",               // Vatata Peer to Peer Protocol
	Port{4189, 6}:    "pcep",                 // Path Computation Element Communication Protocol
	Port{4190, 6}:    "sieve",                // ManageSieve Protocol
	Port{4191, 6}:    "dsmipv6",              // Dual Stack MIPv6 NAT Traversal
	Port{4191, 17}:   "dsmipv6",              // Dual Stack MIPv6 NAT Traversal
	Port{4192, 6}:    "azeti",                // azeti-bd | Azeti Agent Service | azeti blinddate
	Port{4192, 17}:   "azeti-bd",             // azeti blinddate
	Port{4193, 6}:    "pvxplusio",            // PxPlus remote file srvr
	Port{4197, 6}:    "hctl",                 // Harman HControl Protocol
	Port{4199, 6}:    "eims-admin",           // Eudora Internet Mail Service (EIMS) admin | EIMS ADMIN
	Port{4199, 17}:   "eims-admin",           // EIMS ADMIN
	Port{4200, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4200, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4201, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4201, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4202, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4202, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4203, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4203, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4204, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4204, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4205, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4205, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4206, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4206, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4207, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4207, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4208, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4208, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4209, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4209, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4210, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4210, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4211, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4211, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4212, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4212, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4213, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4213, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4214, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4214, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4215, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4215, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4216, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4216, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4217, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4217, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4218, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4218, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4219, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4219, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4220, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4220, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4221, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4221, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4222, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4222, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4223, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4223, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4224, 6}:    "xtell",                // Xtell messenging server
	Port{4224, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4225, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4225, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4226, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4226, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4227, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4227, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4228, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4228, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4229, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4229, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4230, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4230, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4231, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4231, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4232, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4232, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4233, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4233, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4234, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4234, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4235, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4235, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4236, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4236, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4237, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4237, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4238, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4238, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4239, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4239, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4240, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4240, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4241, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4241, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4242, 6}:    "vrml-multi-use",       // VRML Multi User Systems or CrashPlan http:  support.code42.com CrashPlan Latest Configuring Network#Networking_FAQs
	Port{4242, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4243, 6}:    "vrml-multi-use",       // VRML Multi User Systems or CrashPlan http:  support.code42.com CrashPlan Latest Configuring Network#Networking_FAQs
	Port{4243, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4244, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4244, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4245, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4245, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4246, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4246, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4247, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4247, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4248, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4248, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4249, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4249, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4250, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4250, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4251, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4251, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4252, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4252, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4253, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4253, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4254, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4254, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4255, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4255, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4256, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4256, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4257, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4257, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4258, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4258, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4259, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4259, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4260, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4260, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4261, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4261, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4262, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4262, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4263, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4263, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4264, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4264, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4265, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4265, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4266, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4266, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4267, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4267, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4268, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4268, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4269, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4269, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4270, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4270, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4271, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4271, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4272, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4272, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4273, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4273, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4274, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4274, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4275, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4275, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4276, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4276, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4277, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4277, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4278, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4278, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4279, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4279, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4280, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4280, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4281, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4281, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4282, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4282, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4283, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4283, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4284, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4284, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4285, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4285, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4286, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4286, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4287, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4287, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4288, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4288, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4289, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4289, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4290, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4290, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4291, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4291, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4292, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4292, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4293, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4293, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4294, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4294, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4295, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4295, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4296, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4296, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4297, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4297, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4298, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4298, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4299, 6}:    "vrml-multi-use",       // VRML Multi User Systems
	Port{4299, 17}:   "vrml-multi-use",       // VRML Multi User Systems
	Port{4300, 6}:    "corelccam",            // Corel CCam
	Port{4300, 17}:   "corelccam",            // Corel CCam
	Port{4301, 6}:    "d-data",               // Diagnostic Data
	Port{4301, 17}:   "d-data",               // Diagnostic Data
	Port{4302, 6}:    "d-data-control",       // Diagnostic Data Control
	Port{4302, 17}:   "d-data-control",       // Diagnostic Data Control
	Port{4303, 6}:    "srcp",                 // Simple Railroad Command Protocol
	Port{4303, 17}:   "srcp",                 // Simple Railroad Command Protocol
	Port{4304, 6}:    "owserver",             // One-Wire Filesystem Server
	Port{4304, 17}:   "owserver",             // One-Wire Filesystem Server
	Port{4305, 6}:    "batman",               // better approach to mobile ad-hoc networking
	Port{4305, 17}:   "batman",               // better approach to mobile ad-hoc networking
	Port{4306, 6}:    "pinghgl",              // Hellgate London
	Port{4306, 17}:   "pinghgl",              // Hellgate London
	Port{4307, 6}:    "visicron-vs",          // trueconf | Visicron Videoconference Service | TrueConf Videoconference Service
	Port{4307, 17}:   "visicron-vs",          // Visicron Videoconference Service
	Port{4308, 6}:    "compx-lockview",       // Missing description for compx-lockview
	Port{4308, 17}:   "compx-lockview",       // CompX-LockView
	Port{4309, 6}:    "dserver",              // Exsequi Appliance Discovery
	Port{4309, 17}:   "dserver",              // Exsequi Appliance Discovery
	Port{4310, 6}:    "mirrtex",              // Mir-RT exchange service
	Port{4310, 17}:   "mirrtex",              // Mir-RT exchange service
	Port{4311, 6}:    "p6ssmc",               // P6R Secure Server Management Console
	Port{4312, 6}:    "pscl-mgt",             // Parascale Membership Manager
	Port{4313, 6}:    "perrla",               // PERRLA User Services
	Port{4314, 6}:    "choiceview-agt",       // ChoiceView Agent
	Port{4316, 6}:    "choiceview-clt",       // ChoiceView Client
	Port{4320, 6}:    "fdt-rcatp",            // FDT Remote Categorization Protocol
	Port{4320, 17}:   "fdt-rcatp",            // FDT Remote Categorization Protocol
	Port{4321, 6}:    "rwhois",               // Remote Who Is
	Port{4321, 17}:   "rwhois",               // Remote Who Is
	Port{4322, 6}:    "trim-event",           // TRIM Event Service
	Port{4322, 17}:   "trim-event",           // TRIM Event Service
	Port{4323, 6}:    "trim-ice",             // TRIM ICE Service
	Port{4323, 17}:   "trim-ice",             // TRIM ICE Service
	Port{4324, 6}:    "balour",               // Balour Game Server
	Port{4324, 17}:   "balour",               // Balour Game Server
	Port{4325, 6}:    "geognosisman",         // Cadcorp GeognoSIS Manager Service
	Port{4325, 17}:   "geognosisman",         // Cadcorp GeognoSIS Manager Service
	Port{4326, 6}:    "geognosis",            // Cadcorp GeognoSIS Service
	Port{4326, 17}:   "geognosis",            // Cadcorp GeognoSIS Service
	Port{4327, 6}:    "jaxer-web",            // Jaxer Web Protocol
	Port{4327, 17}:   "jaxer-web",            // Jaxer Web Protocol
	Port{4328, 6}:    "jaxer-manager",        // Jaxer Manager Command Protocol
	Port{4328, 17}:   "jaxer-manager",        // Jaxer Manager Command Protocol
	Port{4329, 6}:    "publiqare-sync",       // PubliQare Distributed Environment Synchronisation Engine
	Port{4330, 6}:    "dey-sapi",             // DEY Storage Administration REST API
	Port{4331, 6}:    "ktickets-rest",        // ktickets REST API for event management and ticketing systems (embedded POS devices)
	Port{4333, 6}:    "msql",                 // ahsp | mini-sql server | ArrowHead Service Protocol (AHSP)
	Port{4334, 6}:    "netconf-ch-ssh",       // NETCONF Call Home (SSH)
	Port{4335, 6}:    "netconf-ch-tls",       // NETCONF Call Home (TLS)
	Port{4336, 6}:    "restconf-ch-tls",      // RESTCONF Call Home (TLS)
	Port{4340, 6}:    "gaia",                 // Gaia Connector Protocol
	Port{4340, 17}:   "gaia",                 // Gaia Connector Protocol
	Port{4341, 6}:    "lisp-data",            // LISP Data Packets
	Port{4341, 17}:   "lisp-data",            // LISP Data Packets
	Port{4342, 6}:    "lisp-cons",            // lisp-control | LISP-CONS Control | LISP Control Packets
	Port{4342, 17}:   "lisp-control",         // LISP Data-Triggered Control
	Port{4343, 6}:    "unicall",              // Missing description for unicall
	Port{4343, 17}:   "unicall",              // Missing description for unicall
	Port{4344, 6}:    "vinainstall",          // Missing description for vinainstall
	Port{4344, 17}:   "vinainstall",          // VinaInstall
	Port{4345, 6}:    "m4-network-as",        // Macro 4 Network AS
	Port{4345, 17}:   "m4-network-as",        // Macro 4 Network AS
	Port{4346, 6}:    "elanlm",               // ELAN LM
	Port{4346, 17}:   "elanlm",               // ELAN LM
	Port{4347, 6}:    "lansurveyor",          // LAN Surveyor
	Port{4347, 17}:   "lansurveyor",          // LAN Surveyor
	Port{4348, 6}:    "itose",                // Missing description for itose
	Port{4348, 17}:   "itose",                // ITOSE
	Port{4349, 6}:    "fsportmap",            // File System Port Map
	Port{4349, 17}:   "fsportmap",            // File System Port Map
	Port{4350, 6}:    "net-device",           // Net Device
	Port{4350, 17}:   "net-device",           // Net Device
	Port{4351, 6}:    "plcy-net-svcs",        // PLCY Net Services
	Port{4351, 17}:   "plcy-net-svcs",        // PLCY Net Services
	Port{4352, 6}:    "pjlink",               // Projector Link
	Port{4352, 17}:   "pjlink",               // Projector Link
	Port{4353, 6}:    "f5-iquery",            // F5 iQuery
	Port{4353, 17}:   "f5-iquery",            // F5 iQuery
	Port{4354, 6}:    "qsnet-trans",          // QSNet Transmitter
	Port{4354, 17}:   "qsnet-trans",          // QSNet Transmitter
	Port{4355, 6}:    "qsnet-workst",         // QSNet Workstation
	Port{4355, 17}:   "qsnet-workst",         // QSNet Workstation
	Port{4356, 6}:    "qsnet-assist",         // QSNet Assistant
	Port{4356, 17}:   "qsnet-assist",         // QSNet Assistant
	Port{4357, 6}:    "qsnet-cond",           // QSNet Conductor
	Port{4357, 17}:   "qsnet-cond",           // QSNet Conductor
	Port{4358, 6}:    "qsnet-nucl",           // QSNet Nucleus
	Port{4358, 17}:   "qsnet-nucl",           // QSNet Nucleus
	Port{4359, 6}:    "omabcastltkm",         // OMA BCAST Long-Term Key Messages
	Port{4359, 17}:   "omabcastltkm",         // OMA BCAST Long-Term Key Messages
	Port{4360, 6}:    "matrix_vnet",          // matrix-vnet | Matrix VNet Communication Protocol
	Port{4361, 6}:    "nacnl",                // NavCom Discovery and Control Port
	Port{4361, 17}:   "nacnl",                // NavCom Discovery and Control Port
	Port{4362, 6}:    "afore-vdp-disc",       // AFORE vNode Discovery protocol
	Port{4366, 6}:    "shadowstream",         // ShadowStream System
	Port{4368, 6}:    "wxbrief",              // WeatherBrief Direct
	Port{4368, 17}:   "wxbrief",              // WeatherBrief Direct
	Port{4369, 6}:    "epmd",                 // Erlang Port Mapper Daemon
	Port{4369, 17}:   "epmd",                 // Erlang Port Mapper Daemon
	Port{4370, 6}:    "elpro_tunnel",         // elpro-tunnel | ELPRO V2 Protocol Tunnel
	Port{4370, 17}:   "elpro_tunnel",         // ELPRO V2 Protocol Tunnel
	Port{4371, 6}:    "l2c-control",          // l2c-disc | LAN2CAN Control | LAN2CAN Discovery
	Port{4371, 17}:   "l2c-disc",             // LAN2CAN Discovery
	Port{4372, 6}:    "l2c-data",             // LAN2CAN Data
	Port{4372, 17}:   "l2c-data",             // LAN2CAN Data
	Port{4373, 6}:    "remctl",               // Remote Authenticated Command Service
	Port{4373, 17}:   "remctl",               // Remote Authenticated Command Service
	Port{4374, 6}:    "psi-ptt",              // PSI Push-to-Talk Protocol
	Port{4375, 6}:    "tolteces",             // Toltec EasyShare
	Port{4375, 17}:   "tolteces",             // Toltec EasyShare
	Port{4376, 6}:    "bip",                  // BioAPI Interworking
	Port{4376, 17}:   "bip",                  // BioAPI Interworking
	Port{4377, 6}:    "cp-spxsvr",            // Cambridge Pixel SPx Server
	Port{4377, 17}:   "cp-spxsvr",            // Cambridge Pixel SPx Server
	Port{4378, 6}:    "cp-spxdpy",            // Cambridge Pixel SPx Display
	Port{4378, 17}:   "cp-spxdpy",            // Cambridge Pixel SPx Display
	Port{4379, 6}:    "ctdb",                 // Missing description for ctdb
	Port{4379, 17}:   "ctdb",                 // CTDB
	Port{4389, 6}:    "xandros-cms",          // Xandros Community Management Service
	Port{4389, 17}:   "xandros-cms",          // Xandros Community Management Service
	Port{4390, 6}:    "wiegand",              // Physical Access Control
	Port{4390, 17}:   "wiegand",              // Physical Access Control
	Port{4391, 6}:    "apwi-imserver",        // American Printware IMServer Protocol
	Port{4392, 6}:    "apwi-rxserver",        // American Printware RXServer Protocol
	Port{4393, 6}:    "apwi-rxspooler",       // American Printware RXSpooler Protocol
	Port{4394, 6}:    "apwi-disc",            // American Printware Discovery
	Port{4394, 17}:   "apwi-disc",            // American Printware Discovery
	Port{4395, 6}:    "omnivisionesx",        // OmniVision communication for Virtual environments
	Port{4395, 17}:   "omnivisionesx",        // OmniVision communication for Virtual environments
	Port{4396, 6}:    "fly",                  // Fly Object Space
	Port{4400, 6}:    "ds-srv",               // ASIGRA Services
	Port{4400, 17}:   "ds-srv",               // ASIGRA Services
	Port{4401, 6}:    "ds-srvr",              // ASIGRA Televaulting DS-System Service
	Port{4401, 17}:   "ds-srvr",              // ASIGRA Televaulting DS-System Service
	Port{4402, 6}:    "ds-clnt",              // ASIGRA Televaulting DS-Client Service
	Port{4402, 17}:   "ds-clnt",              // ASIGRA Televaulting DS-Client Service
	Port{4403, 6}:    "ds-user",              // ASIGRA Televaulting DS-Client Monitoring Management
	Port{4403, 17}:   "ds-user",              // ASIGRA Televaulting DS-Client Monitoring Management
	Port{4404, 6}:    "ds-admin",             // ASIGRA Televaulting DS-System Monitoring Management
	Port{4404, 17}:   "ds-admin",             // ASIGRA Televaulting DS-System Monitoring Management
	Port{4405, 6}:    "ds-mail",              // ASIGRA Televaulting Message Level Restore service
	Port{4405, 17}:   "ds-mail",              // ASIGRA Televaulting Message Level Restore service
	Port{4406, 6}:    "ds-slp",               // ASIGRA Televaulting DS-Sleeper Service
	Port{4406, 17}:   "ds-slp",               // ASIGRA Televaulting DS-Sleeper Service
	Port{4407, 6}:    "nacagent",             // Network Access Control Agent
	Port{4408, 6}:    "slscc",                // SLS Technology Control Centre
	Port{4409, 6}:    "netcabinet-com",       // Net-Cabinet comunication
	Port{4410, 6}:    "itwo-server",          // RIB iTWO Application Server
	Port{4411, 6}:    "found",                // Found Messaging Protocol
	Port{4412, 6}:    "smallchat",            // Missing description for smallchat
	Port{4413, 6}:    "avi-nms",              // avi-nms-disc | AVI Systems NMS
	Port{4414, 6}:    "updog",                // Updog Monitoring and Status Framework
	Port{4415, 6}:    "brcd-vr-req",          // Brocade Virtual Router Request
	Port{4416, 6}:    "pjj-player",           // pjj-player-disc | PJJ Media Player | PJJ Media Player discovery
	Port{4417, 6}:    "workflowdir",          // Workflow Director Communication
	Port{4418, 6}:    "axysbridge",           // AXYS communication protocol
	Port{4419, 6}:    "cbp",                  // Colnod Binary Protocol
	Port{4420, 6}:    "nvm-express",          // NVM Express over Fabrics storage access
	Port{4421, 6}:    "scaleft",              // Multi-Platform Remote Management for Cloud Infrastructure
	Port{4422, 6}:    "tsepisp",              // TSEP Installation Service Protocol
	Port{4423, 6}:    "thingkit",             // thingkit secure mesh
	Port{4425, 6}:    "netrockey6",           // NetROCKEY6 SMART Plus Service
	Port{4425, 17}:   "netrockey6",           // NetROCKEY6 SMART Plus Service
	Port{4426, 6}:    "beacon-port-2",        // SMARTS Beacon Port
	Port{4426, 17}:   "beacon-port-2",        // SMARTS Beacon Port
	Port{4427, 6}:    "drizzle",              // Drizzle database server
	Port{4428, 6}:    "omviserver",           // OMV-Investigation Server-Client
	Port{4429, 6}:    "omviagent",            // OMV Investigation Agent-Server
	Port{4430, 6}:    "rsqlserver",           // REAL SQL Server
	Port{4430, 17}:   "rsqlserver",           // REAL SQL Server
	Port{4431, 6}:    "wspipe",               // adWISE Pipe
	Port{4432, 6}:    "l-acoustics",          // L-ACOUSTICS management
	Port{4433, 6}:    "vop",                  // Versile Object Protocol
	Port{4441, 6}:    "netblox",              // Netblox Protocol
	Port{4441, 17}:   "netblox",              // Netblox Protocol
	Port{4442, 6}:    "saris",                // Missing description for saris
	Port{4442, 17}:   "saris",                // Saris
	Port{4443, 6}:    "pharos",               // Missing description for pharos
	Port{4443, 17}:   "pharos",               // Missing description for pharos
	Port{4444, 6}:    "krb524",               // nv-video | Kerberos 5 to 4 ticket xlator | NV Video default
	Port{4444, 17}:   "krb524",               // Missing description for krb524
	Port{4445, 6}:    "upnotifyp",            // Missing description for upnotifyp
	Port{4445, 17}:   "upnotifyp",            // UPNOTIFYP
	Port{4446, 6}:    "n1-fwp",               // Missing description for n1-fwp
	Port{4446, 17}:   "n1-fwp",               // N1-FWP
	Port{4447, 6}:    "n1-rmgmt",             // Missing description for n1-rmgmt
	Port{4447, 17}:   "n1-rmgmt",             // N1-RMGMT
	Port{4448, 6}:    "asc-slmd",             // ASC Licence Manager
	Port{4448, 17}:   "asc-slmd",             // ASC Licence Manager
	Port{4449, 6}:    "privatewire",          // Missing description for privatewire
	Port{4449, 17}:   "privatewire",          // PrivateWire
	Port{4450, 6}:    "camp",                 // Common ASCII Messaging Protocol
	Port{4450, 17}:   "camp",                 // Camp
	Port{4451, 6}:    "ctisystemmsg",         // CTI System Msg
	Port{4451, 17}:   "ctisystemmsg",         // CTI System Msg
	Port{4452, 6}:    "ctiprogramload",       // CTI Program Load
	Port{4452, 17}:   "ctiprogramload",       // CTI Program Load
	Port{4453, 6}:    "nssalertmgr",          // NSS Alert Manager
	Port{4453, 17}:   "nssalertmgr",          // NSS Alert Manager
	Port{4454, 6}:    "nssagentmgr",          // NSS Agent Manager
	Port{4454, 17}:   "nssagentmgr",          // NSS Agent Manager
	Port{4455, 6}:    "prchat-user",          // PR Chat User
	Port{4455, 17}:   "prchat-user",          // PR Chat User
	Port{4456, 6}:    "prchat-server",        // PR Chat Server
	Port{4456, 17}:   "prchat-server",        // PR Chat Server
	Port{4457, 6}:    "prRegister",           // PR Register
	Port{4457, 17}:   "prRegister",           // PR Register
	Port{4458, 6}:    "mcp",                  // Matrix Configuration Protocol
	Port{4458, 17}:   "mcp",                  // Matrix Configuration Protocol
	Port{4480, 6}:    "proxy-plus",           // Proxy+ HTTP proxy port
	Port{4484, 6}:    "hpssmgmt",             // hpssmgmt service
	Port{4484, 17}:   "hpssmgmt",             // hpssmgmt service
	Port{4485, 6}:    "assyst-dr",            // Assyst Data Repository Service
	Port{4486, 6}:    "icms",                 // Integrated Client Message Service
	Port{4486, 17}:   "icms",                 // Integrated Client Message Service
	Port{4487, 6}:    "prex-tcp",             // Protocol for Remote Execution over TCP
	Port{4488, 6}:    "awacs-ice",            // Apple Wide Area Connectivity Service ICE Bootstrap
	Port{4488, 17}:   "awacs-ice",            // Apple Wide Area Connectivity Service ICE Bootstrap
	Port{4500, 6}:    "sae-urn",              // ipsec-nat-t | IPsec NAT-Traversal
	Port{4500, 17}:   "nat-t-ike",            // IKE Nat Traversal negotiation (RFC3947)
	Port{4502, 6}:    "a25-fap-fgw",          // A25 (FAP-FGW)
	Port{4534, 6}:    "armagetronad",         // Armagetron Advanced Game Server
	Port{4535, 6}:    "ehs",                  // Event Heap Server
	Port{4535, 17}:   "ehs",                  // Event Heap Server
	Port{4536, 6}:    "ehs-ssl",              // Event Heap Server SSL
	Port{4536, 17}:   "ehs-ssl",              // Event Heap Server SSL
	Port{4537, 6}:    "wssauthsvc",           // WSS Security Service
	Port{4537, 17}:   "wssauthsvc",           // WSS Security Service
	Port{4538, 6}:    "swx-gate",             // Software Data Exchange Gateway
	Port{4538, 17}:   "swx-gate",             // Software Data Exchange Gateway
	Port{4545, 6}:    "worldscores",          // Missing description for worldscores
	Port{4545, 17}:   "worldscores",          // WorldScores
	Port{4546, 6}:    "sf-lm",                // SF License Manager (Sentinel)
	Port{4546, 17}:   "sf-lm",                // SF License Manager (Sentinel)
	Port{4547, 6}:    "lanner-lm",            // Lanner License Manager
	Port{4547, 17}:   "lanner-lm",            // Lanner License Manager
	Port{4548, 6}:    "synchromesh",          // Missing description for synchromesh
	Port{4548, 17}:   "synchromesh",          // Synchromesh
	Port{4549, 6}:    "aegate",               // Aegate PMR Service
	Port{4549, 17}:   "aegate",               // Aegate PMR Service
	Port{4550, 6}:    "gds-adppiw-db",        // Perman I Interbase Server
	Port{4550, 17}:   "gds-adppiw-db",        // Perman I Interbase Server
	Port{4551, 6}:    "ieee-mih",             // MIH Services
	Port{4551, 17}:   "ieee-mih",             // MIH Services
	Port{4552, 6}:    "menandmice-mon",       // Men and Mice Monitoring
	Port{4552, 17}:   "menandmice-mon",       // Men and Mice Monitoring
	Port{4553, 6}:    "icshostsvc",           // ICS host services
	Port{4554, 6}:    "msfrs",                // MS FRS Replication
	Port{4554, 17}:   "msfrs",                // MS FRS Replication
	Port{4555, 6}:    "rsip",                 // RSIP Port
	Port{4555, 17}:   "rsip",                 // RSIP Port
	Port{4556, 6}:    "dtn-bundle-tcp",       // dtn-bundle | DTN Bundle TCP CL Protocol | DTN Bundle UDP CL Protocol | DTN Bundle DCCP CL Protocol
	Port{4556, 17}:   "dtn-bundle-udp",       // DTN Bundle UDP CL Protocol
	Port{4557, 6}:    "fax",                  // mtcevrunqss | FlexFax FAX transmission service | Marathon everRun Quorum Service Server
	Port{4557, 17}:   "mtcevrunqss",          // Marathon everRun Quorum Service Server
	Port{4558, 6}:    "mtcevrunqman",         // Marathon everRun Quorum Service Manager
	Port{4558, 17}:   "mtcevrunqman",         // Marathon everRun Quorum Service Manager
	Port{4559, 6}:    "hylafax",              // HylaFAX client-server protocol
	Port{4559, 17}:   "hylafax",              // HylaFAX
	Port{4563, 6}:    "amahi-anywhere",       // Amahi Anywhere
	Port{4566, 6}:    "kwtc",                 // Kids Watch Time Control Service
	Port{4566, 17}:   "kwtc",                 // Kids Watch Time Control Service
	Port{4567, 6}:    "tram",                 // Missing description for tram
	Port{4567, 17}:   "tram",                 // TRAM
	Port{4568, 6}:    "bmc-reporting",        // BMC Reporting
	Port{4568, 17}:   "bmc-reporting",        // BMC Reporting
	Port{4569, 6}:    "iax",                  // Inter-Asterisk eXchange
	Port{4569, 17}:   "iax",                  // Inter-Asterisk eXchange
	Port{4570, 6}:    "deploymentmap",        // Service to distribute and update within a site deployment information for Oracle Communications Suite
	Port{4573, 6}:    "cardifftec-back",      // A port for communication between a server and client for a custom backup system
	Port{4590, 6}:    "rid",                  // RID over HTTP TLS
	Port{4591, 6}:    "l3t-at-an",            // HRPD L3T (AT-AN)
	Port{4591, 17}:   "l3t-at-an",            // HRPD L3T (AT-AN)
	Port{4592, 6}:    "hrpd-ith-at-an",       // HRPD-ITH (AT-AN)
	Port{4592, 17}:   "hrpd-ith-at-an",       // HRPD-ITH (AT-AN)
	Port{4593, 6}:    "ipt-anri-anri",        // IPT (ANRI-ANRI)
	Port{4593, 17}:   "ipt-anri-anri",        // IPT (ANRI-ANRI)
	Port{4594, 6}:    "ias-session",          // IAS-Session (ANRI-ANRI)
	Port{4594, 17}:   "ias-session",          // IAS-Session (ANRI-ANRI)
	Port{4595, 6}:    "ias-paging",           // IAS-Paging (ANRI-ANRI)
	Port{4595, 17}:   "ias-paging",           // IAS-Paging (ANRI-ANRI)
	Port{4596, 6}:    "ias-neighbor",         // IAS-Neighbor (ANRI-ANRI)
	Port{4596, 17}:   "ias-neighbor",         // IAS-Neighbor (ANRI-ANRI)
	Port{4597, 6}:    "a21-an-1xbs",          // A21 (AN-1xBS)
	Port{4597, 17}:   "a21-an-1xbs",          // A21 (AN-1xBS)
	Port{4598, 6}:    "a16-an-an",            // A16 (AN-AN)
	Port{4598, 17}:   "a16-an-an",            // A16 (AN-AN)
	Port{4599, 6}:    "a17-an-an",            // A17 (AN-AN)
	Port{4599, 17}:   "a17-an-an",            // A17 (AN-AN)
	Port{4600, 6}:    "piranha1",             // Missing description for piranha1
	Port{4600, 17}:   "piranha1",             // Piranha1
	Port{4601, 6}:    "piranha2",             // Missing description for piranha2
	Port{4601, 17}:   "piranha2",             // Piranha2
	Port{4602, 6}:    "mtsserver",            // EAX MTS Server
	Port{4603, 6}:    "menandmice-upg",       // Men & Mice Upgrade Agent
	Port{4604, 6}:    "irp",                  // Identity Registration Protocol
	Port{4605, 6}:    "sixchat",              // Direct End to End Secure Chat Protocol
	Port{4621, 6}:    "ventoso",              // Bidirectional single port remote radio VOIP and Control stream
	Port{4658, 6}:    "playsta2-app",         // PlayStation2 App Port
	Port{4658, 17}:   "playsta2-app",         // PlayStation2 App Port
	Port{4659, 6}:    "playsta2-lob",         // PlayStation2 Lobby Port
	Port{4659, 17}:   "playsta2-lob",         // PlayStation2 Lobby Port
	Port{4660, 6}:    "mosmig",               // OpenMOSix MIGrates local processes | smaclmgr
	Port{4660, 17}:   "smaclmgr",             // Missing description for smaclmgr
	Port{4661, 6}:    "kar2ouche",            // Kar2ouche Peer location service
	Port{4661, 17}:   "kar2ouche",            // Kar2ouche Peer location service
	Port{4662, 6}:    "edonkey",              // oms | eDonkey file sharing (Donkey) | OrbitNet Message Service
	Port{4662, 17}:   "oms",                  // OrbitNet Message Service
	Port{4663, 6}:    "noteit",               // Note It! Message Service
	Port{4663, 17}:   "noteit",               // Note It! Message Service
	Port{4664, 6}:    "ems",                  // Rimage Messaging Server
	Port{4664, 17}:   "ems",                  // Rimage Messaging Server
	Port{4665, 6}:    "contclientms",         // Container Client Message Service
	Port{4665, 17}:   "contclientms",         // Container Client Message Service
	Port{4666, 6}:    "eportcomm",            // E-Port Message Service
	Port{4666, 17}:   "edonkey",              // eDonkey file sharing (Donkey)
	Port{4667, 6}:    "mmacomm",              // MMA Comm Services
	Port{4667, 17}:   "mmacomm",              // MMA Comm Services
	Port{4668, 6}:    "mmaeds",               // MMA EDS Service
	Port{4668, 17}:   "mmaeds",               // MMA EDS Service
	Port{4669, 6}:    "eportcommdata",        // E-Port Data Service
	Port{4669, 17}:   "eportcommdata",        // E-Port Data Service
	Port{4670, 6}:    "light",                // Light packets transfer protocol
	Port{4670, 17}:   "light",                // Light packets transfer protocol
	Port{4671, 6}:    "acter",                // Bull RSF action server
	Port{4671, 17}:   "acter",                // Bull RSF action server
	Port{4672, 6}:    "rfa",                  // remote file access server
	Port{4672, 17}:   "rfa",                  // remote file access server
	Port{4673, 6}:    "cxws",                 // CXWS Operations
	Port{4673, 17}:   "cxws",                 // CXWS Operations
	Port{4674, 6}:    "appiq-mgmt",           // AppIQ Agent Management
	Port{4674, 17}:   "appiq-mgmt",           // AppIQ Agent Management
	Port{4675, 6}:    "dhct-status",          // BIAP Device Status
	Port{4675, 17}:   "dhct-status",          // BIAP Device Status
	Port{4676, 6}:    "dhct-alerts",          // BIAP Generic Alert
	Port{4676, 17}:   "dhct-alerts",          // BIAP Generic Alert
	Port{4677, 6}:    "bcs",                  // Business Continuity Servi
	Port{4677, 17}:   "bcs",                  // Business Continuity Servi
	Port{4678, 6}:    "traversal",            // boundary traversal
	Port{4678, 17}:   "traversal",            // boundary traversal
	Port{4679, 6}:    "mgesupervision",       // MGE UPS Supervision
	Port{4679, 17}:   "mgesupervision",       // MGE UPS Supervision
	Port{4680, 6}:    "mgemanagement",        // MGE UPS Management
	Port{4680, 17}:   "mgemanagement",        // MGE UPS Management
	Port{4681, 6}:    "parliant",             // Parliant Telephony System
	Port{4681, 17}:   "parliant",             // Parliant Telephony System
	Port{4682, 6}:    "finisar",              // Missing description for finisar
	Port{4682, 17}:   "finisar",              // Missing description for finisar
	Port{4683, 6}:    "spike",                // Spike Clipboard Service
	Port{4683, 17}:   "spike",                // Spike Clipboard Service
	Port{4684, 6}:    "rfid-rp1",             // RFID Reader Protocol 1.0
	Port{4684, 17}:   "rfid-rp1",             // RFID Reader Protocol 1.0
	Port{4685, 6}:    "autopac",              // Autopac Protocol
	Port{4685, 17}:   "autopac",              // Autopac Protocol
	Port{4686, 6}:    "msp-os",               // Manina Service Protocol
	Port{4686, 17}:   "msp-os",               // Manina Service Protocol
	Port{4687, 6}:    "nst",                  // Network Scanner Tool FTP
	Port{4687, 17}:   "nst",                  // Network Scanner Tool FTP
	Port{4688, 6}:    "mobile-p2p",           // Mobile P2P Service
	Port{4688, 17}:   "mobile-p2p",           // Mobile P2P Service
	Port{4689, 6}:    "altovacentral",        // Altova DatabaseCentral
	Port{4689, 17}:   "altovacentral",        // Altova DatabaseCentral
	Port{4690, 6}:    "prelude",              // Prelude IDS message proto
	Port{4690, 17}:   "prelude",              // Prelude IDS message proto
	Port{4691, 6}:    "mtn",                  // monotone Netsync Protocol
	Port{4691, 17}:   "mtn",                  // monotone Netsync Protocol
	Port{4692, 6}:    "conspiracy",           // Conspiracy messaging
	Port{4692, 17}:   "conspiracy",           // Conspiracy messaging
	Port{4700, 6}:    "netxms-agent",         // NetXMS Agent
	Port{4700, 17}:   "netxms-agent",         // NetXMS Agent
	Port{4701, 6}:    "netxms-mgmt",          // NetXMS Management
	Port{4701, 17}:   "netxms-mgmt",          // NetXMS Management
	Port{4702, 6}:    "netxms-sync",          // NetXMS Server Synchronization
	Port{4702, 17}:   "netxms-sync",          // NetXMS Server Synchronization
	Port{4703, 6}:    "npqes-test",           // Network Performance Quality Evaluation System Test Service
	Port{4704, 6}:    "assuria-ins",          // Assuria Insider
	Port{4711, 6}:    "trinity-dist",         // Trinity Trust Network Node Communication
	Port{4713, 6}:    "pulseaudio",           // Pulse Audio UNIX sound framework
	Port{4725, 6}:    "truckstar",            // TruckStar Service
	Port{4725, 17}:   "truckstar",            // TruckStar Service
	Port{4726, 6}:    "a26-fap-fgw",          // A26 (FAP-FGW)
	Port{4726, 17}:   "a26-fap-fgw",          // A26 (FAP-FGW)
	Port{4727, 6}:    "fcis",                 // fcis-disc | F-Link Client Information Service | F-Link Client Information Service Discovery
	Port{4727, 17}:   "fcis-disc",            // F-Link Client Information Service Discovery
	Port{4728, 6}:    "capmux",               // CA Port Multiplexer
	Port{4728, 17}:   "capmux",               // CA Port Multiplexer
	Port{4729, 6}:    "gsmtap",               // GSM Interface Tap
	Port{4729, 17}:   "gsmtap",               // GSM Interface Tap
	Port{4730, 6}:    "gearman",              // Gearman Job Queue System
	Port{4730, 17}:   "gearman",              // Gearman Job Queue System
	Port{4731, 6}:    "remcap",               // Remote Capture Protocol
	Port{4732, 6}:    "ohmtrigger",           // OHM server trigger
	Port{4732, 17}:   "ohmtrigger",           // OHM server trigger
	Port{4733, 6}:    "resorcs",              // RES Orchestration Catalog Services
	Port{4737, 6}:    "ipdr-sp",              // IPDR SP
	Port{4737, 17}:   "ipdr-sp",              // IPDR SP
	Port{4738, 6}:    "solera-lpn",           // SoleraTec Locator
	Port{4738, 17}:   "solera-lpn",           // SoleraTec Locator
	Port{4739, 132}:  "ipfix",                // IP Flow Info Export
	Port{4739, 6}:    "ipfix",                // IP Flow Info Export
	Port{4739, 17}:   "ipfix",                // IP Flow Info Export
	Port{4740, 132}:  "ipfixs",               // IP Flow Info Export over DTLS | ipfix protocol over TLS | ipfix protocol over DTLS
	Port{4740, 6}:    "ipfixs",               // IP Flow Info Export over TLS
	Port{4740, 17}:   "ipfixs",               // IP Flow Info Export over DTLS
	Port{4741, 6}:    "lumimgrd",             // Luminizer Manager
	Port{4741, 17}:   "lumimgrd",             // Luminizer Manager
	Port{4742, 6}:    "sicct",                // sicct-sdp | SICCT Service Discovery Protocol
	Port{4742, 17}:   "sicct-sdp",            // SICCT Service Discovery Protocol
	Port{4743, 6}:    "openhpid",             // openhpi HPI service
	Port{4743, 17}:   "openhpid",             // openhpi HPI service
	Port{4744, 6}:    "ifsp",                 // Internet File Synchronization Protocol
	Port{4744, 17}:   "ifsp",                 // Internet File Synchronization Protocol
	Port{4745, 6}:    "fmp",                  // Funambol Mobile Push
	Port{4745, 17}:   "fmp",                  // Funambol Mobile Push
	Port{4746, 6}:    "intelliadm-disc",      // IntelliAdmin Discovery
	Port{4747, 6}:    "buschtrommel",         // peer-to-peer file exchange protocol
	Port{4749, 6}:    "profilemac",           // Profile for Mac
	Port{4749, 17}:   "profilemac",           // Profile for Mac
	Port{4750, 6}:    "ssad",                 // Simple Service Auto Discovery
	Port{4750, 17}:   "ssad",                 // Simple Service Auto Discovery
	Port{4751, 6}:    "spocp",                // Simple Policy Control Protocol
	Port{4751, 17}:   "spocp",                // Simple Policy Control Protocol
	Port{4752, 6}:    "snap",                 // Simple Network Audio Protocol
	Port{4752, 17}:   "snap",                 // Simple Network Audio Protocol
	Port{4753, 6}:    "simon",                // simon-disc | Simple Invocation of Methods Over Network (SIMON) | Simple Invocation of Methods Over Network (SIMON) Discovery
	Port{4754, 6}:    "gre-in-udp",           // GRE-in-UDP Encapsulation
	Port{4755, 6}:    "gre-udp-dtls",         // GRE-in-UDP Encapsulation with DTLS
	Port{4756, 6}:    "RDCenter",             // Reticle Decision Center
	Port{4774, 6}:    "converge",             // Converge RPC
	Port{4784, 6}:    "bfd-multi-ctl",        // BFD Multihop Control
	Port{4784, 17}:   "bfd-multi-ctl",        // BFD Multihop Control
	Port{4785, 6}:    "cncp",                 // Cisco Nexus Control Protocol
	Port{4785, 17}:   "cncp",                 // Cisco Nexus Control Protocol
	Port{4786, 6}:    "smart-install",        // Smart Install Service
	Port{4787, 6}:    "sia-ctrl-plane",       // Service Insertion Architecture (SIA) Control-Plane
	Port{4788, 6}:    "xmcp",                 // eXtensible Messaging Client Protocol
	Port{4789, 6}:    "vxlan",                // Virtual eXtensible Local Area Network (VXLAN)
	Port{4790, 6}:    "vxlan-gpe",            // Generic Protocol Extension for Virtual eXtensible Local Area Network (VXLAN)
	Port{4791, 6}:    "roce",                 // IP Routable RocE
	Port{4800, 6}:    "iims",                 // Icona Instant Messenging System
	Port{4800, 17}:   "iims",                 // Icona Instant Messenging System
	Port{4801, 6}:    "iwec",                 // Icona Web Embedded Chat
	Port{4801, 17}:   "iwec",                 // Icona Web Embedded Chat
	Port{4802, 6}:    "ilss",                 // Icona License System Server
	Port{4802, 17}:   "ilss",                 // Icona License System Server
	Port{4803, 6}:    "notateit",             // notateit-disc | Notateit Messaging | Notateit Messaging Discovery
	Port{4803, 17}:   "notateit-disc",        // Notateit Messaging Discovery
	Port{4804, 6}:    "aja-ntv4-disc",        // AJA ntv4 Video System Discovery
	Port{4804, 17}:   "aja-ntv4-disc",        // AJA ntv4 Video System Discovery
	Port{4827, 6}:    "htcp",                 // Missing description for htcp
	Port{4827, 17}:   "squid-htcp",           // Squid proxy HTCP port
	Port{4837, 6}:    "varadero-0",           // Missing description for varadero-0
	Port{4837, 17}:   "varadero-0",           // Varadero-0
	Port{4838, 6}:    "varadero-1",           // Missing description for varadero-1
	Port{4838, 17}:   "varadero-1",           // Varadero-1
	Port{4839, 6}:    "varadero-2",           // Missing description for varadero-2
	Port{4839, 17}:   "varadero-2",           // Varadero-2
	Port{4840, 6}:    "opcua-tcp",            // opcua-udp | OPC UA TCP Protocol | OPC UA Connection Protocol | OPC UA Multicast Datagram Protocol
	Port{4840, 17}:   "opcua-udp",            // OPC UA TCP Protocol
	Port{4841, 6}:    "quosa",                // QUOSA Virtual Library Service
	Port{4841, 17}:   "quosa",                // QUOSA Virtual Library Service
	Port{4842, 6}:    "gw-asv",               // nCode ICE-flow Library AppServer
	Port{4842, 17}:   "gw-asv",               // nCode ICE-flow Library AppServer
	Port{4843, 6}:    "opcua-tls",            // OPC UA TCP Protocol over TLS SSL
	Port{4843, 17}:   "opcua-tls",            // OPC UA TCP Protocol over TLS SSL
	Port{4844, 6}:    "gw-log",               // nCode ICE-flow Library LogServer
	Port{4844, 17}:   "gw-log",               // nCode ICE-flow Library LogServer
	Port{4845, 6}:    "wcr-remlib",           // WordCruncher Remote Library Service
	Port{4845, 17}:   "wcr-remlib",           // WordCruncher Remote Library Service
	Port{4846, 6}:    "contamac_icm",         // contamac-icm | Contamac ICM Service
	Port{4846, 17}:   "contamac_icm",         // Contamac ICM Service
	Port{4847, 6}:    "wfc",                  // Web Fresh Communication
	Port{4847, 17}:   "wfc",                  // Web Fresh Communication
	Port{4848, 6}:    "appserv-http",         // App Server - Admin HTTP
	Port{4848, 17}:   "appserv-http",         // App Server - Admin HTTP
	Port{4849, 6}:    "appserv-https",        // App Server - Admin HTTPS
	Port{4849, 17}:   "appserv-https",        // App Server - Admin HTTPS
	Port{4850, 6}:    "sun-as-nodeagt",       // Sun App Server - NA
	Port{4850, 17}:   "sun-as-nodeagt",       // Sun App Server - NA
	Port{4851, 6}:    "derby-repli",          // Apache Derby Replication
	Port{4851, 17}:   "derby-repli",          // Apache Derby Replication
	Port{4867, 6}:    "unify-debug",          // Unify Debugger
	Port{4867, 17}:   "unify-debug",          // Unify Debugger
	Port{4868, 6}:    "phrelay",              // Photon Relay
	Port{4868, 17}:   "phrelay",              // Photon Relay
	Port{4869, 6}:    "phrelaydbg",           // Photon Relay Debug
	Port{4869, 17}:   "phrelaydbg",           // Photon Relay Debug
	Port{4870, 6}:    "cc-tracking",          // Citcom Tracking Service
	Port{4870, 17}:   "cc-tracking",          // Citcom Tracking Service
	Port{4871, 6}:    "wired",                // Missing description for wired
	Port{4871, 17}:   "wired",                // Wired
	Port{4876, 6}:    "tritium-can",          // Tritium CAN Bus Bridge Service
	Port{4877, 6}:    "lmcs",                 // Lighting Management Control System
	Port{4878, 6}:    "inst-discovery",       // Agilent Instrument Discovery
	Port{4879, 6}:    "wsdl-event",           // WSDL Event Receiver
	Port{4880, 6}:    "hislip",               // IVI High-Speed LAN Instrument Protocol
	Port{4881, 6}:    "socp-t",               // SOCP Time Synchronization Protocol
	Port{4881, 17}:   "socp-t",               // SOCP Time Synchronization Protocol
	Port{4882, 6}:    "socp-c",               // SOCP Control Protocol
	Port{4882, 17}:   "socp-c",               // SOCP Control Protocol
	Port{4883, 6}:    "wmlserver",            // Meier-Phelps License Server
	Port{4884, 6}:    "hivestor",             // HiveStor Distributed File System
	Port{4884, 17}:   "hivestor",             // HiveStor Distributed File System
	Port{4885, 6}:    "abbs",                 // Missing description for abbs
	Port{4885, 17}:   "abbs",                 // ABBS
	Port{4894, 6}:    "lyskom",               // LysKOM Protocol A
	Port{4894, 17}:   "lyskom",               // LysKOM Protocol A
	Port{4899, 6}:    "radmin",               // radmin-port | Radmin (www.radmin.com) remote PC control software | RAdmin Port
	Port{4899, 17}:   "radmin-port",          // RAdmin Port
	Port{4900, 6}:    "hfcs",                 // HyperFileSQL Client Server Database Engine | HFSQL Client Server Database Engine
	Port{4900, 17}:   "hfcs",                 // HyperFileSQL Client Server Database Engine
	Port{4901, 6}:    "flr_agent",            // flr-agent | FileLocator Remote Search Agent
	Port{4902, 6}:    "magiccontrol",         // magicCONROL RF and Data Interface
	Port{4912, 6}:    "lutap",                // Technicolor LUT Access Protocol
	Port{4913, 6}:    "lutcp",                // LUTher Control Protocol
	Port{4914, 6}:    "bones",                // Bones Remote Control
	Port{4914, 17}:   "bones",                // Bones Remote Control
	Port{4915, 6}:    "frcs",                 // Fibics Remote Control Service
	Port{4936, 6}:    "an-signaling",         // Signal protocol port for autonomic networking
	Port{4937, 6}:    "atsc-mh-ssc",          // ATSC-M H Service Signaling Channel
	Port{4937, 17}:   "atsc-mh-ssc",          // ATSC-M H Service Signaling Channel
	Port{4940, 6}:    "eq-office-4940",       // Equitrac Office
	Port{4940, 17}:   "eq-office-4940",       // Equitrac Office
	Port{4941, 6}:    "eq-office-4941",       // Equitrac Office
	Port{4941, 17}:   "eq-office-4941",       // Equitrac Office
	Port{4942, 6}:    "eq-office-4942",       // Equitrac Office
	Port{4942, 17}:   "eq-office-4942",       // Equitrac Office
	Port{4949, 6}:    "munin",                // Munin Graphing Framework
	Port{4949, 17}:   "munin",                // Munin Graphing Framework
	Port{4950, 6}:    "sybasesrvmon",         // Sybase Server Monitor
	Port{4950, 17}:   "sybasesrvmon",         // Sybase Server Monitor
	Port{4951, 6}:    "pwgwims",              // PWG WIMS
	Port{4951, 17}:   "pwgwims",              // PWG WIMS
	Port{4952, 6}:    "sagxtsds",             // SAG Directory Server
	Port{4952, 17}:   "sagxtsds",             // SAG Directory Server
	Port{4953, 6}:    "dbsyncarbiter",        // Synchronization Arbiter
	Port{4969, 6}:    "ccss-qmm",             // CCSS QMessageMonitor
	Port{4969, 17}:   "ccss-qmm",             // CCSS QMessageMonitor
	Port{4970, 6}:    "ccss-qsm",             // CCSS QSystemMonitor
	Port{4970, 17}:   "ccss-qsm",             // CCSS QSystemMonitor
	Port{4971, 6}:    "burp",                 // BackUp and Restore Program
	Port{4980, 6}:    "ctxs-vpp",             // Citrix Virtual Path
	Port{4984, 6}:    "webyast",              // Missing description for webyast
	Port{4985, 6}:    "gerhcs",               // GER HC Standard
	Port{4986, 6}:    "mrip",                 // Model Railway Interface Program
	Port{4986, 17}:   "mrip",                 // Model Railway Interface Program
	Port{4987, 6}:    "maybe-veritas",        // smar-se-port1 | SMAR Ethernet Port 1
	Port{4987, 17}:   "smar-se-port1",        // SMAR Ethernet Port 1
	Port{4988, 6}:    "smar-se-port2",        // SMAR Ethernet Port 2
	Port{4988, 17}:   "smar-se-port2",        // SMAR Ethernet Port 2
	Port{4989, 6}:    "parallel",             // Parallel for GAUSS (tm)
	Port{4989, 17}:   "parallel",             // Parallel for GAUSS (tm)
	Port{4990, 6}:    "busycal",              // BusySync Calendar Synch. Protocol
	Port{4990, 17}:   "busycal",              // BusySync Calendar Synch. Protocol
	Port{4991, 6}:    "vrt",                  // VITA Radio Transport
	Port{4991, 17}:   "vrt",                  // VITA Radio Transport
	Port{4998, 6}:    "maybe-veritas",        // Missing description for maybe-veritas
	Port{4999, 6}:    "hfcs-manager",         // HyperFileSQL Client Server Database Engine Manager | HFSQL Client Server Database Engine Manager
	Port{4999, 17}:   "hfcs-manager",         // HyperFileSQL Client Server Database Engine Manager
	Port{5000, 6}:    "upnp",                 // commplex-main | Universal PnP, also Free Internet Chess Server
	Port{5000, 17}:   "upnp",                 // also complex-main
	Port{5001, 6}:    "commplex-link",        // Missing description for commplex-link
	Port{5001, 17}:   "commplex-link",        // Missing description for commplex-link
	Port{5002, 6}:    "rfe",                  // Radio Free Ethernet | radio free ethernet
	Port{5002, 17}:   "rfe",                  // Radio Free Ethernet
	Port{5003, 6}:    "filemaker",            // fmpro-internal | Filemaker Server - http:  www.filemaker.com ti 104289.html | FileMaker, Inc. - Proprietary transport | FileMaker, Inc. - Proprietary name binding
	Port{5003, 17}:   "filemaker",            // Filemaker Server - http:  www.filemaker.com ti 104289.html
	Port{5004, 6}:    "avt-profile-1",        // RTP media data [RFC 3551][RFC 4571] | RTP media data
	Port{5004, 17}:   "avt-profile-1",        // RTP media data [RFC 3551]
	Port{5005, 6}:    "avt-profile-2",        // RTP control protocol [RFC 3551][RFC 4571] | RTP control protocol
	Port{5005, 17}:   "avt-profile-2",        // RTP control protocol [RFC 3551]
	Port{5006, 6}:    "wsm-server",           // wsm server
	Port{5006, 17}:   "wsm-server",           // wsm server
	Port{5007, 6}:    "wsm-server-ssl",       // wsm server ssl
	Port{5007, 17}:   "wsm-server-ssl",       // wsm server ssl
	Port{5008, 6}:    "synapsis-edge",        // Synapsis EDGE
	Port{5008, 17}:   "synapsis-edge",        // Synapsis EDGE
	Port{5009, 6}:    "airport-admin",        // winfs | Apple AirPort WAP Administration | Microsoft Windows Filesystem
	Port{5009, 17}:   "winfs",                // Microsoft Windows Filesystem
	Port{5010, 6}:    "telelpathstart",       // TelepathStart
	Port{5010, 17}:   "telelpathstart",       // Missing description for telelpathstart
	Port{5011, 6}:    "telelpathattack",      // TelepathAttack
	Port{5011, 17}:   "telelpathattack",      // Missing description for telelpathattack
	Port{5012, 6}:    "nsp",                  // NetOnTap Service
	Port{5012, 17}:   "nsp",                  // NetOnTap Service
	Port{5013, 6}:    "fmpro-v6",             // FileMaker, Inc. - Proprietary transport
	Port{5013, 17}:   "fmpro-v6",             // FileMaker, Inc. - Proprietary transport
	Port{5014, 6}:    "onpsocket",            // Overlay Network Protocol
	Port{5014, 17}:   "onpsocket",            // Overlay Network Protocol
	Port{5015, 6}:    "fmwp",                 // FileMaker, Inc. - Web publishing
	Port{5020, 6}:    "zenginkyo-1",          // Missing description for zenginkyo-1
	Port{5020, 17}:   "zenginkyo-1",          // Missing description for zenginkyo-1
	Port{5021, 6}:    "zenginkyo-2",          // Missing description for zenginkyo-2
	Port{5021, 17}:   "zenginkyo-2",          // Missing description for zenginkyo-2
	Port{5022, 6}:    "mice",                 // mice server
	Port{5022, 17}:   "mice",                 // mice server
	Port{5023, 6}:    "htuilsrv",             // Htuil Server for PLD2
	Port{5023, 17}:   "htuilsrv",             // Htuil Server for PLD2
	Port{5024, 6}:    "scpi-telnet",          // Missing description for scpi-telnet
	Port{5024, 17}:   "scpi-telnet",          // SCPI-TELNET
	Port{5025, 6}:    "scpi-raw",             // Missing description for scpi-raw
	Port{5025, 17}:   "scpi-raw",             // SCPI-RAW
	Port{5026, 6}:    "strexec-d",            // Storix I O daemon (data)
	Port{5026, 17}:   "strexec-d",            // Storix I O daemon (data)
	Port{5027, 6}:    "strexec-s",            // Storix I O daemon (stat)
	Port{5027, 17}:   "strexec-s",            // Storix I O daemon (stat)
	Port{5028, 6}:    "qvr",                  // Quiqum Virtual Relais
	Port{5029, 6}:    "infobright",           // Infobright Database Server
	Port{5029, 17}:   "infobright",           // Infobright Database Server
	Port{5030, 6}:    "surfpass",             // Missing description for surfpass
	Port{5030, 17}:   "surfpass",             // SurfPass
	Port{5031, 6}:    "dmp",                  // Direct Message Protocol
	Port{5031, 17}:   "dmp",                  // Direct Message Protocol
	Port{5032, 6}:    "signacert-agent",      // SignaCert Enterprise Trust Server Agent
	Port{5033, 6}:    "jtnetd-server",        // Janstor Secure Data
	Port{5034, 6}:    "jtnetd-status",        // Janstor Status
	Port{5042, 6}:    "asnaacceler8db",       // Missing description for asnaacceler8db
	Port{5042, 17}:   "asnaacceler8db",       // Missing description for asnaacceler8db
	Port{5043, 6}:    "swxadmin",             // ShopWorX Administration
	Port{5043, 17}:   "swxadmin",             // ShopWorX Administration
	Port{5044, 6}:    "lxi-evntsvc",          // LXI Event Service
	Port{5044, 17}:   "lxi-evntsvc",          // LXI Event Service
	Port{5045, 6}:    "osp",                  // Open Settlement Protocol
	Port{5046, 6}:    "vpm-udp",              // Vishay PM UDP Service
	Port{5046, 17}:   "vpm-udp",              // Vishay PM UDP Service
	Port{5047, 6}:    "iscape",               // iSCAPE Data Broadcasting
	Port{5047, 17}:   "iscape",               // iSCAPE Data Broadcasting
	Port{5048, 6}:    "texai",                // Texai Message Service
	Port{5049, 6}:    "ivocalize",            // iVocalize Web Conference
	Port{5049, 17}:   "ivocalize",            // iVocalize Web Conference
	Port{5050, 6}:    "mmcc",                 // multimedia conference control tool
	Port{5050, 17}:   "mmcc",                 // multimedia conference control tool
	Port{5051, 6}:    "ida-agent",            // ita-agent | Symantec Intruder Alert | ITA Agent
	Port{5051, 17}:   "ita-agent",            // ITA Agent
	Port{5052, 6}:    "ita-manager",          // ITA Manager
	Port{5052, 17}:   "ita-manager",          // ITA Manager
	Port{5053, 6}:    "rlm",                  // rlm-disc | RLM License Server | RLM Discovery Server
	Port{5054, 6}:    "rlm-admin",            // RLM administrative interface
	Port{5055, 6}:    "unot",                 // Missing description for unot
	Port{5055, 17}:   "unot",                 // UNOT
	Port{5056, 6}:    "intecom-ps1",          // Intecom Pointspan 1
	Port{5056, 17}:   "intecom-ps1",          // Intecom Pointspan 1
	Port{5057, 6}:    "intecom-ps2",          // Intecom Pointspan 2
	Port{5057, 17}:   "intecom-ps2",          // Intecom Pointspan 2
	Port{5058, 6}:    "locus-disc",           // Locus Discovery
	Port{5058, 17}:   "locus-disc",           // Locus Discovery
	Port{5059, 6}:    "sds",                  // SIP Directory Services
	Port{5059, 17}:   "sds",                  // SIP Directory Services
	Port{5060, 132}:  "sip",                  // Session Initiation Protocol (SIP)
	Port{5060, 6}:    "sip",                  // Session Initiation Protocol (SIP)
	Port{5060, 17}:   "sip",                  // Session Initiation Protocol (SIP)
	Port{5061, 132}:  "sip-tls",              // sips
	Port{5061, 6}:    "sip-tls",              // SIP-TLS
	Port{5061, 17}:   "sip-tls",              // SIP-TLS
	Port{5062, 6}:    "na-localise",          // Localisation access
	Port{5062, 17}:   "na-localise",          // Localisation access
	Port{5063, 6}:    "csrpc",                // centrify secure RPC
	Port{5064, 6}:    "ca-1",                 // Channel Access 1
	Port{5064, 17}:   "ca-1",                 // Channel Access 1
	Port{5065, 6}:    "ca-2",                 // Channel Access 2
	Port{5065, 17}:   "ca-2",                 // Channel Access 2
	Port{5066, 6}:    "stanag-5066",          // STANAG-5066-SUBNET-INTF
	Port{5066, 17}:   "stanag-5066",          // STANAG-5066-SUBNET-INTF
	Port{5067, 6}:    "authentx",             // Authentx Service
	Port{5067, 17}:   "authentx",             // Authentx Service
	Port{5068, 6}:    "bitforestsrv",         // Bitforest Data Service
	Port{5069, 6}:    "i-net-2000-npr",       // I Net 2000-NPR
	Port{5069, 17}:   "i-net-2000-npr",       // I Net 2000-NPR
	Port{5070, 6}:    "vtsas",                // VersaTrans Server Agent Service
	Port{5070, 17}:   "vtsas",                // VersaTrans Server Agent Service
	Port{5071, 6}:    "powerschool",          // Missing description for powerschool
	Port{5071, 17}:   "powerschool",          // PowerSchool
	Port{5072, 6}:    "ayiya",                // Anything In Anything
	Port{5072, 17}:   "ayiya",                // Anything In Anything
	Port{5073, 6}:    "tag-pm",               // Advantage Group Port Mgr
	Port{5073, 17}:   "tag-pm",               // Advantage Group Port Mgr
	Port{5074, 6}:    "alesquery",            // ALES Query
	Port{5074, 17}:   "alesquery",            // ALES Query
	Port{5075, 6}:    "pvaccess",             // Experimental Physics and Industrial Control System
	Port{5078, 6}:    "pixelpusher",          // PixelPusher pixel data
	Port{5079, 6}:    "cp-spxrpts",           // Cambridge Pixel SPx Reports
	Port{5079, 17}:   "cp-spxrpts",           // Cambridge Pixel SPx Reports
	Port{5080, 6}:    "onscreen",             // OnScreen Data Collection Service
	Port{5080, 17}:   "onscreen",             // OnScreen Data Collection Service
	Port{5081, 6}:    "sdl-ets",              // SDL - Ent Trans Server
	Port{5081, 17}:   "sdl-ets",              // SDL - Ent Trans Server
	Port{5082, 6}:    "qcp",                  // Qpur Communication Protocol
	Port{5082, 17}:   "qcp",                  // Qpur Communication Protocol
	Port{5083, 6}:    "qfp",                  // Qpur File Protocol
	Port{5083, 17}:   "qfp",                  // Qpur File Protocol
	Port{5084, 6}:    "llrp",                 // EPCglobal Low-Level Reader Protocol
	Port{5084, 17}:   "llrp",                 // EPCglobal Low-Level Reader Protocol
	Port{5085, 6}:    "encrypted-llrp",       // EPCglobal Encrypted LLRP
	Port{5085, 17}:   "encrypted-llrp",       // EPCglobal Encrypted LLRP
	Port{5086, 6}:    "aprigo-cs",            // Aprigo Collection Service
	Port{5087, 6}:    "biotic",               // BIOTIC - Binary Internet of Things Interoperable Communication
	Port{5090, 132}:  "car",                  // Candidate AR
	Port{5091, 132}:  "cxtp",                 // Context Transfer Protocol
	Port{5092, 6}:    "magpie",               // Magpie Binary
	Port{5092, 17}:   "magpie",               // Magpie Binary
	Port{5093, 6}:    "sentinel-lm",          // Sentinel LM
	Port{5093, 17}:   "sentinel-lm",          // Sentinel LM
	Port{5094, 6}:    "hart-ip",              // Missing description for hart-ip
	Port{5094, 17}:   "hart-ip",              // HART-IP
	Port{5099, 6}:    "sentlm-srv2srv",       // SentLM Srv2Srv
	Port{5099, 17}:   "sentlm-srv2srv",       // SentLM Srv2Srv
	Port{5100, 6}:    "admd",                 // socalia | (chili!soft asp admin port) or Yahoo pager | Socalia service mux
	Port{5100, 17}:   "socalia",              // Socalia service mux
	Port{5101, 6}:    "admdog",               // talarian-udp | talarian-tcp | (chili!soft asp) | Talarian_TCP | Talarian_UDP
	Port{5101, 17}:   "talarian-udp",         // Talarian_UDP
	Port{5102, 6}:    "admeng",               // oms-nonsecure | (chili!soft asp) | Oracle OMS non-secure
	Port{5102, 17}:   "oms-nonsecure",        // Oracle OMS non-secure
	Port{5103, 6}:    "actifio-c2c",          // Actifio C2C
	Port{5104, 6}:    "tinymessage",          // Missing description for tinymessage
	Port{5104, 17}:   "tinymessage",          // TinyMessage
	Port{5105, 6}:    "hughes-ap",            // Hughes Association Protocol
	Port{5105, 17}:   "hughes-ap",            // Hughes Association Protocol
	Port{5106, 6}:    "actifioudsagent",      // Actifio UDS Agent
	Port{5107, 6}:    "actifioreplic",        // Disk to Disk replication between Actifio Clusters
	Port{5111, 6}:    "taep-as-svc",          // TAEP AS service
	Port{5111, 17}:   "taep-as-svc",          // TAEP AS service
	Port{5112, 6}:    "pm-cmdsvr",            // PeerMe Msg Cmd Service
	Port{5112, 17}:   "pm-cmdsvr",            // PeerMe Msg Cmd Service
	Port{5114, 6}:    "ev-services",          // Enterprise Vault Services
	Port{5115, 6}:    "autobuild",            // Symantec Autobuild Service
	Port{5116, 6}:    "emb-proj-cmd",         // EPSON Projecter Image Transfer
	Port{5116, 17}:   "emb-proj-cmd",         // EPSON Projecter Image Transfer
	Port{5117, 6}:    "gradecam",             // GradeCam Image Processing
	Port{5120, 6}:    "barracuda-bbs",        // Barracuda Backup Protocol
	Port{5133, 6}:    "nbt-pc",               // Policy Commander
	Port{5133, 17}:   "nbt-pc",               // Policy Commander
	Port{5134, 6}:    "ppactivation",         // PP ActivationServer
	Port{5135, 6}:    "erp-scale",            // Missing description for erp-scale
	Port{5136, 6}:    "minotaur-sa",          // Minotaur SA
	Port{5136, 17}:   "minotaur-sa",          // Minotaur SA
	Port{5137, 6}:    "ctsd",                 // MyCTS server port
	Port{5137, 17}:   "ctsd",                 // MyCTS server port
	Port{5145, 6}:    "rmonitor_secure",      // rmonitor-secure | RMONITOR SECURE
	Port{5145, 17}:   "rmonitor_secure",      // Missing description for rmonitor_secure
	Port{5146, 6}:    "social-alarm",         // Social Alarm Service
	Port{5150, 6}:    "atmp",                 // Ascend Tunnel Management Protocol
	Port{5150, 17}:   "atmp",                 // Ascend Tunnel Management Protocol
	Port{5151, 6}:    "esri_sde",             // esri-sde | ESRI SDE Instance | ESRI SDE Remote Start
	Port{5151, 17}:   "esri_sde",             // ESRI SDE Remote Start
	Port{5152, 6}:    "sde-discovery",        // ESRI SDE Instance Discovery
	Port{5152, 17}:   "sde-discovery",        // ESRI SDE Instance Discovery
	Port{5153, 6}:    "toruxserver",          // ToruX Game Server
	Port{5154, 6}:    "bzflag",               // BZFlag game server
	Port{5154, 17}:   "bzflag",               // BZFlag game server
	Port{5155, 6}:    "asctrl-agent",         // Oracle asControl Agent
	Port{5155, 17}:   "asctrl-agent",         // Oracle asControl Agent
	Port{5156, 6}:    "rugameonline",         // Russian Online Game
	Port{5157, 6}:    "mediat",               // Mediat Remote Object Exchange
	Port{5161, 6}:    "snmpssh",              // SNMP over SSH Transport Model
	Port{5162, 6}:    "snmpssh-trap",         // SNMP Notification over SSH Transport Model
	Port{5163, 6}:    "sbackup",              // Shadow Backup
	Port{5164, 6}:    "vpa",                  // vpa-disc | Virtual Protocol Adapter | Virtual Protocol Adapter Discovery
	Port{5164, 17}:   "vpa-disc",             // Virtual Protocol Adapter Discovery
	Port{5165, 6}:    "ife_icorp",            // ife-icorp | ife_1corp
	Port{5165, 17}:   "ife_icorp",            // ife_1corp
	Port{5166, 6}:    "winpcs",               // WinPCS Service Connection
	Port{5166, 17}:   "winpcs",               // WinPCS Service Connection
	Port{5167, 6}:    "scte104",              // SCTE104 Connection
	Port{5167, 17}:   "scte104",              // SCTE104 Connection
	Port{5168, 6}:    "scte30",               // SCTE30 Connection
	Port{5168, 17}:   "scte30",               // SCTE30 Connection
	Port{5172, 6}:    "pcoip-mgmt",           // PC over IP Endpoint Management
	Port{5190, 6}:    "aol",                  // America-Online.  Also can be used by ICQ | America-Online
	Port{5190, 17}:   "aol",                  // America-Online.
	Port{5191, 6}:    "aol-1",                // AmericaOnline1
	Port{5191, 17}:   "aol-1",                // AmericaOnline1
	Port{5192, 6}:    "aol-2",                // AmericaOnline2
	Port{5192, 17}:   "aol-2",                // AmericaOnline2
	Port{5193, 6}:    "aol-3",                // AmericaOnline3
	Port{5193, 17}:   "aol-3",                // AmericaOnline3
	Port{5194, 6}:    "cpscomm",              // CipherPoint Config Service
	Port{5195, 6}:    "ampl-lic",             // The protocol is used by a license server and client programs to control use of program licenses that float to networked machines
	Port{5196, 6}:    "ampl-tableproxy",      // The protocol is used by two programs that exchange "table" data used in the AMPL modeling language
	Port{5197, 6}:    "tunstall-lwp",         // Tunstall Lone worker device interface
	Port{5200, 6}:    "targus-getdata",       // TARGUS GetData
	Port{5200, 17}:   "targus-getdata",       // TARGUS GetData
	Port{5201, 6}:    "targus-getdata1",      // TARGUS GetData 1
	Port{5201, 17}:   "targus-getdata1",      // TARGUS GetData 1
	Port{5202, 6}:    "targus-getdata2",      // TARGUS GetData 2
	Port{5202, 17}:   "targus-getdata2",      // TARGUS GetData 2
	Port{5203, 6}:    "targus-getdata3",      // TARGUS GetData 3
	Port{5203, 17}:   "targus-getdata3",      // TARGUS GetData 3
	Port{5209, 6}:    "nomad",                // Nomad Device Video Transfer
	Port{5215, 6}:    "noteza",               // NOTEZA Data Safety Service
	Port{5221, 6}:    "3exmp",                // 3eTI Extensible Management Protocol for OAMP
	Port{5222, 6}:    "xmpp-client",          // XMPP Client Connection
	Port{5222, 17}:   "xmpp-client",          // XMPP Client Connection
	Port{5223, 6}:    "hpvirtgrp",            // HP Virtual Machine Group Management
	Port{5223, 17}:   "hpvirtgrp",            // HP Virtual Machine Group Management
	Port{5224, 6}:    "hpvirtctrl",           // HP Virtual Machine Console Operations
	Port{5224, 17}:   "hpvirtctrl",           // HP Virtual Machine Console Operations
	Port{5225, 6}:    "hp-server",            // HP Server
	Port{5225, 17}:   "hp-server",            // HP Server
	Port{5226, 6}:    "hp-status",            // HP Status
	Port{5226, 17}:   "hp-status",            // HP Status
	Port{5227, 6}:    "perfd",                // HP System Performance Metric Service
	Port{5227, 17}:   "perfd",                // HP System Performance Metric Service
	Port{5228, 6}:    "hpvroom",              // HP Virtual Room Service
	Port{5229, 6}:    "jaxflow",              // Netflow IPFIX sFlow Collector and Forwarder Management
	Port{5230, 6}:    "jaxflow-data",         // JaxMP RealFlow application and protocol data
	Port{5231, 6}:    "crusecontrol",         // Remote Control of Scan Software for Cruse Scanners
	Port{5232, 6}:    "sgi-dgl",              // csedaemon | SGI Distributed Graphics | Cruse Scanning System Service
	Port{5233, 6}:    "enfs",                 // Etinnae Network File Service
	Port{5234, 6}:    "eenet",                // EEnet communications
	Port{5234, 17}:   "eenet",                // EEnet communications
	Port{5235, 6}:    "galaxy-network",       // Galaxy Network Service
	Port{5235, 17}:   "galaxy-network",       // Galaxy Network Service
	Port{5236, 6}:    "padl2sim",             // Missing description for padl2sim
	Port{5236, 17}:   "padl2sim",             // Missing description for padl2sim
	Port{5237, 6}:    "mnet-discovery",       // m-net discovery
	Port{5237, 17}:   "mnet-discovery",       // m-net discovery
	Port{5245, 6}:    "downtools",            // downtools-disc | DownTools Control Protocol | DownTools Discovery Protocol
	Port{5245, 17}:   "downtools-disc",       // DownTools Discovery Protocol
	Port{5246, 6}:    "capwap-control",       // CAPWAP Control Protocol
	Port{5246, 17}:   "capwap-control",       // CAPWAP Control Protocol
	Port{5247, 6}:    "capwap-data",          // CAPWAP Data Protocol
	Port{5247, 17}:   "capwap-data",          // CAPWAP Data Protocol
	Port{5248, 6}:    "caacws",               // CA Access Control Web Service
	Port{5248, 17}:   "caacws",               // CA Access Control Web Service
	Port{5249, 6}:    "caaclang2",            // CA AC Lang Service
	Port{5249, 17}:   "caaclang2",            // CA AC Lang Service
	Port{5250, 6}:    "soagateway",           // Missing description for soagateway
	Port{5250, 17}:   "soagateway",           // soaGateway
	Port{5251, 6}:    "caevms",               // CA eTrust VM Service
	Port{5251, 17}:   "caevms",               // CA eTrust VM Service
	Port{5252, 6}:    "movaz-ssc",            // Movaz SSC
	Port{5252, 17}:   "movaz-ssc",            // Movaz SSC
	Port{5253, 6}:    "kpdp",                 // Kohler Power Device Protocol
	Port{5254, 6}:    "logcabin",             // LogCabin storage service
	Port{5264, 6}:    "3com-njack-1",         // 3Com Network Jack Port 1
	Port{5264, 17}:   "3com-njack-1",         // 3Com Network Jack Port 1
	Port{5265, 6}:    "3com-njack-2",         // 3Com Network Jack Port 2
	Port{5265, 17}:   "3com-njack-2",         // 3Com Network Jack Port 2
	Port{5269, 6}:    "xmpp-server",          // XMPP Server Connection
	Port{5269, 17}:   "xmpp-server",          // XMPP Server Connection
	Port{5270, 6}:    "xmp",                  // cartographerxmp | Cartographer XMP
	Port{5270, 17}:   "xmp",                  // Cartographer XMP
	Port{5271, 6}:    "cuelink",              // cuelink-disc | StageSoft CueLink messaging | StageSoft CueLink discovery
	Port{5271, 17}:   "cuelink-disc",         // StageSoft CueLink discovery
	Port{5272, 6}:    "pk",                   // Missing description for pk
	Port{5272, 17}:   "pk",                   // PK
	Port{5280, 6}:    "xmpp-bosh",            // Bidirectional-streams Over Synchronous HTTP (BOSH)
	Port{5281, 6}:    "undo-lm",              // Undo License Manager
	Port{5282, 6}:    "transmit-port",        // Marimba Transmitter Port
	Port{5282, 17}:   "transmit-port",        // Marimba Transmitter Port
	Port{5298, 6}:    "presence",             // XMPP Link-Local Messaging
	Port{5298, 17}:   "presence",             // XMPP Link-Local Messaging
	Port{5299, 6}:    "nlg-data",             // NLG Data Service
	Port{5299, 17}:   "nlg-data",             // NLG Data Service
	Port{5300, 6}:    "hacl-hb",              // HA cluster heartbeat
	Port{5300, 17}:   "hacl-hb",              // HA cluster heartbeat
	Port{5301, 6}:    "hacl-gs",              // HA cluster general services
	Port{5301, 17}:   "hacl-gs",              // HA cluster general services
	Port{5302, 6}:    "hacl-cfg",             // HA cluster configuration
	Port{5302, 17}:   "hacl-cfg",             // HA cluster configuration
	Port{5303, 6}:    "hacl-probe",           // HA cluster probing
	Port{5303, 17}:   "hacl-probe",           // HA cluster probing
	Port{5304, 6}:    "hacl-local",           // HA Cluster Commands
	Port{5304, 17}:   "hacl-local",           // Missing description for hacl-local
	Port{5305, 6}:    "hacl-test",            // HA Cluster Test
	Port{5305, 17}:   "hacl-test",            // Missing description for hacl-test
	Port{5306, 6}:    "sun-mc-grp",           // Sun MC Group
	Port{5306, 17}:   "sun-mc-grp",           // Sun MC Group
	Port{5307, 6}:    "sco-aip",              // SCO AIP
	Port{5307, 17}:   "sco-aip",              // SCO AIP
	Port{5308, 6}:    "cfengine",             // Missing description for cfengine
	Port{5308, 17}:   "cfengine",             // Missing description for cfengine
	Port{5309, 6}:    "jprinter",             // J Printer
	Port{5309, 17}:   "jprinter",             // J Printer
	Port{5310, 6}:    "outlaws",              // Missing description for outlaws
	Port{5310, 17}:   "outlaws",              // Outlaws
	Port{5312, 6}:    "permabit-cs",          // Permabit Client-Server
	Port{5312, 17}:   "permabit-cs",          // Permabit Client-Server
	Port{5313, 6}:    "rrdp",                 // Real-time & Reliable Data
	Port{5313, 17}:   "rrdp",                 // Real-time & Reliable Data
	Port{5314, 6}:    "opalis-rbt-ipc",       // Missing description for opalis-rbt-ipc
	Port{5314, 17}:   "opalis-rbt-ipc",       // Missing description for opalis-rbt-ipc
	Port{5315, 6}:    "hacl-poll",            // HA Cluster UDP Polling
	Port{5315, 17}:   "hacl-poll",            // HA Cluster UDP Polling
	Port{5316, 6}:    "hpdevms",              // hpbladems | HP Device Monitor Service | HPBladeSystem Monitor Service
	Port{5316, 17}:   "hpdevms",              // HP Device Monitor Service
	Port{5317, 6}:    "hpdevms",              // HP Device Monitor Service
	Port{5318, 6}:    "pkix-cmc",             // PKIX Certificate Management using CMS (CMC)
	Port{5320, 6}:    "bsfserver-zn",         // Webservices-based Zn interface of BSF
	Port{5321, 6}:    "bsfsvr-zn-ssl",        // Webservices-based Zn interface of BSF over SSL
	Port{5343, 6}:    "kfserver",             // Sculptor Database Server
	Port{5343, 17}:   "kfserver",             // Sculptor Database Server
	Port{5344, 6}:    "xkotodrcp",            // xkoto DRCP
	Port{5344, 17}:   "xkotodrcp",            // xkoto DRCP
	Port{5349, 6}:    "stuns",                // stun-behaviors | turns | STUN over TLS | STUN over DTLS | TURN over TLS | TURN over DTLS | STUN Behavior Discovery over TLS | Reserved for a future enhancement of STUN-BEHAVIOR
	Port{5349, 17}:   "stuns",                // Reserved for a future enhancement of STUN
	Port{5350, 6}:    "nat-pmp-status",       // pcp-multicast | NAT-PMP Status Announcements | Port Control Protocol Multicast
	Port{5350, 17}:   "nat-pmp-status",       // NAT-PMP Status Announcements
	Port{5351, 6}:    "nat-pmp",              // pcp | NAT Port Mapping Protocol | Port Control Protocol
	Port{5351, 17}:   "nat-pmp",              // Missing description for nat-pmp
	Port{5352, 6}:    "dns-llq",              // DNS Long-Lived Queries
	Port{5352, 17}:   "dns-llq",              // DNS Long-Lived Queries
	Port{5353, 6}:    "mdns",                 // Multicast DNS
	Port{5353, 17}:   "zeroconf",             // Mac OS X Bonjour Zeroconf port
	Port{5354, 6}:    "mdnsresponder",        // Multicast DNS Responder IPC
	Port{5354, 17}:   "mdnsresponder",        // Multicast DNS Responder IPC
	Port{5355, 6}:    "llmnr",                // Missing description for llmnr
	Port{5355, 17}:   "llmnr",                // LLMNR
	Port{5356, 6}:    "ms-smlbiz",            // Microsoft Small Business
	Port{5356, 17}:   "ms-smlbiz",            // Microsoft Small Business
	Port{5357, 6}:    "wsdapi",               // Web Services for Devices
	Port{5357, 17}:   "wsdapi",               // Web Services for Devices
	Port{5358, 6}:    "wsdapi-s",             // WS for Devices Secured
	Port{5358, 17}:   "wsdapi-s",             // WS for Devices Secured
	Port{5359, 6}:    "ms-alerter",           // Microsoft Alerter
	Port{5359, 17}:   "ms-alerter",           // Microsoft Alerter
	Port{5360, 6}:    "ms-sideshow",          // Protocol for Windows SideShow
	Port{5360, 17}:   "ms-sideshow",          // Protocol for Windows SideShow
	Port{5361, 6}:    "ms-s-sideshow",        // Secure Protocol for Windows SideShow
	Port{5361, 17}:   "ms-s-sideshow",        // Secure Protocol for Windows SideShow
	Port{5362, 6}:    "serverwsd2",           // Microsoft Windows Server WSD2 Service
	Port{5362, 17}:   "serverwsd2",           // Microsoft Windows Server WSD2 Service
	Port{5363, 6}:    "net-projection",       // Windows Network Projection
	Port{5363, 17}:   "net-projection",       // Windows Network Projection
	Port{5364, 6}:    "kdnet",                // Microsoft Kernel Debugger
	Port{5397, 6}:    "stresstester",         // StressTester(tm) Injector
	Port{5397, 17}:   "stresstester",         // StressTester(tm) Injector
	Port{5398, 6}:    "elektron-admin",       // Elektron Administration
	Port{5398, 17}:   "elektron-admin",       // Elektron Administration
	Port{5399, 6}:    "securitychase",        // Missing description for securitychase
	Port{5399, 17}:   "securitychase",        // SecurityChase
	Port{5400, 6}:    "pcduo-old",            // excerpt | RemCon PC-Duo - old port | Excerpt Search
	Port{5400, 17}:   "excerpt",              // Excerpt Search
	Port{5401, 6}:    "excerpts",             // Excerpt Search Secure
	Port{5401, 17}:   "excerpts",             // Excerpt Search Secure
	Port{5402, 6}:    "mftp",                 // OmniCast MFTP
	Port{5402, 17}:   "mftp",                 // OmniCast MFTP
	Port{5403, 6}:    "hpoms-ci-lstn",        // Missing description for hpoms-ci-lstn
	Port{5403, 17}:   "hpoms-ci-lstn",        // HPOMS-CI-LSTN
	Port{5404, 6}:    "hpoms-dps-lstn",       // Missing description for hpoms-dps-lstn
	Port{5404, 17}:   "hpoms-dps-lstn",       // HPOMS-DPS-LSTN
	Port{5405, 6}:    "pcduo",                // netsupport | RemCon PC-Duo - new port | NetSupport
	Port{5405, 17}:   "netsupport",           // NetSupport
	Port{5406, 6}:    "systemics-sox",        // Systemics Sox
	Port{5406, 17}:   "systemics-sox",        // Systemics Sox
	Port{5407, 6}:    "foresyte-clear",       // Missing description for foresyte-clear
	Port{5407, 17}:   "foresyte-clear",       // Foresyte-Clear
	Port{5408, 6}:    "foresyte-sec",         // Missing description for foresyte-sec
	Port{5408, 17}:   "foresyte-sec",         // Foresyte-Sec
	Port{5409, 6}:    "salient-dtasrv",       // Salient Data Server
	Port{5409, 17}:   "salient-dtasrv",       // Salient Data Server
	Port{5410, 6}:    "salient-usrmgr",       // Salient User Manager
	Port{5410, 17}:   "salient-usrmgr",       // Salient User Manager
	Port{5411, 6}:    "actnet",               // Missing description for actnet
	Port{5411, 17}:   "actnet",               // ActNet
	Port{5412, 6}:    "continuus",            // Missing description for continuus
	Port{5412, 17}:   "continuus",            // Continuus
	Port{5413, 6}:    "wwiotalk",             // Missing description for wwiotalk
	Port{5413, 17}:   "wwiotalk",             // WWIOTALK
	Port{5414, 6}:    "statusd",              // Missing description for statusd
	Port{5414, 17}:   "statusd",              // StatusD
	Port{5415, 6}:    "ns-server",            // NS Server
	Port{5415, 17}:   "ns-server",            // NS Server
	Port{5416, 6}:    "sns-gateway",          // SNS Gateway
	Port{5416, 17}:   "sns-gateway",          // SNS Gateway
	Port{5417, 6}:    "sns-agent",            // SNS Agent
	Port{5417, 17}:   "sns-agent",            // SNS Agent
	Port{5418, 6}:    "mcntp",                // Missing description for mcntp
	Port{5418, 17}:   "mcntp",                // MCNTP
	Port{5419, 6}:    "dj-ice",               // Missing description for dj-ice
	Port{5419, 17}:   "dj-ice",               // DJ-ICE
	Port{5420, 6}:    "cylink-c",             // Missing description for cylink-c
	Port{5420, 17}:   "cylink-c",             // Cylink-C
	Port{5421, 6}:    "netsupport2",          // Net Support 2
	Port{5421, 17}:   "netsupport2",          // Net Support 2
	Port{5422, 6}:    "salient-mux",          // Salient MUX
	Port{5422, 17}:   "salient-mux",          // Salient MUX
	Port{5423, 6}:    "virtualuser",          // Missing description for virtualuser
	Port{5423, 17}:   "virtualuser",          // VIRTUALUSER
	Port{5424, 6}:    "beyond-remote",        // Beyond Remote
	Port{5424, 17}:   "beyond-remote",        // Beyond Remote
	Port{5425, 6}:    "br-channel",           // Beyond Remote Command Channel
	Port{5425, 17}:   "br-channel",           // Beyond Remote Command Channel
	Port{5426, 6}:    "devbasic",             // Missing description for devbasic
	Port{5426, 17}:   "devbasic",             // DEVBASIC
	Port{5427, 6}:    "sco-peer-tta",         // Missing description for sco-peer-tta
	Port{5427, 17}:   "sco-peer-tta",         // SCO-PEER-TTA
	Port{5428, 6}:    "telaconsole",          // Missing description for telaconsole
	Port{5428, 17}:   "omid",                 // OpenMosix Info Dissemination
	Port{5429, 6}:    "base",                 // Billing and Accounting System Exchange
	Port{5429, 17}:   "base",                 // Billing and Accounting System Exchange
	Port{5430, 6}:    "radec-corp",           // RADEC CORP
	Port{5430, 17}:   "radec-corp",           // RADEC CORP
	Port{5431, 6}:    "park-agent",           // PARK AGENT
	Port{5431, 17}:   "park-agent",           // Missing description for park-agent
	Port{5432, 6}:    "postgresql",           // PostgreSQL database server | PostgreSQL Database
	Port{5432, 17}:   "postgresql",           // PostgreSQL Database
	Port{5433, 6}:    "pyrrho",               // Pyrrho DBMS
	Port{5433, 17}:   "pyrrho",               // Pyrrho DBMS
	Port{5434, 6}:    "sgi-arrayd",           // SGI Array Services Daemon
	Port{5434, 17}:   "sgi-arrayd",           // SGI Array Services Daemon
	Port{5435, 6}:    "sceanics",             // SCEANICS situation and action notification
	Port{5435, 17}:   "sceanics",             // SCEANICS situation and action notification
	Port{5436, 6}:    "pmip6-cntl",           // Missing description for pmip6-cntl
	Port{5436, 17}:   "pmip6-cntl",           // Missing description for pmip6-cntl
	Port{5437, 6}:    "pmip6-data",           // Missing description for pmip6-data
	Port{5437, 17}:   "pmip6-data",           // Missing description for pmip6-data
	Port{5443, 6}:    "spss",                 // Pearson HTTPS
	Port{5443, 17}:   "spss",                 // Pearson HTTPS
	Port{5445, 6}:    "smbdirect",            // Server Message Block over Remote Direct Memory Access
	Port{5450, 6}:    "tiepie",               // tiepie-disc | TiePie engineering data acquisition | TiePie engineering data acquisition (discovery)
	Port{5453, 6}:    "surebox",              // Missing description for surebox
	Port{5453, 17}:   "surebox",              // SureBox
	Port{5454, 6}:    "apc-5454",             // APC 5454
	Port{5454, 17}:   "apc-5454",             // APC 5454
	Port{5455, 6}:    "apc-5455",             // APC 5455
	Port{5455, 17}:   "apc-5455",             // APC 5455
	Port{5456, 6}:    "apc-5456",             // APC 5456
	Port{5456, 17}:   "apc-5456",             // APC 5456
	Port{5461, 6}:    "silkmeter",            // Missing description for silkmeter
	Port{5461, 17}:   "silkmeter",            // SILKMETER
	Port{5462, 6}:    "ttl-publisher",        // TTL Publisher
	Port{5462, 17}:   "ttl-publisher",        // TTL Publisher
	Port{5463, 6}:    "ttlpriceproxy",        // TTL Price Proxy
	Port{5463, 17}:   "ttlpriceproxy",        // TTL Price Proxy
	Port{5464, 6}:    "quailnet",             // Quail Networks Object Broker
	Port{5464, 17}:   "quailnet",             // Quail Networks Object Broker
	Port{5465, 6}:    "netops-broker",        // Missing description for netops-broker
	Port{5465, 17}:   "netops-broker",        // NETOPS-BROKER
	Port{5470, 6}:    "apsolab-col",          // The Apsolab company's data collection protocol (native api)
	Port{5471, 6}:    "apsolab-cols",         // The Apsolab company's secure data collection protocol (native api)
	Port{5472, 6}:    "apsolab-tag",          // The Apsolab company's dynamic tag protocol
	Port{5473, 6}:    "apsolab-tags",         // The Apsolab company's secure dynamic tag protocol
	Port{5474, 6}:    "apsolab-rpc",          // The Apsolab company's status query protocol
	Port{5475, 6}:    "apsolab-data",         // The Apsolab company's data retrieval protocol
	Port{5490, 6}:    "connect-proxy",        // Many HTTP CONNECT proxies
	Port{5500, 6}:    "hotline",              // Hotline file sharing client server | fcp-addr-srvr1
	Port{5500, 17}:   "securid",              // SecurID
	Port{5501, 6}:    "fcp-addr-srvr2",       // Missing description for fcp-addr-srvr2
	Port{5501, 17}:   "fcp-addr-srvr2",       // Missing description for fcp-addr-srvr2
	Port{5502, 6}:    "fcp-srvr-inst1",       // Missing description for fcp-srvr-inst1
	Port{5502, 17}:   "fcp-srvr-inst1",       // Missing description for fcp-srvr-inst1
	Port{5503, 6}:    "fcp-srvr-inst2",       // Missing description for fcp-srvr-inst2
	Port{5503, 17}:   "fcp-srvr-inst2",       // Missing description for fcp-srvr-inst2
	Port{5504, 6}:    "fcp-cics-gw1",         // Missing description for fcp-cics-gw1
	Port{5504, 17}:   "fcp-cics-gw1",         // Missing description for fcp-cics-gw1
	Port{5505, 6}:    "checkoutdb",           // Checkout Database
	Port{5505, 17}:   "checkoutdb",           // Checkout Database
	Port{5506, 6}:    "amc",                  // Amcom Mobile Connect
	Port{5506, 17}:   "amc",                  // Amcom Mobile Connect
	Port{5507, 6}:    "psl-management",       // PowerSysLab Electrical Management
	Port{5510, 6}:    "secureidprop",         // ACE Server services
	Port{5520, 6}:    "sdlog",                // ACE Server services
	Port{5530, 6}:    "sdserv",               // ACE Server services
	Port{5540, 6}:    "sdreport",             // ACE Server services
	Port{5540, 17}:   "sdxauthd",             // ACE Server services
	Port{5550, 6}:    "sdadmind",             // cbus | ACE Server services | Model Railway control using the CBUS message protocol
	Port{5553, 6}:    "sgi-eventmond",        // SGI Eventmond Port
	Port{5553, 17}:   "sgi-eventmond",        // SGI Eventmond Port
	Port{5554, 6}:    "sgi-esphttp",          // SGI ESP HTTP
	Port{5554, 17}:   "sgi-esphttp",          // SGI ESP HTTP
	Port{5555, 6}:    "freeciv",              // personal-agent | Personal Agent
	Port{5555, 17}:   "rplay",                // Missing description for rplay
	Port{5556, 6}:    "freeciv",              // Freeciv gameplay
	Port{5556, 17}:   "freeciv",              // Freeciv gameplay
	Port{5557, 6}:    "farenet",              // Sandlab FARENET
	Port{5560, 6}:    "isqlplus",             // Oracle web enabled SQL interface (version 10g+)
	Port{5565, 6}:    "hpe-dp-bura",          // HPE Advanced BURA
	Port{5566, 6}:    "westec-connect",       // Westec Connect
	Port{5567, 6}:    "m-oap",                // dof-dps-mc-sec | Multicast Object Access Protocol | DOF Protocol Stack Multicast Secure Transport
	Port{5567, 17}:   "m-oap",                // Multicast Object Access Protocol
	Port{5568, 6}:    "sdt",                  // Session Data Transport Multicast
	Port{5568, 17}:   "sdt",                  // Session Data Transport Multicast
	Port{5569, 6}:    "rdmnet-ctrl",          // rdmnet-device | PLASA E1.33, Remote Device Management (RDM) controller status notifications | PLASA E1.33, Remote Device Management (RDM) messages
	Port{5573, 6}:    "sdmmp",                // SAS Domain Management Messaging Protocol
	Port{5573, 17}:   "sdmmp",                // SAS Domain Management Messaging Protocol
	Port{5574, 6}:    "lsi-bobcat",           // SAS IO Forwarding
	Port{5575, 6}:    "ora-oap",              // Oracle Access Protocol
	Port{5579, 6}:    "fdtracks",             // FleetDisplay Tracking Service
	Port{5580, 6}:    "tmosms0",              // T-Mobile SMS Protocol Message 0
	Port{5580, 17}:   "tmosms0",              // T-Mobile SMS Protocol Message 0
	Port{5581, 6}:    "tmosms1",              // T-Mobile SMS Protocol Message 1
	Port{5581, 17}:   "tmosms1",              // T-Mobile SMS Protocol Message 1
	Port{5582, 6}:    "fac-restore",          // T-Mobile SMS Protocol Message 3
	Port{5582, 17}:   "fac-restore",          // T-Mobile SMS Protocol Message 3
	Port{5583, 6}:    "tmo-icon-sync",        // T-Mobile SMS Protocol Message 2
	Port{5583, 17}:   "tmo-icon-sync",        // T-Mobile SMS Protocol Message 2
	Port{5584, 6}:    "bis-web",              // BeInSync-Web
	Port{5584, 17}:   "bis-web",              // BeInSync-Web
	Port{5585, 6}:    "bis-sync",             // BeInSync-sync
	Port{5585, 17}:   "bis-sync",             // BeInSync-sync
	Port{5586, 6}:    "att-mt-sms",           // Planning to send mobile terminated SMS to the specific port so that the SMS is not visible to the client
	Port{5597, 6}:    "ininmessaging",        // inin secure messaging
	Port{5597, 17}:   "ininmessaging",        // inin secure messaging
	Port{5598, 6}:    "mctfeed",              // MCT Market Data Feed
	Port{5598, 17}:   "mctfeed",              // MCT Market Data Feed
	Port{5599, 6}:    "esinstall",            // Enterprise Security Remote Install
	Port{5599, 17}:   "esinstall",            // Enterprise Security Remote Install
	Port{5600, 6}:    "esmmanager",           // Enterprise Security Manager
	Port{5600, 17}:   "esmmanager",           // Enterprise Security Manager
	Port{5601, 6}:    "esmagent",             // Enterprise Security Agent
	Port{5601, 17}:   "esmagent",             // Enterprise Security Agent
	Port{5602, 6}:    "a1-msc",               // Missing description for a1-msc
	Port{5602, 17}:   "a1-msc",               // A1-MSC
	Port{5603, 6}:    "a1-bs",                // Missing description for a1-bs
	Port{5603, 17}:   "a1-bs",                // A1-BS
	Port{5604, 6}:    "a3-sdunode",           // Missing description for a3-sdunode
	Port{5604, 17}:   "a3-sdunode",           // A3-SDUNode
	Port{5605, 6}:    "a4-sdunode",           // Missing description for a4-sdunode
	Port{5605, 17}:   "a4-sdunode",           // A4-SDUNode
	Port{5618, 6}:    "efr",                  // Fiscal Registering Protocol
	Port{5627, 6}:    "ninaf",                // Node Initiated Network Association Forma
	Port{5627, 17}:   "ninaf",                // Node Initiated Network Association Forma
	Port{5628, 6}:    "htrust",               // HTrust API
	Port{5628, 17}:   "htrust",               // HTrust API
	Port{5629, 6}:    "symantec-sfdb",        // Symantec Storage Foundation for Database
	Port{5629, 17}:   "symantec-sfdb",        // Symantec Storage Foundation for Database
	Port{5630, 6}:    "precise-comm",         // PreciseCommunication
	Port{5630, 17}:   "precise-comm",         // PreciseCommunication
	Port{5631, 6}:    "pcanywheredata",       // Missing description for pcanywheredata
	Port{5631, 17}:   "pcanywheredata",       // pcANYWHEREdata
	Port{5632, 6}:    "pcanywherestat",       // Missing description for pcanywherestat
	Port{5632, 17}:   "pcanywherestat",       // Missing description for pcanywherestat
	Port{5633, 6}:    "beorl",                // BE Operations Request Listener
	Port{5633, 17}:   "beorl",                // BE Operations Request Listener
	Port{5634, 6}:    "xprtld",               // SF Message Service
	Port{5634, 17}:   "xprtld",               // SF Message Service
	Port{5635, 6}:    "sfmsso",               // SFM Authentication Subsystem
	Port{5636, 6}:    "sfm-db-server",        // SFMdb - SFM DB server
	Port{5637, 6}:    "cssc",                 // Symantec CSSC
	Port{5638, 6}:    "flcrs",                // Symantec Fingerprint Lookup and Container Reference Service
	Port{5639, 6}:    "ics",                  // Symantec Integrity Checking Service
	Port{5646, 6}:    "vfmobile",             // Ventureforth Mobile
	Port{5666, 6}:    "nrpe",                 // Nagios NRPE | Nagios Remote Plugin Executor
	Port{5670, 6}:    "filemq",               // zre-disc | ZeroMQ file publish-subscribe protocol | Local area discovery and messaging over ZeroMQ
	Port{5671, 6}:    "amqps",                // amqp protocol over TLS SSL
	Port{5671, 17}:   "amqps",                // amqp protocol over TLS SSL
	Port{5672, 132}:  "amqp",                 // Missing description for amqp
	Port{5672, 6}:    "amqp",                 // AMQP
	Port{5672, 17}:   "amqp",                 // AMQP
	Port{5673, 6}:    "jms",                  // JACL Message Server
	Port{5673, 17}:   "jms",                  // JACL Message Server
	Port{5674, 6}:    "hyperscsi-port",       // HyperSCSI Port
	Port{5674, 17}:   "hyperscsi-port",       // HyperSCSI Port
	Port{5675, 132}:  "v5ua",                 // V5UA application port
	Port{5675, 6}:    "v5ua",                 // V5UA application port
	Port{5675, 17}:   "v5ua",                 // V5UA application port
	Port{5676, 6}:    "raadmin",              // RA Administration
	Port{5676, 17}:   "raadmin",              // RA Administration
	Port{5677, 6}:    "questdb2-lnchr",       // Quest Central DB2 Launchr
	Port{5677, 17}:   "questdb2-lnchr",       // Quest Central DB2 Launchr
	Port{5678, 6}:    "rrac",                 // Remote Replication Agent Connection
	Port{5678, 17}:   "rrac",                 // Remote Replication Agent Connection
	Port{5679, 6}:    "activesync",           // dccm | Microsoft ActiveSync PDY synchronization | Direct Cable Connect Manager
	Port{5679, 17}:   "dccm",                 // Direct Cable Connect Manager
	Port{5680, 6}:    "canna",                // auriga-router | Canna (Japanese Input) | Auriga Router Service
	Port{5680, 17}:   "auriga-router",        // Auriga Router Service
	Port{5681, 6}:    "ncxcp",                // Net-coneX Control Protocol
	Port{5681, 17}:   "ncxcp",                // Net-coneX Control Protocol
	Port{5682, 6}:    "brightcore",           // BrightCore control & data transfer exchange
	Port{5682, 17}:   "brightcore",           // BrightCore control & data transfer exchange
	Port{5683, 6}:    "coap",                 // Constrained Application Protocol | Constrained Application Protocol (CoAP)
	Port{5683, 17}:   "coap",                 // Constrained Application Protocol
	Port{5684, 6}:    "coaps",                // DTLS-secured CoAP | Constrained Application Protocol (CoAP)
	Port{5684, 17}:   "coaps",                // DTLS-secured Constrained Application Protocol
	Port{5687, 6}:    "gog-multiplayer",      // GOG multiplayer game protocol
	Port{5688, 6}:    "ggz",                  // GGZ Gaming Zone
	Port{5688, 17}:   "ggz",                  // GGZ Gaming Zone
	Port{5689, 6}:    "qmvideo",              // QM video network management protocol
	Port{5689, 17}:   "qmvideo",              // QM video network management protocol
	Port{5693, 6}:    "rbsystem",             // Robert Bosch Data Transfer
	Port{5696, 6}:    "kmip",                 // Key Management Interoperability Protocol
	Port{5700, 6}:    "supportassist",        // Dell SupportAssist data center management
	Port{5705, 6}:    "storageos",            // StorageOS REST API
	Port{5713, 6}:    "proshareaudio",        // proshare conf audio
	Port{5713, 17}:   "proshareaudio",        // proshare conf audio
	Port{5714, 6}:    "prosharevideo",        // proshare conf video
	Port{5714, 17}:   "prosharevideo",        // proshare conf video
	Port{5715, 6}:    "prosharedata",         // proshare conf data
	Port{5715, 17}:   "prosharedata",         // proshare conf data
	Port{5716, 6}:    "prosharerequest",      // proshare conf request
	Port{5716, 17}:   "prosharerequest",      // proshare conf request
	Port{5717, 6}:    "prosharenotify",       // proshare conf notify
	Port{5717, 17}:   "prosharenotify",       // proshare conf notify
	Port{5718, 6}:    "dpm",                  // DPM Communication Server
	Port{5718, 17}:   "dpm",                  // DPM Communication Server
	Port{5719, 6}:    "dpm-agent",            // DPM Agent Coordinator
	Port{5719, 17}:   "dpm-agent",            // DPM Agent Coordinator
	Port{5720, 6}:    "ms-licensing",         // Missing description for ms-licensing
	Port{5720, 17}:   "ms-licensing",         // MS-Licensing
	Port{5721, 6}:    "dtpt",                 // Desktop Passthru Service
	Port{5721, 17}:   "dtpt",                 // Desktop Passthru Service
	Port{5722, 6}:    "msdfsr",               // Microsoft DFS Replication Service
	Port{5722, 17}:   "msdfsr",               // Microsoft DFS Replication Service
	Port{5723, 6}:    "omhs",                 // Operations Manager - Health Service
	Port{5723, 17}:   "omhs",                 // Operations Manager - Health Service
	Port{5724, 6}:    "omsdk",                // Operations Manager - SDK Service
	Port{5724, 17}:   "omsdk",                // Operations Manager - SDK Service
	Port{5725, 6}:    "ms-ilm",               // Microsoft Identity Lifecycle Manager
	Port{5726, 6}:    "ms-ilm-sts",           // Microsoft Lifecycle Manager Secure Token Service
	Port{5727, 6}:    "asgenf",               // ASG Event Notification Framework
	Port{5728, 6}:    "io-dist-data",         // io-dist-group | Dist. I O Comm. Service Data and Control | Dist. I O Comm. Service Group Membership
	Port{5728, 17}:   "io-dist-group",        // Dist. I O Comm. Service Group Membership
	Port{5729, 6}:    "openmail",             // Openmail User Agent Layer
	Port{5729, 17}:   "openmail",             // Openmail User Agent Layer
	Port{5730, 6}:    "unieng",               // Steltor's calendar access
	Port{5730, 17}:   "unieng",               // Steltor's calendar access
	Port{5741, 6}:    "ida-discover1",        // IDA Discover Port 1
	Port{5741, 17}:   "ida-discover1",        // IDA Discover Port 1
	Port{5742, 6}:    "ida-discover2",        // IDA Discover Port 2
	Port{5742, 17}:   "ida-discover2",        // IDA Discover Port 2
	Port{5743, 6}:    "watchdoc-pod",         // Watchdoc NetPOD Protocol
	Port{5743, 17}:   "watchdoc-pod",         // Watchdoc NetPOD Protocol
	Port{5744, 6}:    "watchdoc",             // Watchdoc Server
	Port{5744, 17}:   "watchdoc",             // Watchdoc Server
	Port{5745, 6}:    "fcopy-server",         // Missing description for fcopy-server
	Port{5745, 17}:   "fcopy-server",         // Missing description for fcopy-server
	Port{5746, 6}:    "fcopys-server",        // Missing description for fcopys-server
	Port{5746, 17}:   "fcopys-server",        // Missing description for fcopys-server
	Port{5747, 6}:    "tunatic",              // Wildbits Tunatic
	Port{5747, 17}:   "tunatic",              // Wildbits Tunatic
	Port{5748, 6}:    "tunalyzer",            // Wildbits Tunalyzer
	Port{5748, 17}:   "tunalyzer",            // Wildbits Tunalyzer
	Port{5750, 6}:    "rscd",                 // Bladelogic Agent Service
	Port{5750, 17}:   "rscd",                 // Bladelogic Agent Service
	Port{5755, 6}:    "openmailg",            // OpenMail Desk Gateway server
	Port{5755, 17}:   "openmailg",            // OpenMail Desk Gateway server
	Port{5757, 6}:    "x500ms",               // OpenMail X.500 Directory Server
	Port{5757, 17}:   "x500ms",               // OpenMail X.500 Directory Server
	Port{5766, 6}:    "openmailns",           // OpenMail NewMail Server
	Port{5766, 17}:   "openmailns",           // OpenMail NewMail Server
	Port{5767, 6}:    "s-openmail",           // OpenMail Suer Agent Layer (Secure)
	Port{5767, 17}:   "s-openmail",           // OpenMail Suer Agent Layer (Secure)
	Port{5768, 6}:    "openmailpxy",          // OpenMail CMTS Server
	Port{5768, 17}:   "openmailpxy",          // OpenMail CMTS Server
	Port{5769, 6}:    "spramsca",             // x509solutions Internal CA
	Port{5769, 17}:   "spramsca",             // x509solutions Internal CA
	Port{5770, 6}:    "spramsd",              // x509solutions Secure Data
	Port{5770, 17}:   "spramsd",              // x509solutions Secure Data
	Port{5771, 6}:    "netagent",             // Missing description for netagent
	Port{5771, 17}:   "netagent",             // NetAgent
	Port{5777, 6}:    "dali-port",            // DALI Port
	Port{5777, 17}:   "dali-port",            // DALI Port
	Port{5780, 6}:    "vts-rpc",              // Visual Tag System RPC
	Port{5781, 6}:    "3par-evts",            // 3PAR Event Reporting Service
	Port{5781, 17}:   "3par-evts",            // 3PAR Event Reporting Service
	Port{5782, 6}:    "3par-mgmt",            // 3PAR Management Service
	Port{5782, 17}:   "3par-mgmt",            // 3PAR Management Service
	Port{5783, 6}:    "3par-mgmt-ssl",        // 3PAR Management Service with SSL
	Port{5783, 17}:   "3par-mgmt-ssl",        // 3PAR Management Service with SSL
	Port{5784, 6}:    "ibar",                 // Cisco Interbox Application Redundancy
	Port{5784, 17}:   "ibar",                 // Cisco Interbox Application Redundancy
	Port{5785, 6}:    "3par-rcopy",           // 3PAR Inform Remote Copy
	Port{5785, 17}:   "3par-rcopy",           // 3PAR Inform Remote Copy
	Port{5786, 6}:    "cisco-redu",           // redundancy notification
	Port{5786, 17}:   "cisco-redu",           // redundancy notification
	Port{5787, 6}:    "waascluster",          // Cisco WAAS Cluster Protocol
	Port{5793, 6}:    "xtreamx",              // XtreamX Supervised Peer message
	Port{5793, 17}:   "xtreamx",              // XtreamX Supervised Peer message
	Port{5794, 6}:    "spdp",                 // Simple Peered Discovery Protocol
	Port{5794, 17}:   "spdp",                 // Simple Peered Discovery Protocol
	Port{5800, 6}:    "vnc-http",             // Virtual Network Computer HTTP Access, display 0
	Port{5801, 6}:    "vnc-http-1",           // Virtual Network Computer HTTP Access, display 1
	Port{5802, 6}:    "vnc-http-2",           // Virtual Network Computer HTTP Access, display 2
	Port{5803, 6}:    "vnc-http-3",           // Virtual Network Computer HTTP Access, display 3
	Port{5813, 6}:    "icmpd",                // Missing description for icmpd
	Port{5813, 17}:   "icmpd",                // ICMPD
	Port{5814, 6}:    "spt-automation",       // Support Automation
	Port{5814, 17}:   "spt-automation",       // Support Automation
	Port{5841, 6}:    "shiprush-d-ch",        // Z-firm ShipRush interface for web access and bidirectional data
	Port{5842, 6}:    "reversion",            // Reversion Backup Restore
	Port{5859, 6}:    "wherehoo",             // Missing description for wherehoo
	Port{5859, 17}:   "wherehoo",             // WHEREHOO
	Port{5863, 6}:    "ppsuitemsg",           // PlanetPress Suite Messeng
	Port{5863, 17}:   "ppsuitemsg",           // PlanetPress Suite Messeng
	Port{5868, 6}:    "diameters",            // Diameter over TLS TCP | Diameter over DTLS SCTP
	Port{5883, 6}:    "jute",                 // Javascript Unit Test Environment
	Port{5900, 6}:    "vnc",                  // rfb | Virtual Network Computer display 0 | Remote Framebuffer
	Port{5900, 17}:   "rfb",                  // Remote Framebuffer
	Port{5901, 6}:    "vnc-1",                // Virtual Network Computer display 1
	Port{5902, 6}:    "vnc-2",                // Virtual Network Computer display 2
	Port{5903, 6}:    "vnc-3",                // Virtual Network Computer display 3
	Port{5910, 6}:    "cm",                   // Context Management
	Port{5910, 17}:   "cm",                   // Context Management
	Port{5911, 6}:    "cpdlc",                // Controller Pilot Data Link Communication
	Port{5911, 17}:   "cpdlc",                // Controller Pilot Data Link Communication
	Port{5912, 6}:    "fis",                  // Flight Information Services
	Port{5912, 17}:   "fis",                  // Flight Information Services
	Port{5913, 6}:    "ads-c",                // Automatic Dependent Surveillance
	Port{5913, 17}:   "ads-c",                // Automatic Dependent Surveillance
	Port{5938, 6}:    "teamviewer",           // teamviewer - http:  www.teamviewer.com en help 334-Which-ports-are-used-by-TeamViewer.aspx
	Port{5963, 6}:    "indy",                 // Indy Application Server
	Port{5963, 17}:   "indy",                 // Indy Application Server
	Port{5968, 6}:    "mppolicy-v5",          // Missing description for mppolicy-v5
	Port{5968, 17}:   "mppolicy-v5",          // Missing description for mppolicy-v5
	Port{5969, 6}:    "mppolicy-mgr",         // Missing description for mppolicy-mgr
	Port{5969, 17}:   "mppolicy-mgr",         // Missing description for mppolicy-mgr
	Port{5977, 6}:    "ncd-pref-tcp",         // NCD preferences tcp port
	Port{5978, 6}:    "ncd-diag-tcp",         // NCD diagnostic tcp port
	Port{5979, 6}:    "ncd-conf-tcp",         // NCD configuration tcp port
	Port{5984, 6}:    "couchdb",              // Missing description for couchdb
	Port{5984, 17}:   "couchdb",              // CouchDB
	Port{5985, 6}:    "wsman",                // WBEM WS-Management HTTP
	Port{5985, 17}:   "wsman",                // WBEM WS-Management HTTP
	Port{5986, 6}:    "wsmans",               // WBEM WS-Management HTTP over TLS SSL
	Port{5986, 17}:   "wsmans",               // WBEM WS-Management HTTP over TLS SSL
	Port{5987, 6}:    "wbem-rmi",             // WBEM RMI
	Port{5987, 17}:   "wbem-rmi",             // WBEM RMI
	Port{5988, 6}:    "wbem-http",            // WBEM CIM-XML (HTTP)
	Port{5988, 17}:   "wbem-http",            // WBEM CIM-XML (HTTP)
	Port{5989, 6}:    "wbem-https",           // WBEM CIM-XML (HTTPS)
	Port{5989, 17}:   "wbem-https",           // WBEM CIM-XML (HTTPS)
	Port{5990, 6}:    "wbem-exp-https",       // WBEM Export HTTPS
	Port{5990, 17}:   "wbem-exp-https",       // WBEM Export HTTPS
	Port{5991, 6}:    "nuxsl",                // Missing description for nuxsl
	Port{5991, 17}:   "nuxsl",                // NUXSL
	Port{5992, 6}:    "consul-insight",       // Consul InSight Security
	Port{5992, 17}:   "consul-insight",       // Consul InSight Security
	Port{5993, 6}:    "cim-rs",               // DMTF WBEM CIM REST
	Port{5997, 6}:    "ncd-pref",             // NCD preferences telnet port
	Port{5998, 6}:    "ncd-diag",             // NCD diagnostic telnet port
	Port{5999, 6}:    "ncd-conf",             // cvsup | NCD configuration telnet port | CVSup
	Port{5999, 17}:   "cvsup",                // CVSup
	Port{6000, 6}:    "X11",                  // X Window server
	Port{6000, 17}:   "X11",                  // Missing description for X11
	Port{6001, 6}:    "X11:1",                // X Window server
	Port{6001, 17}:   "X11:1",                // Missing description for X11:1
	Port{6002, 6}:    "X11:2",                // X Window server
	Port{6002, 17}:   "X11:2",                // Missing description for X11:2
	Port{6003, 6}:    "X11:3",                // X Window server
	Port{6003, 17}:   "X11:3",                // Missing description for X11:3
	Port{6004, 6}:    "X11:4",                // X Window server
	Port{6004, 17}:   "X11:4",                // Missing description for X11:4
	Port{6005, 6}:    "X11:5",                // X Window server
	Port{6005, 17}:   "X11:5",                // Missing description for X11:5
	Port{6006, 6}:    "X11:6",                // X Window server
	Port{6006, 17}:   "X11:6",                // Missing description for X11:6
	Port{6007, 6}:    "X11:7",                // X Window server
	Port{6007, 17}:   "X11:7",                // Missing description for X11:7
	Port{6008, 6}:    "X11:8",                // X Window server
	Port{6008, 17}:   "X11:8",                // Missing description for X11:8
	Port{6009, 6}:    "X11:9",                // X Window server
	Port{6009, 17}:   "X11:9",                // Missing description for X11:9
	Port{6010, 6}:    "x11",                  // X Window System
	Port{6010, 17}:   "x11",                  // X Window System
	Port{6011, 6}:    "x11",                  // X Window System
	Port{6011, 17}:   "x11",                  // X Window System
	Port{6012, 6}:    "x11",                  // X Window System
	Port{6012, 17}:   "x11",                  // X Window System
	Port{6013, 6}:    "x11",                  // X Window System
	Port{6013, 17}:   "x11",                  // X Window System
	Port{6014, 6}:    "x11",                  // X Window System
	Port{6014, 17}:   "x11",                  // X Window System
	Port{6015, 6}:    "x11",                  // X Window System
	Port{6015, 17}:   "x11",                  // X Window System
	Port{6016, 6}:    "x11",                  // X Window System
	Port{6016, 17}:   "x11",                  // X Window System
	Port{6017, 6}:    "xmail-ctrl",           // XMail CTRL server
	Port{6017, 17}:   "x11",                  // X Window System
	Port{6018, 6}:    "x11",                  // X Window System
	Port{6018, 17}:   "x11",                  // X Window System
	Port{6019, 6}:    "x11",                  // X Window System
	Port{6019, 17}:   "x11",                  // X Window System
	Port{6020, 6}:    "x11",                  // X Window System
	Port{6020, 17}:   "x11",                  // X Window System
	Port{6021, 6}:    "x11",                  // X Window System
	Port{6021, 17}:   "x11",                  // X Window System
	Port{6022, 6}:    "x11",                  // X Window System
	Port{6022, 17}:   "x11",                  // X Window System
	Port{6023, 6}:    "x11",                  // X Window System
	Port{6023, 17}:   "x11",                  // X Window System
	Port{6024, 6}:    "x11",                  // X Window System
	Port{6024, 17}:   "x11",                  // X Window System
	Port{6025, 6}:    "x11",                  // X Window System
	Port{6025, 17}:   "x11",                  // X Window System
	Port{6026, 6}:    "x11",                  // X Window System
	Port{6026, 17}:   "x11",                  // X Window System
	Port{6027, 6}:    "x11",                  // X Window System
	Port{6027, 17}:   "x11",                  // X Window System
	Port{6028, 6}:    "x11",                  // X Window System
	Port{6028, 17}:   "x11",                  // X Window System
	Port{6029, 6}:    "x11",                  // X Window System
	Port{6029, 17}:   "x11",                  // X Window System
	Port{6030, 6}:    "x11",                  // X Window System
	Port{6030, 17}:   "x11",                  // X Window System
	Port{6031, 6}:    "x11",                  // X Window System
	Port{6031, 17}:   "x11",                  // X Window System
	Port{6032, 6}:    "x11",                  // X Window System
	Port{6032, 17}:   "x11",                  // X Window System
	Port{6033, 6}:    "x11",                  // X Window System
	Port{6033, 17}:   "x11",                  // X Window System
	Port{6034, 6}:    "x11",                  // X Window System
	Port{6034, 17}:   "x11",                  // X Window System
	Port{6035, 6}:    "x11",                  // X Window System
	Port{6035, 17}:   "x11",                  // X Window System
	Port{6036, 6}:    "x11",                  // X Window System
	Port{6036, 17}:   "x11",                  // X Window System
	Port{6037, 6}:    "x11",                  // X Window System
	Port{6037, 17}:   "x11",                  // X Window System
	Port{6038, 6}:    "x11",                  // X Window System
	Port{6038, 17}:   "x11",                  // X Window System
	Port{6039, 6}:    "x11",                  // X Window System
	Port{6039, 17}:   "x11",                  // X Window System
	Port{6040, 6}:    "x11",                  // X Window System
	Port{6040, 17}:   "x11",                  // X Window System
	Port{6041, 6}:    "x11",                  // X Window System
	Port{6041, 17}:   "x11",                  // X Window System
	Port{6042, 6}:    "x11",                  // X Window System
	Port{6042, 17}:   "x11",                  // X Window System
	Port{6043, 6}:    "x11",                  // X Window System
	Port{6043, 17}:   "x11",                  // X Window System
	Port{6044, 6}:    "x11",                  // X Window System
	Port{6044, 17}:   "x11",                  // X Window System
	Port{6045, 6}:    "x11",                  // X Window System
	Port{6045, 17}:   "x11",                  // X Window System
	Port{6046, 6}:    "x11",                  // X Window System
	Port{6046, 17}:   "x11",                  // X Window System
	Port{6047, 6}:    "x11",                  // X Window System
	Port{6047, 17}:   "x11",                  // X Window System
	Port{6048, 6}:    "x11",                  // X Window System
	Port{6048, 17}:   "x11",                  // X Window System
	Port{6049, 6}:    "x11",                  // X Window System
	Port{6049, 17}:   "x11",                  // X Window System
	Port{6050, 6}:    "arcserve",             // ARCserve agent
	Port{6050, 17}:   "x11",                  // X Window System
	Port{6051, 6}:    "x11",                  // X Window System
	Port{6051, 17}:   "x11",                  // X Window System
	Port{6052, 6}:    "x11",                  // X Window System
	Port{6052, 17}:   "x11",                  // X Window System
	Port{6053, 6}:    "x11",                  // X Window System
	Port{6053, 17}:   "x11",                  // X Window System
	Port{6054, 6}:    "x11",                  // X Window System
	Port{6054, 17}:   "x11",                  // X Window System
	Port{6055, 6}:    "x11",                  // X Window System
	Port{6055, 17}:   "x11",                  // X Window System
	Port{6056, 6}:    "x11",                  // X Window System
	Port{6056, 17}:   "x11",                  // X Window System
	Port{6057, 6}:    "x11",                  // X Window System
	Port{6057, 17}:   "x11",                  // X Window System
	Port{6058, 6}:    "x11",                  // X Window System
	Port{6058, 17}:   "x11",                  // X Window System
	Port{6059, 6}:    "X11:59",               // X Window server
	Port{6059, 17}:   "X11:59",               // Missing description for X11:59
	Port{6060, 6}:    "x11",                  // X Window System
	Port{6060, 17}:   "x11",                  // X Window System
	Port{6061, 6}:    "x11",                  // X Window System
	Port{6061, 17}:   "x11",                  // X Window System
	Port{6062, 6}:    "x11",                  // X Window System
	Port{6062, 17}:   "x11",                  // X Window System
	Port{6063, 6}:    "x11",                  // X Window System
	Port{6063, 17}:   "x11",                  // X Window System
	Port{6064, 6}:    "ndl-ahp-svc",          // Missing description for ndl-ahp-svc
	Port{6064, 17}:   "ndl-ahp-svc",          // NDL-AHP-SVC
	Port{6065, 6}:    "winpharaoh",           // Missing description for winpharaoh
	Port{6065, 17}:   "winpharaoh",           // WinPharaoh
	Port{6066, 6}:    "ewctsp",               // Missing description for ewctsp
	Port{6066, 17}:   "ewctsp",               // EWCTSP
	Port{6068, 6}:    "gsmp",                 // gsmp-ancp | GSMP ANCP
	Port{6068, 17}:   "gsmp",                 // GSMP
	Port{6069, 6}:    "trip",                 // Missing description for trip
	Port{6069, 17}:   "trip",                 // TRIP
	Port{6070, 6}:    "messageasap",          // Missing description for messageasap
	Port{6070, 17}:   "messageasap",          // Messageasap
	Port{6071, 6}:    "ssdtp",                // Missing description for ssdtp
	Port{6071, 17}:   "ssdtp",                // SSDTP
	Port{6072, 6}:    "diagnose-proc",        // Missing description for diagnose-proc
	Port{6072, 17}:   "diagnose-proc",        // DIAGNOSE-PROC
	Port{6073, 6}:    "directplay8",          // Missing description for directplay8
	Port{6073, 17}:   "directplay8",          // DirectPlay8
	Port{6074, 6}:    "max",                  // Microsoft Max
	Port{6074, 17}:   "max",                  // Microsoft Max
	Port{6075, 6}:    "dpm-acm",              // Microsoft DPM Access Control Manager
	Port{6076, 6}:    "msft-dpm-cert",        // Microsoft DPM WCF Certificates
	Port{6077, 6}:    "iconstructsrv",        // iConstruct Server
	Port{6080, 6}:    "gue",                  // Generic UDP Encapsulation
	Port{6081, 6}:    "geneve",               // Generic Network Virtualization Encapsulation (Geneve)
	Port{6082, 6}:    "p25cai",               // APCO Project 25 Common Air Interface - UDP encapsulation
	Port{6083, 6}:    "miami-bcast",          // telecomsoftware miami broadcast
	Port{6084, 6}:    "p2p-sip",              // reload-config | Peer to Peer Infrastructure Protocol | Peer to Peer Infrastructure Configuration
	Port{6085, 6}:    "konspire2b",           // konspire2b p2p network
	Port{6085, 17}:   "konspire2b",           // konspire2b p2p network
	Port{6086, 6}:    "pdtp",                 // PDTP P2P
	Port{6086, 17}:   "pdtp",                 // PDTP P2P
	Port{6087, 6}:    "ldss",                 // Local Download Sharing Service
	Port{6087, 17}:   "ldss",                 // Local Download Sharing Service
	Port{6088, 6}:    "doglms",               // doglms-notify | SuperDog License Manager | SuperDog License Manager Notifier
	Port{6099, 6}:    "raxa-mgmt",            // RAXA Management
	Port{6100, 6}:    "synchronet-db",        // Missing description for synchronet-db
	Port{6100, 17}:   "synchronet-db",        // SynchroNet-db
	Port{6101, 6}:    "backupexec",           // synchronet-rtc | Backup Exec UNIX and 95 98 ME Aent | SynchroNet-rtc
	Port{6101, 17}:   "synchronet-rtc",       // SynchroNet-rtc
	Port{6102, 6}:    "synchronet-upd",       // Missing description for synchronet-upd
	Port{6102, 17}:   "synchronet-upd",       // SynchroNet-upd
	Port{6103, 6}:    "RETS-or-BackupExec",   // rets | Backup Exec Agent Accelerator and Remote Agent also sql server and cisco works blue | RETS
	Port{6103, 17}:   "rets",                 // RETS
	Port{6104, 6}:    "dbdb",                 // Missing description for dbdb
	Port{6104, 17}:   "dbdb",                 // DBDB
	Port{6105, 6}:    "isdninfo",             // primaserver | Prima Server
	Port{6105, 17}:   "primaserver",          // Prima Server
	Port{6106, 6}:    "isdninfo",             // mpsserver | i4lmond | MPS Server
	Port{6106, 17}:   "mpsserver",            // MPS Server
	Port{6107, 6}:    "etc-control",          // ETC Control
	Port{6107, 17}:   "etc-control",          // ETC Control
	Port{6108, 6}:    "sercomm-scadmin",      // Missing description for sercomm-scadmin
	Port{6108, 17}:   "sercomm-scadmin",      // Sercomm-SCAdmin
	Port{6109, 6}:    "globecast-id",         // Missing description for globecast-id
	Port{6109, 17}:   "globecast-id",         // GLOBECAST-ID
	Port{6110, 6}:    "softcm",               // HP SoftBench CM
	Port{6110, 17}:   "softcm",               // HP SoftBench CM
	Port{6111, 6}:    "spc",                  // HP SoftBench Sub-Process Control
	Port{6111, 17}:   "spc",                  // HP SoftBench Sub-Process Control
	Port{6112, 6}:    "dtspc",                // dtspcd | CDE subprocess control | Desk-Top Sub-Process Control Daemon
	Port{6112, 17}:   "dtspcd",               // Desk-Top Sub-Process Control Daemon
	Port{6113, 6}:    "dayliteserver",        // Daylite Server
	Port{6114, 6}:    "wrspice",              // WRspice IPC Service
	Port{6115, 6}:    "xic",                  // Xic IPC Service
	Port{6116, 6}:    "xtlserv",              // XicTools License Manager Service
	Port{6117, 6}:    "daylitetouch",         // Daylite Touch Sync
	Port{6118, 6}:    "tipc",                 // Transparent Inter Process Communication
	Port{6121, 6}:    "spdy",                 // SPDY for a faster web
	Port{6122, 6}:    "bex-webadmin",         // Backup Express Web Server
	Port{6122, 17}:   "bex-webadmin",         // Backup Express Web Server
	Port{6123, 6}:    "backup-express",       // Backup Express
	Port{6123, 17}:   "backup-express",       // Backup Express
	Port{6124, 6}:    "pnbs",                 // Phlexible Network Backup Service
	Port{6124, 17}:   "pnbs",                 // Phlexible Network Backup Service
	Port{6130, 6}:    "damewaremobgtwy",      // The DameWare Mobile Gateway Service
	Port{6133, 6}:    "nbt-wol",              // New Boundary Tech WOL
	Port{6133, 17}:   "nbt-wol",              // New Boundary Tech WOL
	Port{6140, 6}:    "pulsonixnls",          // Pulsonix Network License Service
	Port{6140, 17}:   "pulsonixnls",          // Pulsonix Network License Service
	Port{6141, 6}:    "meta-corp",            // Meta Corporation License Manager
	Port{6141, 17}:   "meta-corp",            // Meta Corporation License Manager
	Port{6142, 6}:    "aspentec-lm",          // Aspen Technology License Manager
	Port{6142, 17}:   "aspentec-lm",          // Aspen Technology License Manager
	Port{6143, 6}:    "watershed-lm",         // Watershed License Manager
	Port{6143, 17}:   "watershed-lm",         // Watershed License Manager
	Port{6144, 6}:    "statsci1-lm",          // StatSci License Manager - 1
	Port{6144, 17}:   "statsci1-lm",          // StatSci License Manager - 1
	Port{6145, 6}:    "statsci2-lm",          // StatSci License Manager - 2
	Port{6145, 17}:   "statsci2-lm",          // StatSci License Manager - 2
	Port{6146, 6}:    "lonewolf-lm",          // Lone Wolf Systems License Manager
	Port{6146, 17}:   "lonewolf-lm",          // Lone Wolf Systems License Manager
	Port{6147, 6}:    "montage-lm",           // Montage License Manager
	Port{6147, 17}:   "montage-lm",           // Montage License Manager
	Port{6148, 6}:    "ricardo-lm",           // Ricardo North America License Manager
	Port{6148, 17}:   "ricardo-lm",           // Ricardo North America License Manager
	Port{6149, 6}:    "tal-pod",              // Missing description for tal-pod
	Port{6149, 17}:   "tal-pod",              // Missing description for tal-pod
	Port{6159, 6}:    "efb-aci",              // EFB Application Control Interface
	Port{6160, 6}:    "ecmp",                 // ecmp-data | Emerson Extensible Control and Management Protocol | Emerson Extensible Control and Management Protocol Data
	Port{6161, 6}:    "patrol-ism",           // PATROL Internet Srv Mgr
	Port{6161, 17}:   "patrol-ism",           // PATROL Internet Srv Mgr
	Port{6162, 6}:    "patrol-coll",          // PATROL Collector
	Port{6162, 17}:   "patrol-coll",          // PATROL Collector
	Port{6163, 6}:    "pscribe",              // Precision Scribe Cnx Port
	Port{6163, 17}:   "pscribe",              // Precision Scribe Cnx Port
	Port{6200, 6}:    "lm-x",                 // LM-X License Manager by X-Formation
	Port{6200, 17}:   "lm-x",                 // LM-X License Manager by X-Formation
	Port{6201, 6}:    "thermo-calc",          // Management of service nodes in a processing grid for thermodynamic calculations
	Port{6209, 6}:    "qmtps",                // QMTP over TLS
	Port{6222, 6}:    "radmind",              // Radmind protocol | Radmind Access Protocol
	Port{6222, 17}:   "radmind",              // Radmind Access Protocol
	Port{6241, 6}:    "jeol-nsdtp-1",         // jeol-nsddp-1 | JEOL Network Services Data Transport Protocol 1 | JEOL Network Services Dynamic Discovery Protocol 1
	Port{6241, 17}:   "jeol-nsddp-1",         // JEOL Network Services Dynamic Discovery Protocol 1
	Port{6242, 6}:    "jeol-nsdtp-2",         // jeol-nsddp-2 | JEOL Network Services Data Transport Protocol 2 | JEOL Network Services Dynamic Discovery Protocol 2
	Port{6242, 17}:   "jeol-nsddp-2",         // JEOL Network Services Dynamic Discovery Protocol 2
	Port{6243, 6}:    "jeol-nsdtp-3",         // jeol-nsddp-3 | JEOL Network Services Data Transport Protocol 3 | JEOL Network Services Dynamic Discovery Protocol 3
	Port{6243, 17}:   "jeol-nsddp-3",         // JEOL Network Services Dynamic Discovery Protocol 3
	Port{6244, 6}:    "jeol-nsdtp-4",         // jeol-nsddp-4 | JEOL Network Services Data Transport Protocol 4 | JEOL Network Services Dynamic Discovery Protocol 4
	Port{6244, 17}:   "jeol-nsddp-4",         // JEOL Network Services Dynamic Discovery Protocol 4
	Port{6251, 6}:    "tl1-raw-ssl",          // TL1 Raw Over SSL TLS
	Port{6251, 17}:   "tl1-raw-ssl",          // TL1 Raw Over SSL TLS
	Port{6252, 6}:    "tl1-ssh",              // TL1 over SSH
	Port{6252, 17}:   "tl1-ssh",              // TL1 over SSH
	Port{6253, 6}:    "crip",                 // Missing description for crip
	Port{6253, 17}:   "crip",                 // CRIP
	Port{6267, 6}:    "gld",                  // GridLAB-D User Interface
	Port{6268, 6}:    "grid",                 // Grid Authentication
	Port{6268, 17}:   "grid",                 // Grid Authentication
	Port{6269, 6}:    "grid-alt",             // Grid Authentication Alt
	Port{6269, 17}:   "grid-alt",             // Grid Authentication Alt
	Port{6300, 6}:    "bmc-grx",              // BMC GRX
	Port{6300, 17}:   "bmc-grx",              // BMC GRX
	Port{6301, 6}:    "bmc_ctd_ldap",         // bmc-ctd-ldap | BMC CONTROL-D LDAP SERVER
	Port{6301, 17}:   "bmc_ctd_ldap",         // BMC CONTROL-D LDAP SERVER
	Port{6306, 6}:    "ufmp",                 // Unified Fabric Management Protocol
	Port{6306, 17}:   "ufmp",                 // Unified Fabric Management Protocol
	Port{6315, 6}:    "scup",                 // scup-disc | Sensor Control Unit Protocol | Sensor Control Unit Protocol Discovery Protocol
	Port{6315, 17}:   "scup-disc",            // Sensor Control Unit Protocol Discovery Protocol
	Port{6316, 6}:    "abb-escp",             // Ethernet Sensor Communications Protocol
	Port{6316, 17}:   "abb-escp",             // Ethernet Sensor Communications Protocol
	Port{6317, 6}:    "nav-data-cmd",         // nav-data | Navtech Radar Sensor Data Command | Navtech Radar Sensor Data
	Port{6320, 6}:    "repsvc",               // Double-Take Replication Service
	Port{6320, 17}:   "repsvc",               // Double-Take Replication Service
	Port{6321, 6}:    "emp-server1",          // Empress Software Connectivity Server 1
	Port{6321, 17}:   "emp-server1",          // Empress Software Connectivity Server 1
	Port{6322, 6}:    "emp-server2",          // Empress Software Connectivity Server 2
	Port{6322, 17}:   "emp-server2",          // Empress Software Connectivity Server 2
	Port{6324, 6}:    "hrd-ncs",              // hrd-ns-disc | HR Device Network Configuration Service | HR Device Network service
	Port{6325, 6}:    "dt-mgmtsvc",           // Double-Take Management Service
	Port{6326, 6}:    "dt-vra",               // Double-Take Virtual Recovery Assistant
	Port{6343, 6}:    "sflow",                // sFlow traffic monitoring
	Port{6343, 17}:   "sflow",                // sFlow traffic monitoring
	Port{6344, 6}:    "streletz",             // Argus-Spectr security and fire-prevention systems service
	Port{6346, 6}:    "gnutella",             // Gnutella file sharing protocol | gnutella-svc
	Port{6346, 17}:   "gnutella",             // Gnutella file sharing protocol
	Port{6347, 6}:    "gnutella2",            // Gnutella2 file sharing protocol | gnutella-rtr
	Port{6347, 17}:   "gnutella2",            // Gnutella2 file sharing protocol
	Port{6350, 6}:    "adap",                 // App Discovery and Access Protocol
	Port{6350, 17}:   "adap",                 // App Discovery and Access Protocol
	Port{6355, 6}:    "pmcs",                 // PMCS applications
	Port{6355, 17}:   "pmcs",                 // PMCS applications
	Port{6360, 6}:    "metaedit-mu",          // MetaEdit+ Multi-User
	Port{6360, 17}:   "metaedit-mu",          // MetaEdit+ Multi-User
	Port{6363, 6}:    "ndn",                  // Named Data Networking
	Port{6370, 6}:    "metaedit-se",          // MetaEdit+ Server Administration
	Port{6370, 17}:   "metaedit-se",          // MetaEdit+ Server Administration
	Port{6379, 6}:    "redis",                // An advanced key-value cache and store
	Port{6382, 6}:    "metatude-mds",         // Metatude Dialogue Server
	Port{6382, 17}:   "metatude-mds",         // Metatude Dialogue Server
	Port{6389, 6}:    "clariion-evr01",       // Missing description for clariion-evr01
	Port{6389, 17}:   "clariion-evr01",       // Missing description for clariion-evr01
	Port{6390, 6}:    "metaedit-ws",          // MetaEdit+ WebService API
	Port{6390, 17}:   "metaedit-ws",          // MetaEdit+ WebService API
	Port{6400, 6}:    "crystalreports",       // boe-cms | Seagate Crystal Reports | Business Objects CMS contact port
	Port{6400, 17}:   "boe-cms",              // Business Objects CMS contact port
	Port{6401, 6}:    "crystalenterprise",    // Seagate Crystal Enterprise | boe-was
	Port{6401, 17}:   "boe-was",              // Missing description for boe-was
	Port{6402, 6}:    "boe-eventsrv",         // Missing description for boe-eventsrv
	Port{6402, 17}:   "boe-eventsrv",         // Missing description for boe-eventsrv
	Port{6403, 6}:    "boe-cachesvr",         // Missing description for boe-cachesvr
	Port{6403, 17}:   "boe-cachesvr",         // Missing description for boe-cachesvr
	Port{6404, 6}:    "boe-filesvr",          // Business Objects Enterprise internal server
	Port{6404, 17}:   "boe-filesvr",          // Business Objects Enterprise internal server
	Port{6405, 6}:    "boe-pagesvr",          // Business Objects Enterprise internal server
	Port{6405, 17}:   "boe-pagesvr",          // Business Objects Enterprise internal server
	Port{6406, 6}:    "boe-processsvr",       // Business Objects Enterprise internal server
	Port{6406, 17}:   "boe-processsvr",       // Business Objects Enterprise internal server
	Port{6407, 6}:    "boe-resssvr1",         // Business Objects Enterprise internal server
	Port{6407, 17}:   "boe-resssvr1",         // Business Objects Enterprise internal server
	Port{6408, 6}:    "boe-resssvr2",         // Business Objects Enterprise internal server
	Port{6408, 17}:   "boe-resssvr2",         // Business Objects Enterprise internal server
	Port{6409, 6}:    "boe-resssvr3",         // Business Objects Enterprise internal server
	Port{6409, 17}:   "boe-resssvr3",         // Business Objects Enterprise internal server
	Port{6410, 6}:    "boe-resssvr4",         // Business Objects Enterprise internal server
	Port{6410, 17}:   "boe-resssvr4",         // Business Objects Enterprise internal server
	Port{6417, 6}:    "faxcomservice",        // Faxcom Message Service
	Port{6417, 17}:   "faxcomservice",        // Faxcom Message Service
	Port{6418, 6}:    "syserverremote",       // SYserver remote commands
	Port{6419, 6}:    "svdrp",                // svdrp-disc | Simple VDR Protocol | Simple VDR Protocol Discovery
	Port{6420, 6}:    "nim-vdrshell",         // NIM_VDRShell
	Port{6420, 17}:   "nim-vdrshell",         // NIM_VDRShell
	Port{6421, 6}:    "nim-wan",              // NIM_WAN
	Port{6421, 17}:   "nim-wan",              // NIM_WAN
	Port{6432, 6}:    "pgbouncer",            // Missing description for pgbouncer
	Port{6442, 6}:    "tarp",                 // Transitory Application Request Protocol
	Port{6443, 6}:    "sun-sr-https",         // Service Registry Default HTTPS Domain
	Port{6443, 17}:   "sun-sr-https",         // Service Registry Default HTTPS Domain
	Port{6444, 6}:    "sge_qmaster",          // sge-qmaster | Grid Engine Qmaster Service
	Port{6444, 17}:   "sge_qmaster",          // Grid Engine Qmaster Service
	Port{6445, 6}:    "sge_execd",            // sge-execd | Grid Engine Execution Service
	Port{6445, 17}:   "sge_execd",            // Grid Engine Execution Service
	Port{6446, 6}:    "mysql-proxy",          // MySQL Proxy
	Port{6446, 17}:   "mysql-proxy",          // MySQL Proxy
	Port{6455, 6}:    "skip-cert-recv",       // SKIP Certificate Receive
	Port{6456, 6}:    "skip-cert-send",       // SKIP Certificate Send
	Port{6456, 17}:   "skip-cert-send",       // SKIP Certificate Send
	Port{6464, 6}:    "ieee11073-20701",      // Port assignment for medical device communication in accordance to IEEE 11073-20701
	Port{6471, 6}:    "lvision-lm",           // LVision License Manager
	Port{6471, 17}:   "lvision-lm",           // LVision License Manager
	Port{6480, 6}:    "sun-sr-http",          // Service Registry Default HTTP Domain
	Port{6480, 17}:   "sun-sr-http",          // Service Registry Default HTTP Domain
	Port{6481, 6}:    "servicetags",          // Service Tags
	Port{6481, 17}:   "servicetags",          // Service Tags
	Port{6482, 6}:    "ldoms-mgmt",           // Logical Domains Management Interface
	Port{6482, 17}:   "ldoms-mgmt",           // Logical Domains Management Interface
	Port{6483, 6}:    "SunVTS-RMI",           // SunVTS RMI
	Port{6483, 17}:   "SunVTS-RMI",           // SunVTS RMI
	Port{6484, 6}:    "sun-sr-jms",           // Service Registry Default JMS Domain
	Port{6484, 17}:   "sun-sr-jms",           // Service Registry Default JMS Domain
	Port{6485, 6}:    "sun-sr-iiop",          // Service Registry Default IIOP Domain
	Port{6485, 17}:   "sun-sr-iiop",          // Service Registry Default IIOP Domain
	Port{6486, 6}:    "sun-sr-iiops",         // Service Registry Default IIOPS Domain
	Port{6486, 17}:   "sun-sr-iiops",         // Service Registry Default IIOPS Domain
	Port{6487, 6}:    "sun-sr-iiop-aut",      // Service Registry Default IIOPAuth Domain
	Port{6487, 17}:   "sun-sr-iiop-aut",      // Service Registry Default IIOPAuth Domain
	Port{6488, 6}:    "sun-sr-jmx",           // Service Registry Default JMX Domain
	Port{6488, 17}:   "sun-sr-jmx",           // Service Registry Default JMX Domain
	Port{6489, 6}:    "sun-sr-admin",         // Service Registry Default Admin Domain
	Port{6489, 17}:   "sun-sr-admin",         // Service Registry Default Admin Domain
	Port{6500, 6}:    "boks",                 // BoKS Master
	Port{6500, 17}:   "boks",                 // BoKS Master
	Port{6501, 6}:    "boks_servc",           // boks-servc | BoKS Servc
	Port{6501, 17}:   "boks_servc",           // BoKS Servc
	Port{6502, 6}:    "netop-rc",             // boks_servm | boks-servm | NetOp Remote Control (by Danware Data A S) | BoKS Servm
	Port{6502, 17}:   "netop-rc",             // NetOp Remote Control (by Danware Data A S)
	Port{6503, 6}:    "boks_clntd",           // boks-clntd | BoKS Clntd
	Port{6503, 17}:   "boks_clntd",           // BoKS Clntd
	Port{6505, 6}:    "badm_priv",            // badm-priv | BoKS Admin Private Port
	Port{6505, 17}:   "badm_priv",            // BoKS Admin Private Port
	Port{6506, 6}:    "badm_pub",             // badm-pub | BoKS Admin Public Port
	Port{6506, 17}:   "badm_pub",             // BoKS Admin Public Port
	Port{6507, 6}:    "bdir_priv",            // bdir-priv | BoKS Dir Server, Private Port
	Port{6507, 17}:   "bdir_priv",            // BoKS Dir Server, Private Port
	Port{6508, 6}:    "bdir_pub",             // bdir-pub | BoKS Dir Server, Public Port
	Port{6508, 17}:   "bdir_pub",             // BoKS Dir Server, Public Port
	Port{6509, 6}:    "mgcs-mfp-port",        // MGCS-MFP Port
	Port{6509, 17}:   "mgcs-mfp-port",        // MGCS-MFP Port
	Port{6510, 6}:    "mcer-port",            // MCER Port
	Port{6510, 17}:   "mcer-port",            // MCER Port
	Port{6511, 6}:    "dccp-udp",             // Datagram Congestion Control Protocol Encapsulation for NAT Traversal
	Port{6513, 6}:    "netconf-tls",          // NETCONF over TLS
	Port{6514, 6}:    "syslog-tls",           // Syslog over TLS | syslog over DTLS
	Port{6514, 17}:   "syslog-tls",           // syslog over DTLS
	Port{6515, 6}:    "elipse-rec",           // Elipse RPC Protocol
	Port{6515, 17}:   "elipse-rec",           // Elipse RPC Protocol
	Port{6543, 6}:    "mythtv",               // lds-distrib | lds_distrib
	Port{6543, 17}:   "lds-distrib",          // lds_distrib
	Port{6544, 6}:    "mythtv",               // lds-dump | LDS Dump Service
	Port{6544, 17}:   "lds-dump",             // LDS Dump Service
	Port{6547, 6}:    "powerchuteplus",       // apc-6547 | APC 6547
	Port{6547, 17}:   "apc-6547",             // APC 6547
	Port{6548, 6}:    "powerchuteplus",       // apc-6548 | APC 6548
	Port{6548, 17}:   "apc-6548",             // APC 6548
	Port{6549, 6}:    "apc-6549",             // APC 6549
	Port{6549, 17}:   "powerchuteplus",       // Missing description for powerchuteplus
	Port{6550, 6}:    "fg-sysupdate",         // Missing description for fg-sysupdate
	Port{6550, 17}:   "fg-sysupdate",         // Missing description for fg-sysupdate
	Port{6551, 6}:    "sum",                  // Software Update Manager
	Port{6551, 17}:   "sum",                  // Software Update Manager
	Port{6558, 6}:    "xdsxdm",               // Missing description for xdsxdm
	Port{6558, 17}:   "xdsxdm",               // Missing description for xdsxdm
	Port{6566, 6}:    "sane-port",            // SANE Control Port
	Port{6566, 17}:   "sane-port",            // SANE Control Port
	Port{6567, 6}:    "esp",                  // eSilo Storage Protocol
	Port{6567, 17}:   "esp",                  // eSilo Storage Protocol
	Port{6568, 6}:    "canit_store",          // rp-reputation | canit-store | CanIt Storage Manager | Roaring Penguin IP Address Reputation Collection
	Port{6568, 17}:   "rp-reputation",        // Roaring Penguin IP Address Reputation Collection
	Port{6579, 6}:    "affiliate",            // Missing description for affiliate
	Port{6579, 17}:   "affiliate",            // Affiliate
	Port{6580, 6}:    "parsec-master",        // Parsec Masterserver
	Port{6580, 17}:   "parsec-master",        // Parsec Masterserver
	Port{6581, 6}:    "parsec-peer",          // Parsec Peer-to-Peer
	Port{6581, 17}:   "parsec-peer",          // Parsec Peer-to-Peer
	Port{6582, 6}:    "parsec-game",          // Parsec Gameserver
	Port{6582, 17}:   "parsec-game",          // Parsec Gameserver
	Port{6583, 6}:    "joaJewelSuite",        // JOA Jewel Suite
	Port{6583, 17}:   "joaJewelSuite",        // JOA Jewel Suite
	Port{6588, 6}:    "analogx",              // AnalogX HTTP proxy port
	Port{6600, 6}:    "mshvlm",               // Microsoft Hyper-V Live Migration
	Port{6601, 6}:    "mstmg-sstp",           // Microsoft Threat Management Gateway SSTP
	Port{6602, 6}:    "wsscomfrmwk",          // Windows WSS Communication Framework
	Port{6619, 6}:    "odette-ftps",          // ODETTE-FTP over TLS SSL
	Port{6619, 17}:   "odette-ftps",          // ODETTE-FTP over TLS SSL
	Port{6620, 6}:    "kftp-data",            // Kerberos V5 FTP Data
	Port{6620, 17}:   "kftp-data",            // Kerberos V5 FTP Data
	Port{6621, 6}:    "kftp",                 // Kerberos V5 FTP Control
	Port{6621, 17}:   "kftp",                 // Kerberos V5 FTP Control
	Port{6622, 6}:    "mcftp",                // Multicast FTP
	Port{6622, 17}:   "mcftp",                // Multicast FTP
	Port{6623, 6}:    "ktelnet",              // Kerberos V5 Telnet
	Port{6623, 17}:   "ktelnet",              // Kerberos V5 Telnet
	Port{6624, 6}:    "datascaler-db",        // DataScaler database
	Port{6625, 6}:    "datascaler-ctl",       // DataScaler control
	Port{6626, 6}:    "wago-service",         // WAGO Service and Update
	Port{6626, 17}:   "wago-service",         // WAGO Service and Update
	Port{6627, 6}:    "nexgen",               // Allied Electronics NeXGen
	Port{6627, 17}:   "nexgen",               // Allied Electronics NeXGen
	Port{6628, 6}:    "afesc-mc",             // AFE Stock Channel M C
	Port{6628, 17}:   "afesc-mc",             // AFE Stock Channel M C
	Port{6629, 6}:    "nexgen-aux",           // Secondary, (non ANDI) multi-protocol multi-function interface to the Allied ANDI-based family of forecourt controllers
	Port{6632, 6}:    "mxodbc-connect",       // eGenix mxODBC Connect
	Port{6633, 6}:    "cisco-vpath-tun",      // Cisco vPath Services Overlay
	Port{6634, 6}:    "mpls-pm",              // MPLS Performance Measurement out-of-band response
	Port{6635, 6}:    "mpls-udp",             // Encapsulate MPLS packets in UDP tunnels.
	Port{6636, 6}:    "mpls-udp-dtls",        // Encapsulate MPLS packets in UDP tunnels with DTLS.
	Port{6640, 6}:    "ovsdb",                // Open vSwitch Database protocol
	Port{6653, 6}:    "openflow",             // Missing description for openflow
	Port{6655, 6}:    "pcs-sf-ui-man",        // PC SOFT - Software factory UI manager
	Port{6656, 6}:    "emgmsg",               // Emergency Message Control Service
	Port{6657, 6}:    "palcom-disc",          // PalCom Discovery
	Port{6657, 17}:   "palcom-disc",          // PalCom Discovery
	Port{6662, 6}:    "radmind",              // Radmind protocol (deprecated)
	Port{6665, 6}:    "irc",                  // Internet Relay Chat
	Port{6665, 17}:   "ircu",                 // IRCU
	Port{6666, 6}:    "irc",                  // internet relay chat server
	Port{6666, 17}:   "ircu",                 // IRCU
	Port{6667, 6}:    "irc",                  // Internet Relay Chat
	Port{6667, 17}:   "ircu",                 // IRCU
	Port{6668, 6}:    "irc",                  // Internet Relay Chat
	Port{6668, 17}:   "ircu",                 // IRCU
	Port{6669, 6}:    "irc",                  // Internet Relay Chat
	Port{6669, 17}:   "ircu",                 // IRCU
	Port{6670, 6}:    "irc",                  // vocaltec-gold | Internet Relay Chat | Vocaltec Global Online Directory
	Port{6670, 17}:   "vocaltec-gold",        // Vocaltec Global Online Directory
	Port{6671, 6}:    "p4p-portal",           // P4P Portal Service
	Port{6671, 17}:   "p4p-portal",           // P4P Portal Service
	Port{6672, 6}:    "vision_server",        // vision-server
	Port{6672, 17}:   "vision_server",        // Missing description for vision_server
	Port{6673, 6}:    "vision_elmd",          // vision-elmd
	Port{6673, 17}:   "vision_elmd",          // Missing description for vision_elmd
	Port{6678, 6}:    "vfbp",                 // vfbp-disc | Viscount Freedom Bridge Protocol | Viscount Freedom Bridge Discovery
	Port{6679, 6}:    "osaut",                // Osorno Automation
	Port{6687, 6}:    "clever-ctrace",        // CleverView for cTrace Message Service
	Port{6688, 6}:    "clever-tcpip",         // CleverView for TCP IP Message Service
	Port{6689, 6}:    "tsa",                  // Tofino Security Appliance
	Port{6689, 17}:   "tsa",                  // Tofino Security Appliance
	Port{6690, 6}:    "cleverdetect",         // CLEVERDetect Message Service
	Port{6696, 6}:    "babel",                // Babel Routing Protocol
	Port{6697, 6}:    "ircs-u",               // Internet Relay Chat via TLS SSL
	Port{6697, 17}:   "babel",                // Babel Routing Protocol
	Port{6699, 6}:    "napster",              // Napster File (MP3) sharing  software
	Port{6700, 6}:    "carracho",             // Carracho file sharing
	Port{6701, 6}:    "carracho",             // kti-icad-srvr | Carracho file sharing | KTI ICAD Nameserver
	Port{6701, 17}:   "kti-icad-srvr",        // KTI ICAD Nameserver
	Port{6702, 6}:    "e-design-net",         // e-Design network
	Port{6702, 17}:   "e-design-net",         // e-Design network
	Port{6703, 6}:    "e-design-web",         // e-Design web
	Port{6703, 17}:   "e-design-web",         // e-Design web
	Port{6704, 132}:  "frc-hp",               // ForCES HP (High Priority) channel
	Port{6705, 132}:  "frc-mp",               // ForCES MP (Medium Priority) channel
	Port{6706, 132}:  "frc-lp",               // ForCES LP (Low priority) channel
	Port{6714, 6}:    "ibprotocol",           // Internet Backplane Protocol
	Port{6714, 17}:   "ibprotocol",           // Internet Backplane Protocol
	Port{6715, 6}:    "fibotrader-com",       // Fibotrader Communications
	Port{6715, 17}:   "fibotrader-com",       // Fibotrader Communications
	Port{6716, 6}:    "princity-agent",       // Princity Agent
	Port{6767, 6}:    "bmc-perf-agent",       // BMC PERFORM AGENT
	Port{6767, 17}:   "bmc-perf-agent",       // BMC PERFORM AGENT
	Port{6768, 6}:    "bmc-perf-mgrd",        // BMC PERFORM MGRD
	Port{6768, 17}:   "bmc-perf-mgrd",        // BMC PERFORM MGRD
	Port{6769, 6}:    "adi-gxp-srvprt",       // ADInstruments GxP Server
	Port{6769, 17}:   "adi-gxp-srvprt",       // ADInstruments GxP Server
	Port{6770, 6}:    "plysrv-http",          // PolyServe http
	Port{6770, 17}:   "plysrv-http",          // PolyServe http
	Port{6771, 6}:    "plysrv-https",         // PolyServe https
	Port{6771, 17}:   "plysrv-https",         // PolyServe https
	Port{6777, 6}:    "ntz-tracker",          // netTsunami Tracker
	Port{6778, 6}:    "ntz-p2p-storage",      // netTsunami p2p storage system
	Port{6784, 6}:    "bfd-lag",              // Bidirectional Forwarding Detection (BFD) on Link Aggregation Group (LAG) Interfaces
	Port{6785, 6}:    "dgpf-exchg",           // DGPF Individual Exchange
	Port{6785, 17}:   "dgpf-exchg",           // DGPF Individual Exchange
	Port{6786, 6}:    "smc-jmx",              // Sun Java Web Console JMX
	Port{6786, 17}:   "smc-jmx",              // Sun Java Web Console JMX
	Port{6787, 6}:    "smc-admin",            // Sun Web Console Admin
	Port{6787, 17}:   "smc-admin",            // Sun Web Console Admin
	Port{6788, 6}:    "smc-http",             // Missing description for smc-http
	Port{6788, 17}:   "smc-http",             // SMC-HTTP
	Port{6789, 6}:    "ibm-db2-admin",        // radg | smc-https | IBM DB2 | SMC-HTTPS | GSS-API for the Oracle Remote Administration Daemon
	Port{6789, 17}:   "smc-https",            // SMC-HTTPS
	Port{6790, 6}:    "hnmp",                 // Missing description for hnmp
	Port{6790, 17}:   "hnmp",                 // HNMP
	Port{6791, 6}:    "hnm",                  // Halcyon Network Manager
	Port{6791, 17}:   "hnm",                  // Halcyon Network Manager
	Port{6801, 6}:    "acnet",                // ACNET Control System Protocol
	Port{6801, 17}:   "acnet",                // ACNET Control System Protocol
	Port{6817, 6}:    "pentbox-sim",          // PenTBox Secure IM Protocol
	Port{6831, 6}:    "ambit-lm",             // Missing description for ambit-lm
	Port{6831, 17}:   "ambit-lm",             // Missing description for ambit-lm
	Port{6841, 6}:    "netmo-default",        // Netmo Default
	Port{6841, 17}:   "netmo-default",        // Netmo Default
	Port{6842, 6}:    "netmo-http",           // Netmo HTTP
	Port{6842, 17}:   "netmo-http",           // Netmo HTTP
	Port{6850, 6}:    "iccrushmore",          // Missing description for iccrushmore
	Port{6850, 17}:   "iccrushmore",          // ICCRUSHMORE
	Port{6868, 6}:    "acctopus-cc",          // acctopus-st | Acctopus Command Channel | Acctopus Status
	Port{6868, 17}:   "acctopus-st",          // Acctopus Status
	Port{6881, 6}:    "bittorrent-tracker",   // BitTorrent tracker
	Port{6888, 6}:    "muse",                 // Missing description for muse
	Port{6888, 17}:   "muse",                 // MUSE
	Port{6900, 6}:    "rtimeviewer",          // R*TIME Viewer Data Interface
	Port{6901, 6}:    "jetstream",            // Novell Jetstream messaging protocol
	Port{6935, 6}:    "ethoscan",             // EthoScan Service
	Port{6936, 6}:    "xsmsvc",               // XenSource Management Service
	Port{6936, 17}:   "xsmsvc",               // XenSource Management Service
	Port{6946, 6}:    "bioserver",            // Biometrics Server
	Port{6946, 17}:   "bioserver",            // Biometrics Server
	Port{6951, 6}:    "otlp",                 // Missing description for otlp
	Port{6951, 17}:   "otlp",                 // OTLP
	Port{6961, 6}:    "jmact3",               // Missing description for jmact3
	Port{6961, 17}:   "jmact3",               // JMACT3
	Port{6962, 6}:    "jmevt2",               // Missing description for jmevt2
	Port{6962, 17}:   "jmevt2",               // Missing description for jmevt2
	Port{6963, 6}:    "swismgr1",             // Missing description for swismgr1
	Port{6963, 17}:   "swismgr1",             // Missing description for swismgr1
	Port{6964, 6}:    "swismgr2",             // Missing description for swismgr2
	Port{6964, 17}:   "swismgr2",             // Missing description for swismgr2
	Port{6965, 6}:    "swistrap",             // Missing description for swistrap
	Port{6965, 17}:   "swistrap",             // Missing description for swistrap
	Port{6966, 6}:    "swispol",              // Missing description for swispol
	Port{6966, 17}:   "swispol",              // Missing description for swispol
	Port{6969, 6}:    "acmsoda",              // Missing description for acmsoda
	Port{6969, 17}:   "acmsoda",              // Missing description for acmsoda
	Port{6970, 6}:    "conductor",            // conductor-mpx | Conductor test coordination protocol | conductor for multiplex
	Port{6997, 6}:    "MobilitySrv",          // Mobility XE Protocol
	Port{6997, 17}:   "MobilitySrv",          // Mobility XE Protocol
	Port{6998, 6}:    "iatp-highpri",         // Missing description for iatp-highpri
	Port{6998, 17}:   "iatp-highpri",         // IATP-highPri
	Port{6999, 6}:    "iatp-normalpri",       // Missing description for iatp-normalpri
	Port{6999, 17}:   "iatp-normalpri",       // IATP-normalPri
	Port{7000, 6}:    "afs3-fileserver",      // file server itself, msdos | file server itself
	Port{7000, 17}:   "afs3-fileserver",      // file server itself
	Port{7001, 6}:    "afs3-callback",        // callbacks to cache managers
	Port{7001, 17}:   "afs3-callback",        // callbacks to cache managers
	Port{7002, 6}:    "afs3-prserver",        // users & groups database
	Port{7002, 17}:   "afs3-prserver",        // users & groups database
	Port{7003, 6}:    "afs3-vlserver",        // volume location database
	Port{7003, 17}:   "afs3-vlserver",        // volume location database
	Port{7004, 6}:    "afs3-kaserver",        // AFS Kerberos authentication service
	Port{7004, 17}:   "afs3-kaserver",        // AFS Kerberos authentication service
	Port{7005, 6}:    "afs3-volser",          // volume managment server
	Port{7005, 17}:   "afs3-volser",          // volume managment server
	Port{7006, 6}:    "afs3-errors",          // error interpretation service
	Port{7006, 17}:   "afs3-errors",          // error interpretation service
	Port{7007, 6}:    "afs3-bos",             // basic overseer process
	Port{7007, 17}:   "afs3-bos",             // basic overseer process
	Port{7008, 6}:    "afs3-update",          // server-to-server updater
	Port{7008, 17}:   "afs3-update",          // server-to-server updater
	Port{7009, 6}:    "afs3-rmtsys",          // remote cache manager service
	Port{7009, 17}:   "afs3-rmtsys",          // remote cache manager service
	Port{7010, 6}:    "ups-onlinet",          // onlinet uninterruptable power supplies
	Port{7010, 17}:   "ups-onlinet",          // onlinet uninterruptable power supplies
	Port{7011, 6}:    "talon-disc",           // Talon Discovery Port
	Port{7011, 17}:   "talon-disc",           // Talon Discovery Port
	Port{7012, 6}:    "talon-engine",         // Talon Engine
	Port{7012, 17}:   "talon-engine",         // Talon Engine
	Port{7013, 6}:    "microtalon-dis",       // Microtalon Discovery
	Port{7013, 17}:   "microtalon-dis",       // Microtalon Discovery
	Port{7014, 6}:    "microtalon-com",       // Microtalon Communications
	Port{7014, 17}:   "microtalon-com",       // Microtalon Communications
	Port{7015, 6}:    "talon-webserver",      // Talon Webserver
	Port{7015, 17}:   "talon-webserver",      // Talon Webserver
	Port{7016, 6}:    "spg",                  // SPG Controls Carrier
	Port{7017, 6}:    "grasp",                // GeneRic Autonomic Signaling Protocol (TEMPORARY - registered 2017-04-28, expires 2018-04-28) | GeneRic Autonomic Signaling Protocol
	Port{7018, 6}:    "fisa-svc",             // FISA Service
	Port{7019, 6}:    "doceri-ctl",           // doceri-view | doceri drawing service control | doceri drawing service screen view
	Port{7020, 6}:    "dpserve",              // DP Serve
	Port{7020, 17}:   "dpserve",              // DP Serve
	Port{7021, 6}:    "dpserveadmin",         // DP Serve Admin
	Port{7021, 17}:   "dpserveadmin",         // DP Serve Admin
	Port{7022, 6}:    "ctdp",                 // CT Discovery Protocol
	Port{7022, 17}:   "ctdp",                 // CT Discovery Protocol
	Port{7023, 6}:    "ct2nmcs",              // Comtech T2 NMCS
	Port{7023, 17}:   "ct2nmcs",              // Comtech T2 NMCS
	Port{7024, 6}:    "vmsvc",                // Vormetric service
	Port{7024, 17}:   "vmsvc",                // Vormetric service
	Port{7025, 6}:    "vmsvc-2",              // Vormetric Service II
	Port{7025, 17}:   "vmsvc-2",              // Vormetric Service II
	Port{7030, 6}:    "op-probe",             // ObjectPlanet probe
	Port{7030, 17}:   "op-probe",             // ObjectPlanet probe
	Port{7031, 6}:    "iposplanet",           // IPOSPLANET retailing multi devices protocol
	Port{7040, 6}:    "quest-disc",           // Quest application level network service discovery
	Port{7070, 6}:    "realserver",           // arcp | ARCP
	Port{7070, 17}:   "arcp",                 // ARCP
	Port{7071, 6}:    "iwg1",                 // IWGADTS Aircraft Housekeeping Message
	Port{7071, 17}:   "iwg1",                 // IWGADTS Aircraft Housekeeping Message
	Port{7073, 6}:    "martalk",              // MarTalk protocol
	Port{7080, 6}:    "empowerid",            // EmpowerID Communication
	Port{7080, 17}:   "empowerid",            // EmpowerID Communication
	Port{7088, 6}:    "zixi-transport",       // Zixi live video transport protocol
	Port{7095, 6}:    "jdp-disc",             // Java Discovery Protocol
	Port{7099, 6}:    "lazy-ptop",            // Missing description for lazy-ptop
	Port{7099, 17}:   "lazy-ptop",            // Missing description for lazy-ptop
	Port{7100, 6}:    "font-service",         // X Font Service
	Port{7100, 17}:   "font-service",         // X Font Service
	Port{7101, 6}:    "elcn",                 // Embedded Light Control Network
	Port{7101, 17}:   "elcn",                 // Embedded Light Control Network
	Port{7107, 6}:    "aes-x170",             // Missing description for aes-x170
	Port{7117, 6}:    "rothaga",              // Encrypted chat and file transfer service
	Port{7121, 6}:    "virprot-lm",           // Virtual Prototypes License Manager
	Port{7121, 17}:   "virprot-lm",           // Virtual Prototypes License Manager
	Port{7128, 6}:    "scenidm",              // intelligent data manager
	Port{7128, 17}:   "scenidm",              // intelligent data manager
	Port{7129, 6}:    "scenccs",              // Catalog Content Search
	Port{7129, 17}:   "scenccs",              // Catalog Content Search
	Port{7161, 6}:    "cabsm-comm",           // CA BSM Comm
	Port{7161, 17}:   "cabsm-comm",           // CA BSM Comm
	Port{7162, 6}:    "caistoragemgr",        // CA Storage Manager
	Port{7162, 17}:   "caistoragemgr",        // CA Storage Manager
	Port{7163, 6}:    "cacsambroker",         // CA Connection Broker
	Port{7163, 17}:   "cacsambroker",         // CA Connection Broker
	Port{7164, 6}:    "fsr",                  // File System Repository Agent
	Port{7164, 17}:   "fsr",                  // File System Repository Agent
	Port{7165, 6}:    "doc-server",           // Document WCF Server
	Port{7165, 17}:   "doc-server",           // Document WCF Server
	Port{7166, 6}:    "aruba-server",         // Aruba eDiscovery Server
	Port{7166, 17}:   "aruba-server",         // Aruba eDiscovery Server
	Port{7167, 6}:    "casrmagent",           // CA SRM Agent
	Port{7168, 6}:    "cnckadserver",         // cncKadServer DB & Inventory Services
	Port{7169, 6}:    "ccag-pib",             // Consequor Consulting Process Integration Bridge
	Port{7169, 17}:   "ccag-pib",             // Consequor Consulting Process Integration Bridge
	Port{7170, 6}:    "nsrp",                 // Adaptive Name Service Resolution
	Port{7170, 17}:   "nsrp",                 // Adaptive Name Service Resolution
	Port{7171, 6}:    "drm-production",       // Discovery and Retention Mgt Production
	Port{7171, 17}:   "drm-production",       // Discovery and Retention Mgt Production
	Port{7172, 6}:    "metalbend",            // Port used for MetalBend programmable interface
	Port{7173, 6}:    "zsecure",              // zSecure Server
	Port{7174, 6}:    "clutild",              // Missing description for clutild
	Port{7174, 17}:   "clutild",              // Clutild
	Port{7181, 6}:    "janus-disc",           // Janus Guidewire Enterprise Discovery Service Bus
	Port{7200, 6}:    "fodms",                // FODMS FLIP
	Port{7200, 17}:   "fodms",                // FODMS FLIP
	Port{7201, 6}:    "dlip",                 // Missing description for dlip
	Port{7201, 17}:   "dlip",                 // Missing description for dlip
	Port{7202, 6}:    "pon-ictp",             // Inter-Channel Termination Protocol (ICTP) for multi-wavelength PON (Passive Optical Network) systems
	Port{7215, 6}:    "PS-Server",            // Communication ports for PaperStream Server services
	Port{7216, 6}:    "PS-Capture-Pro",       // PaperStream Capture Professional
	Port{7227, 6}:    "ramp",                 // Registry A & M Protocol | Registry A $ M Protocol
	Port{7227, 17}:   "ramp",                 // Registry A $ M Protocol
	Port{7228, 6}:    "citrixupp",            // Citrix Universal Printing Port
	Port{7229, 6}:    "citrixuppg",           // Citrix UPP Gateway
	Port{7235, 6}:    "aspcoordination",      // ASP Coordination Protocol
	Port{7236, 6}:    "display",              // Wi-Fi Alliance Wi-Fi Display Protocol
	Port{7237, 6}:    "pads",                 // PADS (Public Area Display System) Server
	Port{7244, 6}:    "frc-hicp",             // frc-hicp-disc | FrontRow Calypso Human Interface Control Protocol
	Port{7262, 6}:    "cnap",                 // Calypso Network Access Protocol
	Port{7262, 17}:   "cnap",                 // Calypso Network Access Protocol
	Port{7272, 6}:    "watchme-7272",         // WatchMe Monitoring 7272
	Port{7272, 17}:   "watchme-7272",         // WatchMe Monitoring 7272
	Port{7273, 6}:    "openmanage",           // oma-rlp | Dell OpenManage | OMA Roaming Location
	Port{7273, 17}:   "oma-rlp",              // OMA Roaming Location
	Port{7274, 6}:    "oma-rlp-s",            // OMA Roaming Location SEC
	Port{7274, 17}:   "oma-rlp-s",            // OMA Roaming Location SEC
	Port{7275, 6}:    "oma-ulp",              // OMA UserPlane Location
	Port{7275, 17}:   "oma-ulp",              // OMA UserPlane Location
	Port{7276, 6}:    "oma-ilp",              // OMA Internal Location Protocol
	Port{7276, 17}:   "oma-ilp",              // OMA Internal Location Protocol
	Port{7277, 6}:    "oma-ilp-s",            // OMA Internal Location Secure Protocol
	Port{7277, 17}:   "oma-ilp-s",            // OMA Internal Location Secure Protocol
	Port{7278, 6}:    "oma-dcdocbs",          // OMA Dynamic Content Delivery over CBS
	Port{7278, 17}:   "oma-dcdocbs",          // OMA Dynamic Content Delivery over CBS
	Port{7279, 6}:    "ctxlic",               // Citrix Licensing
	Port{7279, 17}:   "ctxlic",               // Citrix Licensing
	Port{7280, 6}:    "itactionserver1",      // ITACTIONSERVER 1
	Port{7280, 17}:   "itactionserver1",      // ITACTIONSERVER 1
	Port{7281, 6}:    "itactionserver2",      // ITACTIONSERVER 2
	Port{7281, 17}:   "itactionserver2",      // ITACTIONSERVER 2
	Port{7282, 6}:    "mzca-action",          // mzca-alert | eventACTION ussACTION (MZCA) server | eventACTION ussACTION (MZCA) alert
	Port{7282, 17}:   "mzca-alert",           // eventACTION ussACTION (MZCA) alert
	Port{7283, 6}:    "genstat",              // General Statistics Rendezvous Protocol
	Port{7300, 6}:    "swx",                  // The Swiss Exchange
	Port{7300, 17}:   "swx",                  // The Swiss Exchange
	Port{7301, 6}:    "swx",                  // The Swiss Exchange
	Port{7301, 17}:   "swx",                  // The Swiss Exchange
	Port{7302, 6}:    "swx",                  // The Swiss Exchange
	Port{7302, 17}:   "swx",                  // The Swiss Exchange
	Port{7303, 6}:    "swx",                  // The Swiss Exchange
	Port{7303, 17}:   "swx",                  // The Swiss Exchange
	Port{7304, 6}:    "swx",                  // The Swiss Exchange
	Port{7304, 17}:   "swx",                  // The Swiss Exchange
	Port{7305, 6}:    "swx",                  // The Swiss Exchange
	Port{7305, 17}:   "swx",                  // The Swiss Exchange
	Port{7306, 6}:    "swx",                  // The Swiss Exchange
	Port{7306, 17}:   "swx",                  // The Swiss Exchange
	Port{7307, 6}:    "swx",                  // The Swiss Exchange
	Port{7307, 17}:   "swx",                  // The Swiss Exchange
	Port{7308, 6}:    "swx",                  // The Swiss Exchange
	Port{7308, 17}:   "swx",                  // The Swiss Exchange
	Port{7309, 6}:    "swx",                  // The Swiss Exchange
	Port{7309, 17}:   "swx",                  // The Swiss Exchange
	Port{7310, 6}:    "swx",                  // The Swiss Exchange
	Port{7310, 17}:   "swx",                  // The Swiss Exchange
	Port{7311, 6}:    "swx",                  // The Swiss Exchange
	Port{7311, 17}:   "swx",                  // The Swiss Exchange
	Port{7312, 6}:    "swx",                  // The Swiss Exchange
	Port{7312, 17}:   "swx",                  // The Swiss Exchange
	Port{7313, 6}:    "swx",                  // The Swiss Exchange
	Port{7313, 17}:   "swx",                  // The Swiss Exchange
	Port{7314, 6}:    "swx",                  // The Swiss Exchange
	Port{7314, 17}:   "swx",                  // The Swiss Exchange
	Port{7315, 6}:    "swx",                  // The Swiss Exchange
	Port{7315, 17}:   "swx",                  // The Swiss Exchange
	Port{7316, 6}:    "swx",                  // The Swiss Exchange
	Port{7316, 17}:   "swx",                  // The Swiss Exchange
	Port{7317, 6}:    "swx",                  // The Swiss Exchange
	Port{7317, 17}:   "swx",                  // The Swiss Exchange
	Port{7318, 6}:    "swx",                  // The Swiss Exchange
	Port{7318, 17}:   "swx",                  // The Swiss Exchange
	Port{7319, 6}:    "swx",                  // The Swiss Exchange
	Port{7319, 17}:   "swx",                  // The Swiss Exchange
	Port{7320, 6}:    "swx",                  // The Swiss Exchange
	Port{7320, 17}:   "swx",                  // The Swiss Exchange
	Port{7321, 6}:    "swx",                  // The Swiss Exchange
	Port{7321, 17}:   "swx",                  // The Swiss Exchange
	Port{7322, 6}:    "swx",                  // The Swiss Exchange
	Port{7322, 17}:   "swx",                  // The Swiss Exchange
	Port{7323, 6}:    "swx",                  // The Swiss Exchange
	Port{7323, 17}:   "swx",                  // The Swiss Exchange
	Port{7324, 6}:    "swx",                  // The Swiss Exchange
	Port{7324, 17}:   "swx",                  // The Swiss Exchange
	Port{7325, 6}:    "swx",                  // The Swiss Exchange
	Port{7325, 17}:   "swx",                  // The Swiss Exchange
	Port{7326, 6}:    "icb",                  // Internet Citizen's Band
	Port{7326, 17}:   "swx",                  // The Swiss Exchange
	Port{7327, 6}:    "swx",                  // The Swiss Exchange
	Port{7327, 17}:   "swx",                  // The Swiss Exchange
	Port{7328, 6}:    "swx",                  // The Swiss Exchange
	Port{7328, 17}:   "swx",                  // The Swiss Exchange
	Port{7329, 6}:    "swx",                  // The Swiss Exchange
	Port{7329, 17}:   "swx",                  // The Swiss Exchange
	Port{7330, 6}:    "swx",                  // The Swiss Exchange
	Port{7330, 17}:   "swx",                  // The Swiss Exchange
	Port{7331, 6}:    "swx",                  // The Swiss Exchange
	Port{7331, 17}:   "swx",                  // The Swiss Exchange
	Port{7332, 6}:    "swx",                  // The Swiss Exchange
	Port{7332, 17}:   "swx",                  // The Swiss Exchange
	Port{7333, 6}:    "swx",                  // The Swiss Exchange
	Port{7333, 17}:   "swx",                  // The Swiss Exchange
	Port{7334, 6}:    "swx",                  // The Swiss Exchange
	Port{7334, 17}:   "swx",                  // The Swiss Exchange
	Port{7335, 6}:    "swx",                  // The Swiss Exchange
	Port{7335, 17}:   "swx",                  // The Swiss Exchange
	Port{7336, 6}:    "swx",                  // The Swiss Exchange
	Port{7336, 17}:   "swx",                  // The Swiss Exchange
	Port{7337, 6}:    "swx",                  // The Swiss Exchange
	Port{7337, 17}:   "swx",                  // The Swiss Exchange
	Port{7338, 6}:    "swx",                  // The Swiss Exchange
	Port{7338, 17}:   "swx",                  // The Swiss Exchange
	Port{7339, 6}:    "swx",                  // The Swiss Exchange
	Port{7339, 17}:   "swx",                  // The Swiss Exchange
	Port{7340, 6}:    "swx",                  // The Swiss Exchange
	Port{7340, 17}:   "swx",                  // The Swiss Exchange
	Port{7341, 6}:    "swx",                  // The Swiss Exchange
	Port{7341, 17}:   "swx",                  // The Swiss Exchange
	Port{7342, 6}:    "swx",                  // The Swiss Exchange
	Port{7342, 17}:   "swx",                  // The Swiss Exchange
	Port{7343, 6}:    "swx",                  // The Swiss Exchange
	Port{7343, 17}:   "swx",                  // The Swiss Exchange
	Port{7344, 6}:    "swx",                  // The Swiss Exchange
	Port{7344, 17}:   "swx",                  // The Swiss Exchange
	Port{7345, 6}:    "swx",                  // The Swiss Exchange
	Port{7345, 17}:   "swx",                  // The Swiss Exchange
	Port{7346, 6}:    "swx",                  // The Swiss Exchange
	Port{7346, 17}:   "swx",                  // The Swiss Exchange
	Port{7347, 6}:    "swx",                  // The Swiss Exchange
	Port{7347, 17}:   "swx",                  // The Swiss Exchange
	Port{7348, 6}:    "swx",                  // The Swiss Exchange
	Port{7348, 17}:   "swx",                  // The Swiss Exchange
	Port{7349, 6}:    "swx",                  // The Swiss Exchange
	Port{7349, 17}:   "swx",                  // The Swiss Exchange
	Port{7350, 6}:    "swx",                  // The Swiss Exchange
	Port{7350, 17}:   "swx",                  // The Swiss Exchange
	Port{7351, 6}:    "swx",                  // The Swiss Exchange
	Port{7351, 17}:   "swx",                  // The Swiss Exchange
	Port{7352, 6}:    "swx",                  // The Swiss Exchange
	Port{7352, 17}:   "swx",                  // The Swiss Exchange
	Port{7353, 6}:    "swx",                  // The Swiss Exchange
	Port{7353, 17}:   "swx",                  // The Swiss Exchange
	Port{7354, 6}:    "swx",                  // The Swiss Exchange
	Port{7354, 17}:   "swx",                  // The Swiss Exchange
	Port{7355, 6}:    "swx",                  // The Swiss Exchange
	Port{7355, 17}:   "swx",                  // The Swiss Exchange
	Port{7356, 6}:    "swx",                  // The Swiss Exchange
	Port{7356, 17}:   "swx",                  // The Swiss Exchange
	Port{7357, 6}:    "swx",                  // The Swiss Exchange
	Port{7357, 17}:   "swx",                  // The Swiss Exchange
	Port{7358, 6}:    "swx",                  // The Swiss Exchange
	Port{7358, 17}:   "swx",                  // The Swiss Exchange
	Port{7359, 6}:    "swx",                  // The Swiss Exchange
	Port{7359, 17}:   "swx",                  // The Swiss Exchange
	Port{7365, 6}:    "lcm-server",           // LifeKeeper Communications
	Port{7365, 17}:   "lcm-server",           // LifeKeeper Communications
	Port{7391, 6}:    "mindfilesys",          // mind-file system server
	Port{7391, 17}:   "mindfilesys",          // mind-file system server
	Port{7392, 6}:    "mrssrendezvous",       // mrss-rendezvous server
	Port{7392, 17}:   "mrssrendezvous",       // mrss-rendezvous server
	Port{7393, 6}:    "nfoldman",             // nFoldMan Remote Publish
	Port{7393, 17}:   "nfoldman",             // nFoldMan Remote Publish
	Port{7394, 6}:    "fse",                  // File system export of backup images
	Port{7394, 17}:   "fse",                  // File system export of backup images
	Port{7395, 6}:    "winqedit",             // Missing description for winqedit
	Port{7395, 17}:   "winqedit",             // Missing description for winqedit
	Port{7397, 6}:    "hexarc",               // Hexarc Command Language
	Port{7397, 17}:   "hexarc",               // Hexarc Command Language
	Port{7400, 6}:    "rtps-discovery",       // RTPS Discovery
	Port{7400, 17}:   "rtps-discovery",       // RTPS Discovery
	Port{7401, 6}:    "rtps-dd-ut",           // RTPS Data-Distribution User-Traffic
	Port{7401, 17}:   "rtps-dd-ut",           // RTPS Data-Distribution User-Traffic
	Port{7402, 6}:    "rtps-dd-mt",           // RTPS Data-Distribution Meta-Traffic
	Port{7402, 17}:   "rtps-dd-mt",           // RTPS Data-Distribution Meta-Traffic
	Port{7410, 6}:    "ionixnetmon",          // Ionix Network Monitor
	Port{7410, 17}:   "ionixnetmon",          // Ionix Network Monitor
	Port{7411, 6}:    "daqstream",            // Streaming of measurement data
	Port{7420, 6}:    "ipluminary",           // Multichannel real-time lighting control
	Port{7421, 6}:    "mtportmon",            // Matisse Port Monitor
	Port{7421, 17}:   "mtportmon",            // Matisse Port Monitor
	Port{7426, 6}:    "pmdmgr",               // OpenView DM Postmaster Manager
	Port{7426, 17}:   "pmdmgr",               // OpenView DM Postmaster Manager
	Port{7427, 6}:    "oveadmgr",             // OpenView DM Event Agent Manager
	Port{7427, 17}:   "oveadmgr",             // OpenView DM Event Agent Manager
	Port{7428, 6}:    "ovladmgr",             // OpenView DM Log Agent Manager
	Port{7428, 17}:   "ovladmgr",             // OpenView DM Log Agent Manager
	Port{7429, 6}:    "opi-sock",             // OpenView DM rqt communication
	Port{7429, 17}:   "opi-sock",             // OpenView DM rqt communication
	Port{7430, 6}:    "xmpv7",                // OpenView DM xmpv7 api pipe
	Port{7430, 17}:   "xmpv7",                // OpenView DM xmpv7 api pipe
	Port{7431, 6}:    "pmd",                  // OpenView DM ovc xmpv3 api pipe
	Port{7431, 17}:   "pmd",                  // OpenView DM ovc xmpv3 api pipe
	Port{7437, 6}:    "faximum",              // Missing description for faximum
	Port{7437, 17}:   "faximum",              // Faximum
	Port{7443, 6}:    "oracleas-https",       // Oracle Application Server HTTPS
	Port{7443, 17}:   "oracleas-https",       // Oracle Application Server HTTPS
	Port{7464, 6}:    "pythonds",             // Python Documentation Server
	Port{7471, 6}:    "sttunnel",             // Stateless Transport Tunneling Protocol
	Port{7473, 6}:    "rise",                 // Rise: The Vieneo Province
	Port{7473, 17}:   "rise",                 // Rise: The Vieneo Province
	Port{7474, 6}:    "neo4j",                // Neo4j Graph Database
	Port{7478, 6}:    "openit",               // IT Asset Management
	Port{7491, 6}:    "telops-lmd",           // Missing description for telops-lmd
	Port{7491, 17}:   "telops-lmd",           // Missing description for telops-lmd
	Port{7500, 6}:    "silhouette",           // Silhouette User
	Port{7500, 17}:   "silhouette",           // Silhouette User
	Port{7501, 6}:    "ovbus",                // HP OpenView Bus Daemon
	Port{7501, 17}:   "ovbus",                // HP OpenView Bus Daemon
	Port{7508, 6}:    "adcp",                 // Automation Device Configuration Protocol
	Port{7509, 6}:    "acplt",                // ACPLT - process automation service
	Port{7510, 6}:    "ovhpas",               // HP OpenView Application Server
	Port{7510, 17}:   "ovhpas",               // HP OpenView Application Server
	Port{7511, 6}:    "pafec-lm",             // Missing description for pafec-lm
	Port{7511, 17}:   "pafec-lm",             // Missing description for pafec-lm
	Port{7542, 6}:    "saratoga",             // Saratoga Transfer Protocol
	Port{7542, 17}:   "saratoga",             // Saratoga Transfer Protocol
	Port{7543, 6}:    "atul",                 // atul server
	Port{7543, 17}:   "atul",                 // atul server
	Port{7544, 6}:    "nta-ds",               // FlowAnalyzer DisplayServer
	Port{7544, 17}:   "nta-ds",               // FlowAnalyzer DisplayServer
	Port{7545, 6}:    "nta-us",               // FlowAnalyzer UtilityServer
	Port{7545, 17}:   "nta-us",               // FlowAnalyzer UtilityServer
	Port{7546, 6}:    "cfs",                  // Cisco Fabric service
	Port{7546, 17}:   "cfs",                  // Cisco Fabric service
	Port{7547, 6}:    "cwmp",                 // DSL Forum CWMP
	Port{7547, 17}:   "cwmp",                 // DSL Forum CWMP
	Port{7548, 6}:    "tidp",                 // Threat Information Distribution Protocol
	Port{7548, 17}:   "tidp",                 // Threat Information Distribution Protocol
	Port{7549, 6}:    "nls-tl",               // Network Layer Signaling Transport Layer
	Port{7549, 17}:   "nls-tl",               // Network Layer Signaling Transport Layer
	Port{7550, 6}:    "cloudsignaling",       // Cloud Signaling Service
	Port{7551, 6}:    "controlone-con",       // ControlONE Console signaling
	Port{7560, 6}:    "sncp",                 // Sniffer Command Protocol
	Port{7560, 17}:   "sncp",                 // Sniffer Command Protocol
	Port{7563, 6}:    "cfw",                  // Control Framework
	Port{7566, 6}:    "vsi-omega",            // VSI Omega
	Port{7566, 17}:   "vsi-omega",            // VSI Omega
	Port{7569, 6}:    "dell-eql-asm",         // Dell EqualLogic Host Group Management
	Port{7570, 6}:    "aries-kfinder",        // Aries Kfinder
	Port{7570, 17}:   "aries-kfinder",        // Aries Kfinder
	Port{7574, 6}:    "coherence",            // coherence-disc | Oracle Coherence Cluster Service | Oracle Coherence Cluster discovery service
	Port{7588, 6}:    "sun-lm",               // Sun License Manager
	Port{7588, 17}:   "sun-lm",               // Sun License Manager
	Port{7597, 6}:    "qaz",                  // Quaz trojan worm
	Port{7606, 6}:    "mipi-debug",           // MIPI Alliance Debug
	Port{7624, 6}:    "indi",                 // Instrument Neutral Distributed Interface
	Port{7624, 17}:   "indi",                 // Instrument Neutral Distributed Interface
	Port{7626, 132}:  "simco",                // SImple Middlebox COnfiguration (SIMCO) | SImple Middlebox COnfiguration (SIMCO) Server
	Port{7626, 6}:    "simco",                // SImple Middlebox COnfiguration (SIMCO) Server
	Port{7627, 6}:    "soap-http",            // SOAP Service Port
	Port{7627, 17}:   "soap-http",            // SOAP Service Port
	Port{7628, 6}:    "zen-pawn",             // Primary Agent Work Notification
	Port{7628, 17}:   "zen-pawn",             // Primary Agent Work Notification
	Port{7629, 6}:    "xdas",                 // OpenXDAS Wire Protocol
	Port{7629, 17}:   "xdas",                 // OpenXDAS Wire Protocol
	Port{7630, 6}:    "hawk",                 // HA Web Konsole
	Port{7631, 6}:    "tesla-sys-msg",        // TESLA System Messaging
	Port{7633, 6}:    "pmdfmgt",              // PMDF Management
	Port{7633, 17}:   "pmdfmgt",              // PMDF Management
	Port{7634, 6}:    "hddtemp",              // A cross-platform hard disk temperature monitoring daemon
	Port{7648, 6}:    "cuseeme",              // bonjour-cuseeme
	Port{7648, 17}:   "cucme-1",              // cucme live video audio server
	Port{7649, 17}:   "cucme-2",              // cucme live video audio server
	Port{7650, 17}:   "cucme-3",              // cucme live video audio server
	Port{7651, 17}:   "cucme-4",              // cucme live video audio server
	Port{7663, 6}:    "rome",                 // Proprietary immutable distributed data storage
	Port{7672, 6}:    "imqstomp",             // iMQ STOMP Server
	Port{7673, 6}:    "imqstomps",            // iMQ STOMP Server over SSL
	Port{7674, 6}:    "imqtunnels",           // iMQ SSL tunnel
	Port{7674, 17}:   "imqtunnels",           // iMQ SSL tunnel
	Port{7675, 6}:    "imqtunnel",            // iMQ Tunnel
	Port{7675, 17}:   "imqtunnel",            // iMQ Tunnel
	Port{7676, 6}:    "imqbrokerd",           // iMQ Broker Rendezvous
	Port{7676, 17}:   "imqbrokerd",           // iMQ Broker Rendezvous
	Port{7677, 6}:    "sun-user-https",       // Sun App Server - HTTPS
	Port{7677, 17}:   "sun-user-https",       // Sun App Server - HTTPS
	Port{7680, 6}:    "pando-pub",            // Pando Media Public Distribution
	Port{7680, 17}:   "pando-pub",            // Pando Media Public Distribution
	Port{7683, 6}:    "dmt",                  // Cleondris DMT
	Port{7687, 6}:    "bolt",                 // Bolt database connection
	Port{7689, 6}:    "collaber",             // Collaber Network Service
	Port{7689, 17}:   "collaber",             // Collaber Network Service
	Port{7697, 6}:    "klio",                 // KLIO communications
	Port{7697, 17}:   "klio",                 // KLIO communications
	Port{7700, 6}:    "em7-secom",            // EM7 Secure Communications
	Port{7701, 6}:    "nfapi",                // SCF nFAPI defining MAC PHY split
	Port{7707, 6}:    "sync-em7",             // EM7 Dynamic Updates
	Port{7707, 17}:   "sync-em7",             // EM7 Dynamic Updates
	Port{7708, 6}:    "scinet",               // scientia.net
	Port{7708, 17}:   "scinet",               // scientia.net
	Port{7720, 6}:    "medimageportal",       // MedImage Portal
	Port{7720, 17}:   "medimageportal",       // MedImage Portal
	Port{7724, 6}:    "nsdeepfreezectl",      // Novell Snap-in Deep Freeze Control
	Port{7724, 17}:   "nsdeepfreezectl",      // Novell Snap-in Deep Freeze Control
	Port{7725, 6}:    "nitrogen",             // Nitrogen Service
	Port{7725, 17}:   "nitrogen",             // Nitrogen Service
	Port{7726, 6}:    "freezexservice",       // FreezeX Console Service
	Port{7726, 17}:   "freezexservice",       // FreezeX Console Service
	Port{7727, 6}:    "trident-data",         // Trident Systems Data
	Port{7727, 17}:   "trident-data",         // Trident Systems Data
	Port{7728, 6}:    "osvr",                 // Open-Source Virtual Reality
	Port{7734, 6}:    "smip",                 // Smith Protocol over IP
	Port{7734, 17}:   "smip",                 // Smith Protocol over IP
	Port{7738, 6}:    "aiagent",              // HP Enterprise Discovery Agent
	Port{7738, 17}:   "aiagent",              // HP Enterprise Discovery Agent
	Port{7741, 6}:    "scriptview",           // ScriptView Network
	Port{7741, 17}:   "scriptview",           // ScriptView Network
	Port{7742, 6}:    "msss",                 // Mugginsoft Script Server Service
	Port{7743, 6}:    "sstp-1",               // Sakura Script Transfer Protocol
	Port{7743, 17}:   "sstp-1",               // Sakura Script Transfer Protocol
	Port{7744, 6}:    "raqmon-pdu",           // RAQMON PDU
	Port{7744, 17}:   "raqmon-pdu",           // RAQMON PDU
	Port{7747, 6}:    "prgp",                 // Put Run Get Protocol
	Port{7747, 17}:   "prgp",                 // Put Run Get Protocol
	Port{7775, 6}:    "inetfs",               // A File System using TLS over a wide area network
	Port{7777, 6}:    "cbt",                  // Missing description for cbt
	Port{7777, 17}:   "cbt",                  // Missing description for cbt
	Port{7778, 6}:    "interwise",            // Missing description for interwise
	Port{7778, 17}:   "interwise",            // Interwise
	Port{7779, 6}:    "vstat",                // Missing description for vstat
	Port{7779, 17}:   "vstat",                // VSTAT
	Port{7781, 6}:    "accu-lmgr",            // Missing description for accu-lmgr
	Port{7781, 17}:   "accu-lmgr",            // Missing description for accu-lmgr
	Port{7784, 6}:    "s-bfd",                // Seamless Bidirectional Forwarding Detection (S-BFD)
	Port{7786, 6}:    "minivend",             // Missing description for minivend
	Port{7786, 17}:   "minivend",             // MINIVEND
	Port{7787, 6}:    "popup-reminders",      // Popup Reminders Receive
	Port{7787, 17}:   "popup-reminders",      // Popup Reminders Receive
	Port{7789, 6}:    "office-tools",         // Office Tools Pro Receive
	Port{7789, 17}:   "office-tools",         // Office Tools Pro Receive
	Port{7794, 6}:    "q3ade",                // Q3ADE Cluster Service
	Port{7794, 17}:   "q3ade",                // Q3ADE Cluster Service
	Port{7797, 6}:    "pnet-conn",            // Propel Connector port
	Port{7797, 17}:   "pnet-conn",            // Propel Connector port
	Port{7798, 6}:    "pnet-enc",             // Propel Encoder port
	Port{7798, 17}:   "pnet-enc",             // Propel Encoder port
	Port{7799, 6}:    "altbsdp",              // Alternate BSDP Service
	Port{7799, 17}:   "altbsdp",              // Alternate BSDP Service
	Port{7800, 6}:    "asr",                  // Apple Software Restore
	Port{7800, 17}:   "asr",                  // Apple Software Restore
	Port{7801, 6}:    "ssp-client",           // Secure Server Protocol - client
	Port{7801, 17}:   "ssp-client",           // Secure Server Protocol - client
	Port{7802, 6}:    "vns-tp",               // Virtualized Network Services Tunnel Protocol
	Port{7810, 6}:    "rbt-wanopt",           // Riverbed WAN Optimization Protocol
	Port{7810, 17}:   "rbt-wanopt",           // Riverbed WAN Optimization Protocol
	Port{7845, 6}:    "apc-7845",             // APC 7845
	Port{7845, 17}:   "apc-7845",             // APC 7845
	Port{7846, 6}:    "apc-7846",             // APC 7846
	Port{7846, 17}:   "apc-7846",             // APC 7846
	Port{7847, 6}:    "csoauth",              // A product key authentication protocol made by CSO
	Port{7869, 6}:    "mobileanalyzer",       // MobileAnalyzer& MobileMonitor
	Port{7870, 6}:    "rbt-smc",              // Riverbed Steelhead Mobile Service
	Port{7871, 6}:    "mdm",                  // Mobile Device Management
	Port{7872, 6}:    "mipv6tls",             // TLS-based Mobile IPv6 Security
	Port{7878, 6}:    "owms",                 // Opswise Message Service
	Port{7880, 6}:    "pss",                  // Pearson
	Port{7880, 17}:   "pss",                  // Pearson
	Port{7887, 6}:    "ubroker",              // Universal Broker
	Port{7887, 17}:   "ubroker",              // Universal Broker
	Port{7900, 6}:    "mevent",               // Multicast Event
	Port{7900, 17}:   "mevent",               // Multicast Event
	Port{7901, 6}:    "tnos-sp",              // TNOS Service Protocol
	Port{7901, 17}:   "tnos-sp",              // TNOS Service Protocol
	Port{7902, 6}:    "tnos-dp",              // TNOS shell Protocol
	Port{7902, 17}:   "tnos-dp",              // TNOS shell Protocol
	Port{7903, 6}:    "tnos-dps",             // TNOS Secure DiaguardProtocol
	Port{7903, 17}:   "tnos-dps",             // TNOS Secure DiaguardProtocol
	Port{7913, 6}:    "qo-secure",            // QuickObjects secure port
	Port{7913, 17}:   "qo-secure",            // QuickObjects secure port
	Port{7932, 6}:    "t2-drm",               // Tier 2 Data Resource Manager
	Port{7932, 17}:   "t2-drm",               // Tier 2 Data Resource Manager
	Port{7933, 6}:    "t2-brm",               // Tier 2 Business Rules Manager
	Port{7933, 17}:   "t2-brm",               // Tier 2 Business Rules Manager
	Port{7937, 6}:    "nsrexecd",             // Legato NetWorker
	Port{7938, 6}:    "lgtomapper",           // Legato portmapper
	Port{7962, 6}:    "generalsync",          // Encrypted, extendable, general-purpose synchronization protocol
	Port{7967, 6}:    "supercell",            // Missing description for supercell
	Port{7967, 17}:   "supercell",            // Supercell
	Port{7979, 6}:    "micromuse-ncps",       // Missing description for micromuse-ncps
	Port{7979, 17}:   "micromuse-ncps",       // Micromuse-ncps
	Port{7980, 6}:    "quest-vista",          // Quest Vista
	Port{7980, 17}:   "quest-vista",          // Quest Vista
	Port{7981, 6}:    "sossd-collect",        // Spotlight on SQL Server Desktop Collect
	Port{7982, 6}:    "sossd-agent",          // sossd-disc | Spotlight on SQL Server Desktop Agent | Spotlight on SQL Server Desktop Agent Discovery
	Port{7982, 17}:   "sossd-disc",           // Spotlight on SQL Server Desktop Agent Discovery
	Port{7997, 6}:    "pushns",               // PUSH Notification Service
	Port{7998, 6}:    "usicontentpush",       // USI Content Push Service
	Port{7998, 17}:   "usicontentpush",       // USI Content Push Service
	Port{7999, 6}:    "irdmi2",               // Missing description for irdmi2
	Port{7999, 17}:   "irdmi2",               // iRDMI2
	Port{8000, 6}:    "http-alt",             // irdmi | A common alternative http port | iRDMI
	Port{8000, 17}:   "irdmi",                // iRDMI
	Port{8001, 6}:    "vcom-tunnel",          // VCOM Tunnel
	Port{8001, 17}:   "vcom-tunnel",          // VCOM Tunnel
	Port{8002, 6}:    "teradataordbms",       // Teradata ORDBMS
	Port{8002, 17}:   "teradataordbms",       // Teradata ORDBMS
	Port{8003, 6}:    "mcreport",             // Mulberry Connect Reporting Service
	Port{8003, 17}:   "mcreport",             // Mulberry Connect Reporting Service
	Port{8004, 6}:    "p2pevolvenet",         // Opensource Evolv Enterprise Platform P2P Network Node Connection Protocol
	Port{8005, 6}:    "mxi",                  // MXI Generation II for z OS
	Port{8005, 17}:   "mxi",                  // MXI Generation II for z OS
	Port{8006, 6}:    "wpl-analytics",        // wpl-disc | World Programming analytics | World Programming analytics discovery
	Port{8007, 6}:    "ajp12",                // warppipe | Apache JServ Protocol 1.x | I O oriented cluster computing software
	Port{8008, 6}:    "http",                 // http-alt | IBM HTTP server | HTTP Alternate
	Port{8008, 17}:   "http-alt",             // HTTP Alternate
	Port{8009, 6}:    "ajp13",                // Apache JServ Protocol 1.3
	Port{8010, 6}:    "xmpp",                 // XMPP File Transfer
	Port{8019, 6}:    "qbdb",                 // QB DB Dynamic Port
	Port{8019, 17}:   "qbdb",                 // QB DB Dynamic Port
	Port{8020, 6}:    "intu-ec-svcdisc",      // Intuit Entitlement Service and Discovery
	Port{8020, 17}:   "intu-ec-svcdisc",      // Intuit Entitlement Service and Discovery
	Port{8021, 6}:    "ftp-proxy",            // intu-ec-client | Common FTP proxy port | Intuit Entitlement Client
	Port{8021, 17}:   "intu-ec-client",       // Intuit Entitlement Client
	Port{8022, 6}:    "oa-system",            // Missing description for oa-system
	Port{8022, 17}:   "oa-system",            // Missing description for oa-system
	Port{8025, 6}:    "ca-audit-da",          // CA Audit Distribution Agent
	Port{8025, 17}:   "ca-audit-da",          // CA Audit Distribution Agent
	Port{8026, 6}:    "ca-audit-ds",          // CA Audit Distribution Server
	Port{8026, 17}:   "ca-audit-ds",          // CA Audit Distribution Server
	Port{8032, 6}:    "pro-ed",               // ProEd
	Port{8032, 17}:   "pro-ed",               // ProEd
	Port{8033, 6}:    "mindprint",            // Missing description for mindprint
	Port{8033, 17}:   "mindprint",            // MindPrint
	Port{8034, 6}:    "vantronix-mgmt",       // .vantronix Management
	Port{8034, 17}:   "vantronix-mgmt",       // .vantronix Management
	Port{8040, 6}:    "ampify",               // Ampify Messaging Protocol
	Port{8040, 17}:   "ampify",               // Ampify Messaging Protocol
	Port{8041, 6}:    "enguity-xccetp",       // Xcorpeon ASIC Carrier Ethernet Transport
	Port{8042, 6}:    "fs-agent",             // FireScope Agent
	Port{8043, 6}:    "fs-server",            // FireScope Server
	Port{8044, 6}:    "fs-mgmt",              // FireScope Management Interface
	Port{8051, 6}:    "rocrail",              // Rocrail Client Service
	Port{8052, 6}:    "senomix01",            // Senomix Timesheets Server
	Port{8052, 17}:   "senomix01",            // Senomix Timesheets Server
	Port{8053, 6}:    "senomix02",            // Senomix Timesheets Client [1 year assignment]
	Port{8053, 17}:   "senomix02",            // Senomix Timesheets Client [1 year assignment]
	Port{8054, 6}:    "senomix03",            // Senomix Timesheets Server [1 year assignment]
	Port{8054, 17}:   "senomix03",            // Senomix Timesheets Server [1 year assignment]
	Port{8055, 6}:    "senomix04",            // Senomix Timesheets Server [1 year assignment]
	Port{8055, 17}:   "senomix04",            // Senomix Timesheets Server [1 year assignment]
	Port{8056, 6}:    "senomix05",            // Senomix Timesheets Server [1 year assignment]
	Port{8056, 17}:   "senomix05",            // Senomix Timesheets Server [1 year assignment]
	Port{8057, 6}:    "senomix06",            // Senomix Timesheets Client [1 year assignment]
	Port{8057, 17}:   "senomix06",            // Senomix Timesheets Client [1 year assignment]
	Port{8058, 6}:    "senomix07",            // Senomix Timesheets Client [1 year assignment]
	Port{8058, 17}:   "senomix07",            // Senomix Timesheets Client [1 year assignment]
	Port{8059, 6}:    "senomix08",            // Senomix Timesheets Client [1 year assignment]
	Port{8059, 17}:   "senomix08",            // Senomix Timesheets Client [1 year assignment]
	Port{8060, 6}:    "aero",                 // Asymmetric Extended Route Optimization (AERO)
	Port{8066, 6}:    "toad-bi-appsrvr",      // Toad BI Application Server
	Port{8067, 6}:    "infi-async",           // Infinidat async replication
	Port{8070, 6}:    "ucs-isc",              // Oracle Unified Communication Suite's Indexed Search Converter
	Port{8074, 6}:    "gadugadu",             // Gadu-Gadu
	Port{8074, 17}:   "gadugadu",             // Gadu-Gadu
	Port{8076, 6}:    "slnp",                 // SLNP (Simple Library Network Protocol) by Sisis Informationssysteme GmbH
	Port{8077, 6}:    "mles",                 // Mles is a client-server data distribution protocol targeted to serve as a lightweight and reliable distributed publish subscribe database service.
	Port{8080, 6}:    "http-proxy",           // http-alt | Common HTTP proxy second web server port | HTTP Alternate (see port 80)
	Port{8080, 17}:   "http-alt",             // HTTP Alternate (see port 80)
	Port{8081, 6}:    "blackice-icecap",      // sunproxyadmin | ICECap user console | Sun Proxy Admin Service
	Port{8081, 17}:   "sunproxyadmin",        // Sun Proxy Admin Service
	Port{8082, 6}:    "blackice-alerts",      // us-cli | BlackIce Alerts sent to this port | Utilistor (Client)
	Port{8082, 17}:   "us-cli",               // Utilistor (Client)
	Port{8083, 6}:    "us-srv",               // Utilistor (Server)
	Port{8083, 17}:   "us-srv",               // Utilistor (Server)
	Port{8086, 6}:    "d-s-n",                // Distributed SCADA Networking Rendezvous Port
	Port{8086, 17}:   "d-s-n",                // Distributed SCADA Networking Rendezvous Port
	Port{8087, 6}:    "simplifymedia",        // Simplify Media SPP Protocol
	Port{8087, 17}:   "simplifymedia",        // Simplify Media SPP Protocol
	Port{8088, 6}:    "radan-http",           // Radan HTTP
	Port{8088, 17}:   "radan-http",           // Radan HTTP
	Port{8090, 6}:    "opsmessaging",         // Vehicle to station messaging
	Port{8091, 6}:    "jamlink",              // Jam Link Framework
	Port{8097, 6}:    "sac",                  // SAC Port Id
	Port{8097, 17}:   "sac",                  // SAC Port Id
	Port{8100, 6}:    "xprint-server",        // Xprint Server
	Port{8100, 17}:   "xprint-server",        // Xprint Server
	Port{8101, 6}:    "ldoms-migr",           // Logical Domains Migration
	Port{8102, 6}:    "kz-migr",              // Oracle Kernel zones migration server
	Port{8115, 6}:    "mtl8000-matrix",       // MTL8000 Matrix
	Port{8115, 17}:   "mtl8000-matrix",       // MTL8000 Matrix
	Port{8116, 6}:    "cp-cluster",           // Check Point Clustering
	Port{8116, 17}:   "cp-cluster",           // Check Point Clustering
	Port{8117, 6}:    "purityrpc",            // Purity replication clustering and remote management
	Port{8118, 6}:    "privoxy",              // Privoxy, www.privoxy.org | Privoxy HTTP proxy
	Port{8118, 17}:   "privoxy",              // Privoxy HTTP proxy
	Port{8121, 6}:    "apollo-data",          // Apollo Data Port
	Port{8121, 17}:   "apollo-data",          // Apollo Data Port
	Port{8122, 6}:    "apollo-admin",         // Apollo Admin Port
	Port{8122, 17}:   "apollo-admin",         // Apollo Admin Port
	Port{8123, 6}:    "polipo",               // Polipo open source web proxy cache
	Port{8128, 6}:    "paycash-online",       // PayCash Online Protocol
	Port{8128, 17}:   "paycash-online",       // PayCash Online Protocol
	Port{8129, 6}:    "paycash-wbp",          // PayCash Wallet-Browser
	Port{8129, 17}:   "paycash-wbp",          // PayCash Wallet-Browser
	Port{8130, 6}:    "indigo-vrmi",          // Missing description for indigo-vrmi
	Port{8130, 17}:   "indigo-vrmi",          // INDIGO-VRMI
	Port{8131, 6}:    "indigo-vbcp",          // Missing description for indigo-vbcp
	Port{8131, 17}:   "indigo-vbcp",          // INDIGO-VBCP
	Port{8132, 6}:    "dbabble",              // Missing description for dbabble
	Port{8132, 17}:   "dbabble",              // Missing description for dbabble
	Port{8140, 6}:    "puppet",               // The Puppet master service
	Port{8148, 6}:    "isdd",                 // i-SDD file transfer
	Port{8148, 17}:   "isdd",                 // i-SDD file transfer
	Port{8149, 6}:    "eor-game",             // Edge of Reality game data
	Port{8153, 6}:    "quantastor",           // QuantaStor Management Interface
	Port{8160, 6}:    "patrol",               // Missing description for patrol
	Port{8160, 17}:   "patrol",               // Patrol
	Port{8161, 6}:    "patrol-snmp",          // Patrol SNMP
	Port{8161, 17}:   "patrol-snmp",          // Patrol SNMP
	Port{8162, 6}:    "lpar2rrd",             // LPAR2RRD client server communication
	Port{8181, 6}:    "intermapper",          // Intermapper network management system
	Port{8182, 6}:    "vmware-fdm",           // VMware Fault Domain Manager
	Port{8182, 17}:   "vmware-fdm",           // VMware Fault Domain Manager
	Port{8183, 6}:    "proremote",            // Missing description for proremote
	Port{8184, 6}:    "itach",                // Remote iTach Connection
	Port{8184, 17}:   "itach",                // Remote iTach Connection
	Port{8190, 6}:    "gcp-rphy",             // Generic control plane for RPHY
	Port{8191, 6}:    "limnerpressure",       // Limner Pressure
	Port{8192, 6}:    "sophos",               // spytechphone | Sophos Remote Management System | SpyTech Phone Service
	Port{8192, 17}:   "sophos",               // Sophos Remote Management System
	Port{8193, 6}:    "sophos",               // Sophos Remote Management System
	Port{8193, 17}:   "sophos",               // Sophos Remote Management System
	Port{8194, 6}:    "sophos",               // blp1 | Sophos Remote Management System | Bloomberg data API
	Port{8194, 17}:   "sophos",               // Sophos Remote Management System
	Port{8195, 6}:    "blp2",                 // Bloomberg feed
	Port{8195, 17}:   "blp2",                 // Bloomberg feed
	Port{8199, 6}:    "vvr-data",             // VVR DATA
	Port{8199, 17}:   "vvr-data",             // VVR DATA
	Port{8200, 6}:    "trivnet1",             // TRIVNET
	Port{8200, 17}:   "trivnet1",             // TRIVNET
	Port{8201, 6}:    "trivnet2",             // TRIVNET
	Port{8201, 17}:   "trivnet2",             // TRIVNET
	Port{8202, 6}:    "aesop",                // Audio+Ethernet Standard Open Protocol
	Port{8204, 6}:    "lm-perfworks",         // LM Perfworks
	Port{8204, 17}:   "lm-perfworks",         // LM Perfworks
	Port{8205, 6}:    "lm-instmgr",           // LM Instmgr
	Port{8205, 17}:   "lm-instmgr",           // LM Instmgr
	Port{8206, 6}:    "lm-dta",               // LM Dta
	Port{8206, 17}:   "lm-dta",               // LM Dta
	Port{8207, 6}:    "lm-sserver",           // LM SServer
	Port{8207, 17}:   "lm-sserver",           // LM SServer
	Port{8208, 6}:    "lm-webwatcher",        // LM Webwatcher
	Port{8208, 17}:   "lm-webwatcher",        // LM Webwatcher
	Port{8230, 6}:    "rexecj",               // RexecJ Server
	Port{8230, 17}:   "rexecj",               // RexecJ Server
	Port{8231, 6}:    "hncp-udp-port",        // HNCP
	Port{8232, 6}:    "hncp-dtls-port",       // HNCP over DTLS
	Port{8243, 6}:    "synapse-nhttps",       // Synapse Non Blocking HTTPS
	Port{8243, 17}:   "synapse-nhttps",       // Synapse Non Blocking HTTPS
	Port{8270, 6}:    "robot-remote",         // Robot Framework Remote Library Interface
	Port{8276, 6}:    "pando-sec",            // Pando Media Controlled Distribution
	Port{8276, 17}:   "pando-sec",            // Pando Media Controlled Distribution
	Port{8280, 6}:    "synapse-nhttp",        // Synapse Non Blocking HTTP
	Port{8280, 17}:   "synapse-nhttp",        // Synapse Non Blocking HTTP
	Port{8282, 6}:    "libelle",              // libelle-disc | Libelle EnterpriseBus | Libelle EnterpriseBus discovery
	Port{8292, 6}:    "blp3",                 // Bloomberg professional
	Port{8292, 17}:   "blp3",                 // Bloomberg professional
	Port{8293, 6}:    "hiperscan-id",         // Hiperscan Identification Service
	Port{8294, 6}:    "blp4",                 // Bloomberg intelligent client
	Port{8294, 17}:   "blp4",                 // Bloomberg intelligent client
	Port{8300, 6}:    "tmi",                  // Transport Management Interface
	Port{8300, 17}:   "tmi",                  // Transport Management Interface
	Port{8301, 6}:    "amberon",              // Amberon PPC PPS
	Port{8301, 17}:   "amberon",              // Amberon PPC PPS
	Port{8313, 6}:    "hub-open-net",         // Hub Open Network
	Port{8320, 6}:    "tnp-discover",         // Thin(ium) Network Protocol
	Port{8320, 17}:   "tnp-discover",         // Thin(ium) Network Protocol
	Port{8321, 6}:    "tnp",                  // Thin(ium) Network Protocol
	Port{8321, 17}:   "tnp",                  // Thin(ium) Network Protocol
	Port{8322, 6}:    "garmin-marine",        // Garmin Marine
	Port{8333, 6}:    "bitcoin",              // Bitcoin crypto currency - https:  en.bitcoin.it wiki Running_Bitcoin
	Port{8351, 6}:    "server-find",          // Server Find
	Port{8351, 17}:   "server-find",          // Server Find
	Port{8376, 6}:    "cruise-enum",          // Cruise ENUM
	Port{8376, 17}:   "cruise-enum",          // Cruise ENUM
	Port{8377, 6}:    "cruise-swroute",       // Cruise SWROUTE
	Port{8377, 17}:   "cruise-swroute",       // Cruise SWROUTE
	Port{8378, 6}:    "cruise-config",        // Cruise CONFIG
	Port{8378, 17}:   "cruise-config",        // Cruise CONFIG
	Port{8379, 6}:    "cruise-diags",         // Cruise DIAGS
	Port{8379, 17}:   "cruise-diags",         // Cruise DIAGS
	Port{8380, 6}:    "cruise-update",        // Cruise UPDATE
	Port{8380, 17}:   "cruise-update",        // Cruise UPDATE
	Port{8383, 6}:    "m2mservices",          // M2m Services
	Port{8383, 17}:   "m2mservices",          // M2m Services
	Port{8384, 6}:    "marathontp",           // Marathon Transport Protocol
	Port{8400, 6}:    "cvd",                  // Missing description for cvd
	Port{8400, 17}:   "cvd",                  // Missing description for cvd
	Port{8401, 6}:    "sabarsd",              // Missing description for sabarsd
	Port{8401, 17}:   "sabarsd",              // Missing description for sabarsd
	Port{8402, 6}:    "abarsd",               // Missing description for abarsd
	Port{8402, 17}:   "abarsd",               // Missing description for abarsd
	Port{8403, 6}:    "admind",               // Missing description for admind
	Port{8403, 17}:   "admind",               // Missing description for admind
	Port{8404, 6}:    "svcloud",              // SuperVault Cloud
	Port{8405, 6}:    "svbackup",             // SuperVault Backup
	Port{8415, 6}:    "dlpx-sp",              // Delphix Session Protocol
	Port{8416, 6}:    "espeech",              // eSpeech Session Protocol
	Port{8416, 17}:   "espeech",              // eSpeech Session Protocol
	Port{8417, 6}:    "espeech-rtp",          // eSpeech RTP Protocol
	Port{8417, 17}:   "espeech-rtp",          // eSpeech RTP Protocol
	Port{8423, 6}:    "aritts",               // Aristech text-to-speech server
	Port{8442, 6}:    "cybro-a-bus",          // CyBro A-bus Protocol
	Port{8442, 17}:   "cybro-a-bus",          // CyBro A-bus Protocol
	Port{8443, 6}:    "https-alt",            // pcsync-https | Common alternative https port | PCsync HTTPS
	Port{8443, 17}:   "pcsync-https",         // PCsync HTTPS
	Port{8444, 6}:    "pcsync-http",          // PCsync HTTP
	Port{8444, 17}:   "pcsync-http",          // PCsync HTTP
	Port{8445, 6}:    "copy",                 // copy-disc | Port for copy peer sync feature | Port for copy discovery
	Port{8450, 6}:    "npmp",                 // Missing description for npmp
	Port{8450, 17}:   "npmp",                 // Missing description for npmp
	Port{8457, 6}:    "nexentamv",            // Nexenta Management GUI
	Port{8470, 6}:    "cisco-avp",            // Cisco Address Validation Protocol
	Port{8471, 132}:  "pim-port",             // PIM over Reliable Transport
	Port{8471, 6}:    "pim-port",             // PIM over Reliable Transport
	Port{8471, 17}:   "pim-port",             // PIM over Reliable Transport
	Port{8472, 6}:    "otv",                  // Overlay Transport Virtualization (OTV)
	Port{8472, 17}:   "otv",                  // Overlay Transport Virtualization (OTV)
	Port{8473, 6}:    "vp2p",                 // Virtual Point to Point
	Port{8473, 17}:   "vp2p",                 // Virtual Point to Point
	Port{8474, 6}:    "noteshare",            // AquaMinds NoteShare
	Port{8474, 17}:   "noteshare",            // AquaMinds NoteShare
	Port{8500, 6}:    "fmtp",                 // Flight Message Transfer Protocol
	Port{8500, 17}:   "fmtp",                 // Flight Message Transfer Protocol
	Port{8501, 6}:    "cmtp-mgt",             // cmtp-av | CYTEL Message Transfer Management | CYTEL Message Transfer Audio and Video
	Port{8502, 6}:    "ftnmtp",               // FTN Message Transfer Protocol
	Port{8503, 6}:    "lsp-self-ping",        // MPLS LSP Self-Ping
	Port{8554, 6}:    "rtsp-alt",             // RTSP Alternate (see port 554)
	Port{8554, 17}:   "rtsp-alt",             // RTSP Alternate (see port 554)
	Port{8555, 6}:    "d-fence",              // SYMAX D-FENCE
	Port{8555, 17}:   "d-fence",              // SYMAX D-FENCE
	Port{8567, 6}:    "oap-admin",            // dof-tunnel | Object Access Protocol Administration | DOF Tunneling Protocol
	Port{8567, 17}:   "oap-admin",            // Object Access Protocol Administration
	Port{8600, 6}:    "asterix",              // Surveillance Data
	Port{8600, 17}:   "asterix",              // Surveillance Data
	Port{8609, 6}:    "canon-cpp-disc",       // Canon Compact Printer Protocol Discovery
	Port{8610, 6}:    "canon-mfnp",           // Canon MFNP Service
	Port{8610, 17}:   "canon-mfnp",           // Canon MFNP Service
	Port{8611, 6}:    "canon-bjnp1",          // Canon BJNP Port 1
	Port{8611, 17}:   "canon-bjnp1",          // Canon BJNP Port 1
	Port{8612, 6}:    "canon-bjnp2",          // Canon BJNP Port 2
	Port{8612, 17}:   "canon-bjnp2",          // Canon BJNP Port 2
	Port{8613, 6}:    "canon-bjnp3",          // Canon BJNP Port 3
	Port{8613, 17}:   "canon-bjnp3",          // Canon BJNP Port 3
	Port{8614, 6}:    "canon-bjnp4",          // Canon BJNP Port 4
	Port{8614, 17}:   "canon-bjnp4",          // Canon BJNP Port 4
	Port{8615, 6}:    "imink",                // Imink Service Control
	Port{8665, 6}:    "monetra",              // Missing description for monetra
	Port{8666, 6}:    "monetra-admin",        // Monetra Administrative Access
	Port{8675, 6}:    "msi-cps-rm",           // msi-cps-rm-disc | Motorola Solutions Customer Programming Software for Radio Management | Motorola Solutions Customer Programming Software for Radio Management Discovery
	Port{8686, 6}:    "sun-as-jmxrmi",        // Sun App Server - JMX RMI
	Port{8686, 17}:   "sun-as-jmxrmi",        // Sun App Server - JMX RMI
	Port{8688, 6}:    "openremote-ctrl",      // OpenRemote Controller HTTP REST
	Port{8699, 6}:    "vnyx",                 // VNYX Primary Port
	Port{8699, 17}:   "vnyx",                 // VNYX Primary Port
	Port{8711, 6}:    "nvc",                  // Nuance Voice Control
	Port{8732, 6}:    "dtp-net",              // DASGIP Net Services
	Port{8732, 17}:   "dtp-net",              // DASGIP Net Services
	Port{8733, 6}:    "ibus",                 // Missing description for ibus
	Port{8733, 17}:   "ibus",                 // iBus
	Port{8750, 6}:    "dey-keyneg",           // DEY Storage Key Negotiation
	Port{8763, 6}:    "mc-appserver",         // Missing description for mc-appserver
	Port{8763, 17}:   "mc-appserver",         // MC-APPSERVER
	Port{8764, 6}:    "openqueue",            // Missing description for openqueue
	Port{8764, 17}:   "openqueue",            // OPENQUEUE
	Port{8765, 6}:    "ultraseek-http",       // Ultraseek HTTP
	Port{8765, 17}:   "ultraseek-http",       // Ultraseek HTTP
	Port{8766, 6}:    "amcs",                 // Agilent Connectivity Service
	Port{8770, 6}:    "apple-iphoto",         // dpap | Apple iPhoto sharing | Digital Photo Access Protocol (iPhoto)
	Port{8770, 17}:   "dpap",                 // Digital Photo Access Protocol
	Port{8778, 6}:    "uec",                  // Stonebranch Universal Enterprise Controller
	Port{8786, 6}:    "msgclnt",              // Message Client
	Port{8786, 17}:   "msgclnt",              // Message Client
	Port{8787, 6}:    "msgsrvr",              // Message Server
	Port{8787, 17}:   "msgsrvr",              // Message Server
	Port{8793, 6}:    "acd-pm",               // Accedian Performance Measurement
	Port{8800, 6}:    "sunwebadmin",          // Sun Web Server Admin Service
	Port{8800, 17}:   "sunwebadmin",          // Sun Web Server Admin Service
	Port{8804, 6}:    "truecm",               // Missing description for truecm
	Port{8804, 17}:   "truecm",               // Missing description for truecm
	Port{8805, 6}:    "pfcp",                 // Destination Port number for PFCP
	Port{8808, 6}:    "ssports-bcast",        // STATSports Broadcast Service
	Port{8834, 6}:    "nessus-xmlrpc",        // Missing description for nessus-xmlrpc
	Port{8873, 6}:    "dxspider",             // dxspider linking protocol
	Port{8873, 17}:   "dxspider",             // dxspider linking protocol
	Port{8880, 6}:    "cddbp-alt",            // CDDBP
	Port{8880, 17}:   "cddbp-alt",            // CDDBP
	Port{8881, 6}:    "galaxy4d",             // Galaxy4D Online Game Engine
	Port{8883, 6}:    "secure-mqtt",          // Secure MQTT
	Port{8883, 17}:   "secure-mqtt",          // Secure MQTT
	Port{8888, 6}:    "sun-answerbook",       // ddi-udp-1 | ddi-tcp-1 | Sun Answerbook HTTP server.  Or gnump3d streaming music server | NewsEDGE server TCP (TCP 1) | NewsEDGE server UDP (UDP 1)
	Port{8888, 17}:   "ddi-udp-1",            // NewsEDGE server UDP (UDP 1)
	Port{8889, 6}:    "ddi-tcp-2",            // ddi-udp-2 | Desktop Data TCP 1 | NewsEDGE server broadcast
	Port{8889, 17}:   "ddi-udp-2",            // NewsEDGE server broadcast
	Port{8890, 6}:    "ddi-tcp-3",            // ddi-udp-3 | Desktop Data TCP 2 | NewsEDGE client broadcast
	Port{8890, 17}:   "ddi-udp-3",            // NewsEDGE client broadcast
	Port{8891, 6}:    "ddi-tcp-4",            // ddi-udp-4 | Desktop Data TCP 3: NESS application | Desktop Data UDP 3: NESS application
	Port{8891, 17}:   "ddi-udp-4",            // Desktop Data UDP 3: NESS application
	Port{8892, 6}:    "seosload",             // ddi-udp-5 | ddi-tcp-5 | From the new Computer Associates eTrust ACX | Desktop Data TCP 4: FARM product | Desktop Data UDP 4: FARM product
	Port{8892, 17}:   "ddi-udp-5",            // Desktop Data UDP 4: FARM product
	Port{8893, 6}:    "ddi-tcp-6",            // ddi-udp-6 | Desktop Data TCP 5: NewsEDGE Web application | Desktop Data UDP 5: NewsEDGE Web application
	Port{8893, 17}:   "ddi-udp-6",            // Desktop Data UDP 5: NewsEDGE Web application
	Port{8894, 6}:    "ddi-tcp-7",            // ddi-udp-7 | Desktop Data TCP 6: COAL application | Desktop Data UDP 6: COAL application
	Port{8894, 17}:   "ddi-udp-7",            // Desktop Data UDP 6: COAL application
	Port{8899, 6}:    "ospf-lite",            // Missing description for ospf-lite
	Port{8899, 17}:   "ospf-lite",            // Missing description for ospf-lite
	Port{8900, 6}:    "jmb-cds1",             // JMB-CDS 1
	Port{8900, 17}:   "jmb-cds1",             // JMB-CDS 1
	Port{8901, 6}:    "jmb-cds2",             // JMB-CDS 2
	Port{8901, 17}:   "jmb-cds2",             // JMB-CDS 2
	Port{8910, 6}:    "manyone-http",         // Missing description for manyone-http
	Port{8910, 17}:   "manyone-http",         // Missing description for manyone-http
	Port{8911, 6}:    "manyone-xml",          // Missing description for manyone-xml
	Port{8911, 17}:   "manyone-xml",          // Missing description for manyone-xml
	Port{8912, 6}:    "wcbackup",             // Windows Client Backup
	Port{8912, 17}:   "wcbackup",             // Windows Client Backup
	Port{8913, 6}:    "dragonfly",            // Dragonfly System Service
	Port{8913, 17}:   "dragonfly",            // Dragonfly System Service
	Port{8937, 6}:    "twds",                 // Transaction Warehouse Data Service
	Port{8953, 6}:    "ub-dns-control",       // unbound dns nameserver control
	Port{8954, 6}:    "cumulus-admin",        // Cumulus Admin Port
	Port{8954, 17}:   "cumulus-admin",        // Cumulus Admin Port
	Port{8980, 6}:    "nod-provider",         // Network of Devices Provider
	Port{8981, 6}:    "nod-client",           // Network of Devices Client
	Port{8989, 6}:    "sunwebadmins",         // Sun Web Server SSL Admin Service
	Port{8989, 17}:   "sunwebadmins",         // Sun Web Server SSL Admin Service
	Port{8990, 6}:    "http-wmap",            // webmail HTTP service
	Port{8990, 17}:   "http-wmap",            // webmail HTTP service
	Port{8991, 6}:    "https-wmap",           // webmail HTTPS service
	Port{8991, 17}:   "https-wmap",           // webmail HTTPS service
	Port{8997, 6}:    "oracle-ms-ens",        // Oracle Messaging Server Event Notification Service
	Port{8998, 6}:    "canto-roboflow",       // Canto RoboFlow Control
	Port{8999, 6}:    "bctp",                 // Brodos Crypto Trade Protocol
	Port{8999, 17}:   "bctp",                 // Brodos Crypto Trade Protocol
	Port{9000, 6}:    "cslistener",           // Missing description for cslistener
	Port{9000, 17}:   "cslistener",           // CSlistener
	Port{9001, 6}:    "tor-orport",           // etlservicemgr | Tor ORPort | ETL Service Manager
	Port{9001, 17}:   "etlservicemgr",        // ETL Service Manager
	Port{9002, 6}:    "dynamid",              // DynamID authentication
	Port{9002, 17}:   "dynamid",              // DynamID authentication
	Port{9005, 6}:    "golem",                // Golem Inter-System RPC
	Port{9007, 6}:    "ogs-client",           // Open Grid Services Client
	Port{9007, 17}:   "ogs-client",           // Open Grid Services Client
	Port{9008, 6}:    "ogs-server",           // Open Grid Services Server
	Port{9009, 6}:    "pichat",               // Pichat Server
	Port{9009, 17}:   "pichat",               // Pichat Server
	Port{9010, 6}:    "sdr",                  // Secure Data Replicator Protocol
	Port{9020, 6}:    "tambora",              // Missing description for tambora
	Port{9020, 17}:   "tambora",              // TAMBORA
	Port{9021, 6}:    "panagolin-ident",      // Pangolin Identification
	Port{9021, 17}:   "panagolin-ident",      // Pangolin Identification
	Port{9022, 6}:    "paragent",             // PrivateArk Remote Agent
	Port{9022, 17}:   "paragent",             // PrivateArk Remote Agent
	Port{9023, 6}:    "swa-1",                // Secure Web Access - 1
	Port{9023, 17}:   "swa-1",                // Secure Web Access - 1
	Port{9024, 6}:    "swa-2",                // Secure Web Access - 2
	Port{9024, 17}:   "swa-2",                // Secure Web Access - 2
	Port{9025, 6}:    "swa-3",                // Secure Web Access - 3
	Port{9025, 17}:   "swa-3",                // Secure Web Access - 3
	Port{9026, 6}:    "swa-4",                // Secure Web Access - 4
	Port{9026, 17}:   "swa-4",                // Secure Web Access - 4
	Port{9040, 6}:    "tor-trans",            // Tor TransPort, www.torproject.org
	Port{9050, 6}:    "tor-socks",            // versiera | Tor SocksPort, www.torproject.org | Versiera Agent Listener
	Port{9051, 6}:    "tor-control",          // fio-cmgmt | Tor ControlPort, www.torproject.org | Fusion-io Central Manager Service
	Port{9060, 6}:    "CardWeb-IO",           // CardWeb-RT | CardWeb request-response I O exchange | CardWeb realtime device data
	Port{9080, 6}:    "glrpc",                // Groove GLRPC
	Port{9080, 17}:   "glrpc",                // Groove GLRPC
	Port{9081, 6}:    "cisco-aqos",           // Required for Adaptive Quality of Service
	Port{9082, 132}:  "lcs-ap",               // LCS Application Protocol
	Port{9083, 6}:    "emc-pp-mgmtsvc",       // EMC PowerPath Mgmt Service
	Port{9084, 132}:  "aurora",               // IBM AURORA Performance Visualizer
	Port{9084, 6}:    "aurora",               // IBM AURORA Performance Visualizer
	Port{9084, 17}:   "aurora",               // IBM AURORA Performance Visualizer
	Port{9085, 6}:    "ibm-rsyscon",          // IBM Remote System Console
	Port{9085, 17}:   "ibm-rsyscon",          // IBM Remote System Console
	Port{9086, 6}:    "net2display",          // Vesa Net2Display
	Port{9086, 17}:   "net2display",          // Vesa Net2Display
	Port{9087, 6}:    "classic",              // Classic Data Server
	Port{9087, 17}:   "classic",              // Classic Data Server
	Port{9088, 6}:    "sqlexec",              // IBM Informix SQL Interface
	Port{9088, 17}:   "sqlexec",              // IBM Informix SQL Interface
	Port{9089, 6}:    "sqlexec-ssl",          // IBM Informix SQL Interface - Encrypted
	Port{9089, 17}:   "sqlexec-ssl",          // IBM Informix SQL Interface - Encrypted
	Port{9090, 6}:    "zeus-admin",           // websm | Zeus admin server | WebSM
	Port{9090, 17}:   "websm",                // WebSM
	Port{9091, 6}:    "xmltec-xmlmail",       // Missing description for xmltec-xmlmail
	Port{9091, 17}:   "xmltec-xmlmail",       // Missing description for xmltec-xmlmail
	Port{9092, 6}:    "XmlIpcRegSvc",         // Xml-Ipc Server Reg
	Port{9092, 17}:   "XmlIpcRegSvc",         // Xml-Ipc Server Reg
	Port{9093, 6}:    "copycat",              // Copycat database replication service
	Port{9100, 6}:    "jetdirect",            // pdl-datastream | hp-pdl-datastr | HP JetDirect card | PDL Data Streaming Port | Printer PDL Data Stream
	Port{9100, 17}:   "hp-pdl-datastr",       // PDL Data Streaming Port
	Port{9101, 6}:    "jetdirect",            // bacula-dir | HP JetDirect card | Bacula Director
	Port{9101, 17}:   "bacula-dir",           // Bacula Director
	Port{9102, 6}:    "jetdirect",            // bacula-fd | HP JetDirect card. Also used (and officially registered for) Bacula File Daemon (an open source backup system) | Bacula File Daemon
	Port{9102, 17}:   "bacula-fd",            // Bacula File Daemon
	Port{9103, 6}:    "jetdirect",            // bacula-sd | HP JetDirect card | Bacula Storage Daemon
	Port{9103, 17}:   "bacula-sd",            // Bacula Storage Daemon
	Port{9104, 6}:    "jetdirect",            // peerwire | HP JetDirect card | PeerWire
	Port{9104, 17}:   "peerwire",             // PeerWire
	Port{9105, 6}:    "jetdirect",            // xadmin | HP JetDirect card | Xadmin Control Service
	Port{9105, 17}:   "xadmin",               // Xadmin Control Service
	Port{9106, 6}:    "jetdirect",            // astergate-disc | astergate | HP JetDirect card | Astergate Control Service | Astergate Discovery Service
	Port{9106, 17}:   "astergate-disc",       // Astergate Discovery Service
	Port{9107, 6}:    "jetdirect",            // astergatefax | HP JetDirect card | AstergateFax Control Service
	Port{9111, 6}:    "DragonIDSConsole",     // Dragon IDS Console
	Port{9119, 6}:    "mxit",                 // MXit Instant Messaging
	Port{9119, 17}:   "mxit",                 // MXit Instant Messaging
	Port{9122, 6}:    "grcmp",                // Global Relay compliant mobile instant messaging protocol
	Port{9123, 6}:    "grcp",                 // Global Relay compliant instant messaging protocol
	Port{9131, 6}:    "dddp",                 // Dynamic Device Discovery
	Port{9131, 17}:   "dddp",                 // Dynamic Device Discovery
	Port{9152, 6}:    "ms-sql2000",           // Missing description for ms-sql2000
	Port{9160, 6}:    "apani1",               // Missing description for apani1
	Port{9160, 17}:   "apani1",               // Missing description for apani1
	Port{9161, 6}:    "apani2",               // Missing description for apani2
	Port{9161, 17}:   "apani2",               // Missing description for apani2
	Port{9162, 6}:    "apani3",               // Missing description for apani3
	Port{9162, 17}:   "apani3",               // Missing description for apani3
	Port{9163, 6}:    "apani4",               // Missing description for apani4
	Port{9163, 17}:   "apani4",               // Missing description for apani4
	Port{9164, 6}:    "apani5",               // Missing description for apani5
	Port{9164, 17}:   "apani5",               // Missing description for apani5
	Port{9191, 6}:    "sun-as-jpda",          // Sun AppSvr JPDA
	Port{9191, 17}:   "sun-as-jpda",          // Sun AppSvr JPDA
	Port{9200, 6}:    "wap-wsp",              // WAP connectionless session services | WAP connectionless session service
	Port{9200, 17}:   "wap-wsp",              // WAP connectionless session services
	Port{9201, 6}:    "wap-wsp-wtp",          // WAP session service
	Port{9201, 17}:   "wap-wsp-wtp",          // WAP session service
	Port{9202, 6}:    "wap-wsp-s",            // WAP secure connectionless session service
	Port{9202, 17}:   "wap-wsp-s",            // WAP secure connectionless session service
	Port{9203, 6}:    "wap-wsp-wtp-s",        // WAP secure session service
	Port{9203, 17}:   "wap-wsp-wtp-s",        // WAP secure session service
	Port{9204, 6}:    "wap-vcard",            // WAP vCard
	Port{9204, 17}:   "wap-vcard",            // WAP vCard
	Port{9205, 6}:    "wap-vcal",             // WAP vCal
	Port{9205, 17}:   "wap-vcal",             // WAP vCal
	Port{9206, 6}:    "wap-vcard-s",          // WAP vCard Secure
	Port{9206, 17}:   "wap-vcard-s",          // WAP vCard Secure
	Port{9207, 6}:    "wap-vcal-s",           // WAP vCal Secure
	Port{9207, 17}:   "wap-vcal-s",           // WAP vCal Secure
	Port{9208, 6}:    "rjcdb-vcards",         // rjcdb vCard
	Port{9208, 17}:   "rjcdb-vcards",         // rjcdb vCard
	Port{9209, 6}:    "almobile-system",      // ALMobile System Service
	Port{9209, 17}:   "almobile-system",      // ALMobile System Service
	Port{9210, 6}:    "oma-mlp",              // OMA Mobile Location Protocol
	Port{9210, 17}:   "oma-mlp",              // OMA Mobile Location Protocol
	Port{9211, 6}:    "oma-mlp-s",            // OMA Mobile Location Protocol Secure
	Port{9211, 17}:   "oma-mlp-s",            // OMA Mobile Location Protocol Secure
	Port{9212, 6}:    "serverviewdbms",       // Server View dbms access [January 2005] | Server View dbms access
	Port{9212, 17}:   "serverviewdbms",       // Server View dbms access [January 2005]
	Port{9213, 6}:    "serverstart",          // ServerStart RemoteControl [August 2005] | ServerStart RemoteControl
	Port{9213, 17}:   "serverstart",          // ServerStart RemoteControl [August 2005]
	Port{9214, 6}:    "ipdcesgbs",            // IPDC ESG BootstrapService
	Port{9214, 17}:   "ipdcesgbs",            // IPDC ESG BootstrapService
	Port{9215, 6}:    "insis",                // Integrated Setup and Install Service
	Port{9215, 17}:   "insis",                // Integrated Setup and Install Service
	Port{9216, 6}:    "acme",                 // Aionex Communication Management Engine
	Port{9216, 17}:   "acme",                 // Aionex Communication Management Engine
	Port{9217, 6}:    "fsc-port",             // FSC Communication Port
	Port{9217, 17}:   "fsc-port",             // FSC Communication Port
	Port{9222, 6}:    "teamcoherence",        // QSC Team Coherence
	Port{9222, 17}:   "teamcoherence",        // QSC Team Coherence
	Port{9255, 6}:    "mon",                  // Manager On Network
	Port{9255, 17}:   "mon",                  // Manager On Network
	Port{9277, 6}:    "traingpsdata",         // GPS Data transmitted from train to ground network
	Port{9278, 6}:    "pegasus",              // Pegasus GPS Platform
	Port{9278, 17}:   "pegasus",              // Pegasus GPS Platform
	Port{9279, 6}:    "pegasus-ctl",          // Pegaus GPS System Control Interface
	Port{9279, 17}:   "pegasus-ctl",          // Pegaus GPS System Control Interface
	Port{9280, 6}:    "pgps",                 // Predicted GPS
	Port{9280, 17}:   "pgps",                 // Predicted GPS
	Port{9281, 6}:    "swtp-port1",           // SofaWare transport port 1
	Port{9281, 17}:   "swtp-port1",           // SofaWare transport port 1
	Port{9282, 6}:    "swtp-port2",           // SofaWare transport port 2
	Port{9282, 17}:   "swtp-port2",           // SofaWare transport port 2
	Port{9283, 6}:    "callwaveiam",          // Missing description for callwaveiam
	Port{9283, 17}:   "callwaveiam",          // CallWaveIAM
	Port{9284, 6}:    "visd",                 // VERITAS Information Serve
	Port{9284, 17}:   "visd",                 // VERITAS Information Serve
	Port{9285, 6}:    "n2h2server",           // N2H2 Filter Service Port
	Port{9285, 17}:   "n2h2server",           // N2H2 Filter Service Port
	Port{9286, 6}:    "n2receive",            // n2 monitoring receiver
	Port{9287, 6}:    "cumulus",              // Missing description for cumulus
	Port{9287, 17}:   "cumulus",              // Cumulus
	Port{9292, 6}:    "armtechdaemon",        // ArmTech Daemon
	Port{9292, 17}:   "armtechdaemon",        // ArmTech Daemon
	Port{9293, 6}:    "storview",             // StorView Client
	Port{9293, 17}:   "storview",             // StorView Client
	Port{9294, 6}:    "armcenterhttp",        // ARMCenter http Service
	Port{9294, 17}:   "armcenterhttp",        // ARMCenter http Service
	Port{9295, 6}:    "armcenterhttps",       // ARMCenter https Service
	Port{9295, 17}:   "armcenterhttps",       // ARMCenter https Service
	Port{9300, 6}:    "vrace",                // Virtual Racing Service
	Port{9300, 17}:   "vrace",                // Virtual Racing Service
	Port{9306, 6}:    "sphinxql",             // Sphinx search server (MySQL listener)
	Port{9312, 6}:    "sphinxapi",            // Sphinx search server
	Port{9318, 6}:    "secure-ts",            // PKIX TimeStamp over TLS
	Port{9318, 17}:   "secure-ts",            // PKIX TimeStamp over TLS
	Port{9321, 6}:    "guibase",              // Missing description for guibase
	Port{9321, 17}:   "guibase",              // Missing description for guibase
	Port{9333, 6}:    "litecoin",             // Litecoin crypto currency - https:  litecoin.info Litecoin.conf
	Port{9343, 6}:    "mpidcmgr",             // Missing description for mpidcmgr
	Port{9343, 17}:   "mpidcmgr",             // MpIdcMgr
	Port{9344, 6}:    "mphlpdmc",             // Missing description for mphlpdmc
	Port{9344, 17}:   "mphlpdmc",             // Mphlpdmc
	Port{9345, 6}:    "rancher",              // Rancher Agent
	Port{9346, 6}:    "ctechlicensing",       // C Tech Licensing
	Port{9346, 17}:   "ctechlicensing",       // C Tech Licensing
	Port{9374, 6}:    "fjdmimgr",             // Missing description for fjdmimgr
	Port{9374, 17}:   "fjdmimgr",             // Missing description for fjdmimgr
	Port{9380, 6}:    "boxp",                 // Brivs! Open Extensible Protocol
	Port{9380, 17}:   "boxp",                 // Brivs! Open Extensible Protocol
	Port{9387, 6}:    "d2dconfig",            // D2D Configuration Service
	Port{9388, 6}:    "d2ddatatrans",         // D2D Data Transfer Service
	Port{9389, 6}:    "adws",                 // Active Directory Web Services
	Port{9390, 6}:    "otp",                  // OpenVAS Transfer Protocol
	Port{9396, 6}:    "fjinvmgr",             // Missing description for fjinvmgr
	Port{9396, 17}:   "fjinvmgr",             // Missing description for fjinvmgr
	Port{9397, 6}:    "mpidcagt",             // Missing description for mpidcagt
	Port{9397, 17}:   "mpidcagt",             // MpIdcAgt
	Port{9400, 6}:    "sec-t4net-srv",        // Samsung Twain for Network Server
	Port{9400, 17}:   "sec-t4net-srv",        // Samsung Twain for Network Server
	Port{9401, 6}:    "sec-t4net-clt",        // Samsung Twain for Network Client
	Port{9401, 17}:   "sec-t4net-clt",        // Samsung Twain for Network Client
	Port{9402, 6}:    "sec-pc2fax-srv",       // Samsung PC2FAX for Network Server
	Port{9402, 17}:   "sec-pc2fax-srv",       // Samsung PC2FAX for Network Server
	Port{9418, 6}:    "git",                  // Git revision control system | git pack transfer service
	Port{9418, 17}:   "git",                  // git pack transfer service
	Port{9443, 6}:    "tungsten-https",       // WSO2 Tungsten HTTPS
	Port{9443, 17}:   "tungsten-https",       // WSO2 Tungsten HTTPS
	Port{9444, 6}:    "wso2esb-console",      // WSO2 ESB Administration Console HTTPS
	Port{9444, 17}:   "wso2esb-console",      // WSO2 ESB Administration Console HTTPS
	Port{9445, 6}:    "mindarray-ca",         // MindArray Systems Console Agent
	Port{9450, 6}:    "sntlkeyssrvr",         // Sentinel Keys Server
	Port{9450, 17}:   "sntlkeyssrvr",         // Sentinel Keys Server
	Port{9500, 6}:    "ismserver",            // Missing description for ismserver
	Port{9500, 17}:   "ismserver",            // Missing description for ismserver
	Port{9522, 6}:    "sma-spw",              // SMA Speedwire
	Port{9535, 6}:    "man",                  // mngsuite | Management Suite Remote Control
	Port{9535, 17}:   "man",                  // Missing description for man
	Port{9536, 6}:    "laes-bf",              // Surveillance buffering function
	Port{9536, 17}:   "laes-bf",              // Surveillance buffering function
	Port{9555, 6}:    "trispen-sra",          // Trispen Secure Remote Access
	Port{9555, 17}:   "trispen-sra",          // Trispen Secure Remote Access
	Port{9592, 6}:    "ldgateway",            // LANDesk Gateway
	Port{9592, 17}:   "ldgateway",            // LANDesk Gateway
	Port{9593, 6}:    "cba8",                 // LANDesk Management Agent (cba8)
	Port{9593, 17}:   "cba8",                 // LANDesk Management Agent (cba8)
	Port{9594, 6}:    "msgsys",               // Message System
	Port{9594, 17}:   "msgsys",               // Message System
	Port{9595, 6}:    "pds",                  // Ping Discovery System | Ping Discovery Service
	Port{9595, 17}:   "pds",                  // Ping Discovery System
	Port{9596, 6}:    "mercury-disc",         // Mercury Discovery
	Port{9596, 17}:   "mercury-disc",         // Mercury Discovery
	Port{9597, 6}:    "pd-admin",             // PD Administration
	Port{9597, 17}:   "pd-admin",             // PD Administration
	Port{9598, 6}:    "vscp",                 // Very Simple Ctrl Protocol
	Port{9598, 17}:   "vscp",                 // Very Simple Ctrl Protocol
	Port{9599, 6}:    "robix",                // Missing description for robix
	Port{9599, 17}:   "robix",                // Robix
	Port{9600, 6}:    "micromuse-ncpw",       // Missing description for micromuse-ncpw
	Port{9600, 17}:   "micromuse-ncpw",       // MICROMUSE-NCPW
	Port{9612, 6}:    "streamcomm-ds",        // StreamComm User Directory
	Port{9612, 17}:   "streamcomm-ds",        // StreamComm User Directory
	Port{9614, 6}:    "iadt-tls",             // iADT Protocol over TLS
	Port{9616, 6}:    "erunbook_agent",       // erunbook-agent | eRunbook Agent
	Port{9617, 6}:    "erunbook_server",      // erunbook-server | eRunbook Server
	Port{9618, 6}:    "condor",               // Condor Collector Service
	Port{9618, 17}:   "condor",               // Condor Collector Service
	Port{9628, 6}:    "odbcpathway",          // ODBC Pathway Service
	Port{9628, 17}:   "odbcpathway",          // ODBC Pathway Service
	Port{9629, 6}:    "uniport",              // UniPort SSO Controller
	Port{9629, 17}:   "uniport",              // UniPort SSO Controller
	Port{9630, 6}:    "peoctlr",              // Peovica Controller
	Port{9631, 6}:    "peocoll",              // Peovica Collector
	Port{9632, 6}:    "mc-comm",              // Mobile-C Communications
	Port{9632, 17}:   "mc-comm",              // Mobile-C Communications
	Port{9640, 6}:    "pqsflows",             // ProQueSys Flows Service
	Port{9666, 6}:    "zoomcp",               // Zoom Control Panel Game Server Management
	Port{9667, 6}:    "xmms2",                // Cross-platform Music Multiplexing System
	Port{9667, 17}:   "xmms2",                // Cross-platform Music Multiplexing System
	Port{9668, 6}:    "tec5-sdctp",           // tec5 Spectral Device Control Protocol
	Port{9668, 17}:   "tec5-sdctp",           // tec5 Spectral Device Control Protocol
	Port{9694, 6}:    "client-wakeup",        // T-Mobile Client Wakeup Message
	Port{9694, 17}:   "client-wakeup",        // T-Mobile Client Wakeup Message
	Port{9695, 6}:    "ccnx",                 // Content Centric Networking
	Port{9695, 17}:   "ccnx",                 // Content Centric Networking
	Port{9700, 6}:    "board-roar",           // Board M.I.T. Service
	Port{9700, 17}:   "board-roar",           // Board M.I.T. Service
	Port{9747, 6}:    "l5nas-parchan",        // L5NAS Parallel Channel
	Port{9747, 17}:   "l5nas-parchan",        // L5NAS Parallel Channel
	Port{9750, 6}:    "board-voip",           // Board M.I.T. Synchronous Collaboration
	Port{9750, 17}:   "board-voip",           // Board M.I.T. Synchronous Collaboration
	Port{9753, 6}:    "rasadv",               // Missing description for rasadv
	Port{9753, 17}:   "rasadv",               // Missing description for rasadv
	Port{9762, 6}:    "tungsten-http",        // WSO2 Tungsten HTTP
	Port{9762, 17}:   "tungsten-http",        // WSO2 Tungsten HTTP
	Port{9800, 6}:    "davsrc",               // WebDav Source Port
	Port{9800, 17}:   "davsrc",               // WebDav Source Port
	Port{9801, 6}:    "sstp-2",               // Sakura Script Transfer Protocol-2
	Port{9801, 17}:   "sstp-2",               // Sakura Script Transfer Protocol-2
	Port{9802, 6}:    "davsrcs",              // WebDAV Source TLS SSL
	Port{9802, 17}:   "davsrcs",              // WebDAV Source TLS SSL
	Port{9875, 6}:    "sapv1",                // Session Announcement v1
	Port{9875, 17}:   "sapv1",                // Session Announcement v1
	Port{9876, 6}:    "sd",                   // Session Director
	Port{9876, 17}:   "sd",                   // Session Director
	Port{9878, 6}:    "kca-service",          // The KX509 Kerberized Certificate Issuance Protocol in Use in 2012
	Port{9888, 6}:    "cyborg-systems",       // CYBORG Systems
	Port{9888, 17}:   "cyborg-systems",       // CYBORG Systems
	Port{9889, 6}:    "gt-proxy",             // Port for Cable network related data proxy or repeater
	Port{9889, 17}:   "gt-proxy",             // Port for Cable network related data proxy or repeater
	Port{9898, 6}:    "monkeycom",            // Missing description for monkeycom
	Port{9898, 17}:   "monkeycom",            // MonkeyCom
	Port{9899, 132}:  "sctp-tunneling",       // SCTP Tunneling (misconfiguration) | SCTP TUNNELING
	Port{9899, 6}:    "sctp-tunneling",       // SCTP TUNNELING
	Port{9899, 17}:   "sctp-tunneling",       // SCTP Tunneling
	Port{9900, 132}:  "iua",                  // Missing description for iua
	Port{9900, 6}:    "iua",                  // IUA
	Port{9900, 17}:   "iua",                  // IUA
	Port{9901, 132}:  "enrp-sctp",            // enrp | ENRP server channel | enrp server channel
	Port{9901, 17}:   "enrp",                 // ENRP server channel
	Port{9902, 132}:  "enrp-sctp-tls",        // ENRP TLS server channel | enrp tls server channel
	Port{9903, 6}:    "multicast-ping",       // Multicast Ping Protocol
	Port{9909, 6}:    "domaintime",           // Missing description for domaintime
	Port{9909, 17}:   "domaintime",           // Missing description for domaintime
	Port{9911, 6}:    "sype-transport",       // SYPECom Transport Protocol
	Port{9911, 17}:   "sype-transport",       // SYPECom Transport Protocol
	Port{9925, 6}:    "xybrid-cloud",         // XYBRID Cloud
	Port{9929, 6}:    "nping-echo",           // Nping echo server mode - http:  nmap.org book nping-man-echo-mode.html - The port frequency is made up to keep it (barely) in top 1000 TCP
	Port{9950, 6}:    "apc-9950",             // APC 9950
	Port{9950, 17}:   "apc-9950",             // APC 9950
	Port{9951, 6}:    "apc-9951",             // APC 9951
	Port{9951, 17}:   "apc-9951",             // APC 9951
	Port{9952, 6}:    "apc-9952",             // APC 9952
	Port{9952, 17}:   "apc-9952",             // APC 9952
	Port{9953, 6}:    "acis",                 // 9953
	Port{9953, 17}:   "acis",                 // 9953
	Port{9954, 6}:    "hinp",                 // HaloteC Instrument Network Protocol
	Port{9955, 6}:    "alljoyn-stm",          // alljoyn-mcm | Contact Port for AllJoyn standard messaging | Contact Port for AllJoyn multiplexed constrained messaging
	Port{9956, 6}:    "alljoyn",              // Alljoyn Name Service
	Port{9966, 6}:    "odnsp",                // OKI Data Network Setting Protocol
	Port{9966, 17}:   "odnsp",                // OKI Data Network Setting Protocol
	Port{9978, 6}:    "xybrid-rt",            // XYBRID RT Server
	Port{9979, 6}:    "visweather",           // Valley Information Systems Weather station data
	Port{9981, 6}:    "pumpkindb",            // Event sourcing database engine with a built-in programming language
	Port{9987, 6}:    "dsm-scm-target",       // DSM SCM Target Interface
	Port{9987, 17}:   "dsm-scm-target",       // DSM SCM Target Interface
	Port{9988, 6}:    "nsesrvr",              // Software Essentials Secure HTTP server
	Port{9990, 6}:    "osm-appsrvr",          // OSM Applet Server
	Port{9990, 17}:   "osm-appsrvr",          // OSM Applet Server
	Port{9991, 6}:    "issa",                 // osm-oev | ISS System Scanner Agent | OSM Event Server
	Port{9991, 17}:   "osm-oev",              // OSM Event Server
	Port{9992, 6}:    "issc",                 // palace-1 | ISS System Scanner Console | OnLive-1
	Port{9992, 17}:   "palace-1",             // OnLive-1
	Port{9993, 6}:    "palace-2",             // OnLive-2
	Port{9993, 17}:   "palace-2",             // OnLive-2
	Port{9994, 6}:    "palace-3",             // OnLive-3
	Port{9994, 17}:   "palace-3",             // OnLive-3
	Port{9995, 6}:    "palace-4",             // Missing description for palace-4
	Port{9995, 17}:   "palace-4",             // Palace-4
	Port{9996, 6}:    "palace-5",             // Missing description for palace-5
	Port{9996, 17}:   "palace-5",             // Palace-5
	Port{9997, 6}:    "palace-6",             // Missing description for palace-6
	Port{9997, 17}:   "palace-6",             // Palace-6
	Port{9998, 6}:    "distinct32",           // Missing description for distinct32
	Port{9998, 17}:   "distinct32",           // Distinct32
	Port{9999, 6}:    "abyss",                // Abyss web server remote web management interface | distinct
	Port{9999, 17}:   "distinct",             // Missing description for distinct
	Port{10000, 6}:   "snet-sensor-mgmt",     // ndmp | SecureNet Pro Sensor https management server or apple airport admin | Network Data Management Protocol
	Port{10000, 17}:  "ndmp",                 // Network Data Management Protocol
	Port{10001, 6}:   "scp-config",           // SCP Configuration
	Port{10001, 17}:  "scp-config",           // SCP Configuration
	Port{10002, 6}:   "documentum",           // EMC-Documentum Content Server Product
	Port{10002, 17}:  "documentum",           // EMC-Documentum Content Server Product
	Port{10003, 6}:   "documentum_s",         // documentum-s | EMC-Documentum Content Server Product
	Port{10003, 17}:  "documentum_s",         // EMC-Documentum Content Server Product
	Port{10004, 6}:   "emcrmirccd",           // EMC Replication Manager Client
	Port{10005, 6}:   "stel",                 // emcrmird | Secure telnet | EMC Replication Manager Server
	Port{10006, 6}:   "netapp-sync",          // Sync replication protocol among different NetApp platforms
	Port{10007, 6}:   "mvs-capacity",         // MVS Capacity
	Port{10007, 17}:  "mvs-capacity",         // MVS Capacity
	Port{10008, 6}:   "octopus",              // Octopus Multiplexer
	Port{10008, 17}:  "octopus",              // Octopus Multiplexer
	Port{10009, 6}:   "swdtp-sv",             // Systemwalker Desktop Patrol
	Port{10009, 17}:  "swdtp-sv",             // Systemwalker Desktop Patrol
	Port{10010, 6}:   "rxapi",                // ooRexx rxapi services
	Port{10020, 6}:   "abb-hw",               // Hardware configuration and maintenance
	Port{10023, 6}:   "cefd-vmp",             // Comtech EF-Data's Vipersat Management Protocol
	Port{10050, 6}:   "zabbix-agent",         // Zabbix Agent
	Port{10050, 17}:  "zabbix-agent",         // Zabbix Agent
	Port{10051, 6}:   "zabbix-trapper",       // Zabbix Trapper
	Port{10051, 17}:  "zabbix-trapper",       // Zabbix Trapper
	Port{10055, 6}:   "qptlmd",               // Quantapoint FLEXlm Licensing Service
	Port{10080, 6}:   "amanda",               // Missing description for amanda
	Port{10080, 17}:  "amanda",               // Amanda Backup Util
	Port{10081, 6}:   "famdc",                // FAM Archive Server
	Port{10081, 17}:  "famdc",                // FAM Archive Server
	Port{10082, 6}:   "amandaidx",            // Amanda indexing
	Port{10083, 6}:   "amidxtape",            // Amanda tape indexing
	Port{10100, 6}:   "itap-ddtp",            // VERITAS ITAP DDTP
	Port{10100, 17}:  "itap-ddtp",            // VERITAS ITAP DDTP
	Port{10101, 6}:   "ezmeeting-2",          // eZmeeting
	Port{10101, 17}:  "ezmeeting-2",          // eZmeeting
	Port{10102, 6}:   "ezproxy-2",            // eZproxy
	Port{10102, 17}:  "ezproxy-2",            // eZproxy
	Port{10103, 6}:   "ezrelay",              // Missing description for ezrelay
	Port{10103, 17}:  "ezrelay",              // eZrelay
	Port{10104, 6}:   "swdtp",                // Systemwalker Desktop Patrol
	Port{10104, 17}:  "swdtp",                // Systemwalker Desktop Patrol
	Port{10107, 6}:   "bctp-server",          // VERITAS BCTP, server
	Port{10107, 17}:  "bctp-server",          // VERITAS BCTP, server
	Port{10110, 6}:   "nmea-0183",            // NMEA-0183 Navigational Data
	Port{10110, 17}:  "nmea-0183",            // NMEA-0183 Navigational Data
	Port{10111, 6}:   "nmea-onenet",          // NMEA OneNet multicast messaging
	Port{10113, 6}:   "netiq-endpoint",       // NetIQ Endpoint
	Port{10113, 17}:  "netiq-endpoint",       // NetIQ Endpoint
	Port{10114, 6}:   "netiq-qcheck",         // NetIQ Qcheck
	Port{10114, 17}:  "netiq-qcheck",         // NetIQ Qcheck
	Port{10115, 6}:   "netiq-endpt",          // NetIQ Endpoint
	Port{10115, 17}:  "netiq-endpt",          // NetIQ Endpoint
	Port{10116, 6}:   "netiq-voipa",          // NetIQ VoIP Assessor
	Port{10116, 17}:  "netiq-voipa",          // NetIQ VoIP Assessor
	Port{10117, 6}:   "iqrm",                 // NetIQ IQCResource Managament Svc
	Port{10117, 17}:  "iqrm",                 // NetIQ IQCResource Managament Svc
	Port{10125, 6}:   "cimple",               // HotLink CIMple REST API
	Port{10128, 6}:   "bmc-perf-sd",          // BMC-PERFORM-SERVICE DAEMON
	Port{10128, 17}:  "bmc-perf-sd",          // BMC-PERFORM-SERVICE DAEMON
	Port{10129, 6}:   "bmc-gms",              // BMC General Manager Server
	Port{10160, 6}:   "qb-db-server",         // QB Database Server
	Port{10160, 17}:  "qb-db-server",         // QB Database Server
	Port{10161, 6}:   "snmptls",              // snmpdtls | SNMP-TLS | SNMP-DTLS
	Port{10161, 17}:  "snmpdtls",             // SNMP-DTLS
	Port{10162, 6}:   "snmptls-trap",         // snmpdtls-trap | SNMP-Trap-TLS | SNMP-Trap-DTLS
	Port{10162, 17}:  "snmpdtls-trap",        // SNMP-Trap-DTLS
	Port{10200, 6}:   "trisoap",              // Trigence AE Soap Service
	Port{10200, 17}:  "trisoap",              // Trigence AE Soap Service
	Port{10201, 6}:   "rsms",                 // rscs | Remote Server Management Service | Remote Server Control and Test Service
	Port{10201, 17}:  "rscs",                 // Remote Server Control and Test Service
	Port{10252, 6}:   "apollo-relay",         // Apollo Relay Port
	Port{10252, 17}:  "apollo-relay",         // Apollo Relay Port
	Port{10253, 6}:   "eapol-relay",          // Relay of EAPOL frames
	Port{10260, 6}:   "axis-wimp-port",       // Axis WIMP Port
	Port{10260, 17}:  "axis-wimp-port",       // Axis WIMP Port
	Port{10261, 6}:   "tile-ml",              // Tile remote machine learning
	Port{10288, 6}:   "blocks",               // Missing description for blocks
	Port{10288, 17}:  "blocks",               // Blocks
	Port{10321, 6}:   "cosir",                // Computer Op System Information Report
	Port{10439, 6}:   "bngsync",              // BalanceNG session table synchronization protocol
	Port{10500, 6}:   "hip-nat-t",            // HIP NAT-Traversal
	Port{10500, 17}:  "hip-nat-t",            // HIP NAT-Traversal
	Port{10540, 6}:   "MOS-lower",            // MOS Media Object Metadata Port
	Port{10540, 17}:  "MOS-lower",            // MOS Media Object Metadata Port
	Port{10541, 6}:   "MOS-upper",            // MOS Running Order Port
	Port{10541, 17}:  "MOS-upper",            // MOS Running Order Port
	Port{10542, 6}:   "MOS-aux",              // MOS Low Priority Port
	Port{10542, 17}:  "MOS-aux",              // MOS Low Priority Port
	Port{10543, 6}:   "MOS-soap",             // MOS SOAP Default Port
	Port{10543, 17}:  "MOS-soap",             // MOS SOAP Default Port
	Port{10544, 6}:   "MOS-soap-opt",         // MOS SOAP Optional Port
	Port{10544, 17}:  "MOS-soap-opt",         // MOS SOAP Optional Port
	Port{10548, 6}:   "serverdocs",           // Apple Document Sharing Service
	Port{10631, 6}:   "printopia",            // Printopia Serve
	Port{10800, 6}:   "gap",                  // Gestor de Acaparamiento para Pocket PCs
	Port{10800, 17}:  "gap",                  // Gestor de Acaparamiento para Pocket PCs
	Port{10805, 6}:   "lpdg",                 // LUCIA Pareja Data Group
	Port{10805, 17}:  "lpdg",                 // LUCIA Pareja Data Group
	Port{10809, 6}:   "nbd",                  // Linux Network Block Device
	Port{10810, 6}:   "nmc-disc",             // Nuance Mobile Care Discovery
	Port{10810, 17}:  "nmc-disc",             // Nuance Mobile Care Discovery
	Port{10860, 6}:   "helix",                // Helix Client Server
	Port{10860, 17}:  "helix",                // Helix Client Server
	Port{10880, 6}:   "bveapi",               // BVEssentials HTTP API
	Port{10933, 6}:   "octopustentacle",      // Listen port used by the Octopus Deploy Tentacle deployment agent
	Port{10990, 6}:   "rmiaux",               // Auxiliary RMI Port
	Port{10990, 17}:  "rmiaux",               // Auxiliary RMI Port
	Port{11000, 6}:   "irisa",                // Missing description for irisa
	Port{11000, 17}:  "irisa",                // IRISA
	Port{11001, 6}:   "metasys",              // Missing description for metasys
	Port{11001, 17}:  "metasys",              // Metasys
	Port{11095, 6}:   "weave",                // Nest device-to-device and device-to-service application protocol
	Port{11103, 6}:   "origo-sync",           // OrigoDB Server Sync Interface
	Port{11104, 6}:   "netapp-icmgmt",        // NetApp Intercluster Management
	Port{11105, 6}:   "netapp-icdata",        // NetApp Intercluster Data
	Port{11106, 6}:   "sgi-lk",               // SGI LK Licensing service
	Port{11106, 17}:  "sgi-lk",               // SGI LK Licensing service
	Port{11108, 6}:   "myq-termlink",         // Hardware Terminals Discovery and Low-Level Communication Protocol
	Port{11109, 6}:   "sgi-dmfmgr",           // Data migration facility Manager (DMF) is a browser based interface to DMF
	Port{11110, 6}:   "sgi-soap",             // Data migration facility (DMF) SOAP is a web server protocol to support remote access to DMF
	Port{11111, 6}:   "vce",                  // Viral Computing Environment (VCE)
	Port{11111, 17}:  "vce",                  // Viral Computing Environment (VCE)
	Port{11112, 6}:   "dicom",                // Missing description for dicom
	Port{11112, 17}:  "dicom",                // DICOM
	Port{11161, 6}:   "suncacao-snmp",        // sun cacao snmp access point
	Port{11161, 17}:  "suncacao-snmp",        // sun cacao snmp access point
	Port{11162, 6}:   "suncacao-jmxmp",       // sun cacao JMX-remoting access point
	Port{11162, 17}:  "suncacao-jmxmp",       // sun cacao JMX-remoting access point
	Port{11163, 6}:   "suncacao-rmi",         // sun cacao rmi registry access point
	Port{11163, 17}:  "suncacao-rmi",         // sun cacao rmi registry access point
	Port{11164, 6}:   "suncacao-csa",         // sun cacao command-streaming access point
	Port{11164, 17}:  "suncacao-csa",         // sun cacao command-streaming access point
	Port{11165, 6}:   "suncacao-websvc",      // sun cacao web service access point
	Port{11165, 17}:  "suncacao-websvc",      // sun cacao web service access point
	Port{11171, 6}:   "snss",                 // Surgical Notes Security Service Discovery (SNSS)
	Port{11171, 17}:  "snss",                 // Surgical Notes Security Service Discovery (SNSS)
	Port{11172, 6}:   "oemcacao-jmxmp",       // OEM cacao JMX-remoting access point
	Port{11173, 6}:   "t5-straton",           // Straton Runtime Programing
	Port{11174, 6}:   "oemcacao-rmi",         // OEM cacao rmi registry access point
	Port{11175, 6}:   "oemcacao-websvc",      // OEM cacao web service access point
	Port{11201, 6}:   "smsqp",                // Missing description for smsqp
	Port{11201, 17}:  "smsqp",                // Missing description for smsqp
	Port{11202, 6}:   "dcsl-backup",          // DCSL Network Backup Services
	Port{11208, 6}:   "wifree",               // WiFree Service
	Port{11208, 17}:  "wifree",               // WiFree Service
	Port{11211, 6}:   "memcache",             // Memory cache service
	Port{11211, 17}:  "memcache",             // Memory cache service
	Port{11319, 6}:   "imip",                 // Missing description for imip
	Port{11319, 17}:  "imip",                 // IMIP
	Port{11320, 6}:   "imip-channels",        // IMIP Channels Port
	Port{11320, 17}:  "imip-channels",        // IMIP Channels Port
	Port{11321, 6}:   "arena-server",         // Arena Server Listen
	Port{11321, 17}:  "arena-server",         // Arena Server Listen
	Port{11367, 6}:   "atm-uhas",             // ATM UHAS
	Port{11367, 17}:  "atm-uhas",             // ATM UHAS
	Port{11371, 6}:   "pksd",                 // hkp | PGP Public Key Server | OpenPGP HTTP Keyserver
	Port{11371, 17}:  "hkp",                  // OpenPGP HTTP Keyserver
	Port{11430, 6}:   "lsdp",                 // Lenbrook Service Discovery Protocol
	Port{11489, 6}:   "asgcypresstcps",       // ASG Cypress Secure Only
	Port{11600, 6}:   "tempest-port",         // Tempest Protocol Port
	Port{11600, 17}:  "tempest-port",         // Tempest Protocol Port
	Port{11623, 6}:   "emc-xsw-dconfig",      // EMC XtremSW distributed config
	Port{11720, 6}:   "h323callsigalt",       // h323 Call Signal Alternate | H.323 Call Control Signalling Alternate
	Port{11720, 17}:  "h323callsigalt",       // h323 Call Signal Alternate
	Port{11723, 6}:   "emc-xsw-dcache",       // EMC XtremSW distributed cache
	Port{11751, 6}:   "intrepid-ssl",         // Intrepid SSL
	Port{11751, 17}:  "intrepid-ssl",         // Intrepid SSL
	Port{11796, 6}:   "lanschool",            // lanschool-mpt | Lanschool Multipoint
	Port{11876, 6}:   "xoraya",               // X2E Xoraya Multichannel protocol
	Port{11876, 17}:  "xoraya",               // X2E Xoraya Multichannel protocol
	Port{11877, 6}:   "x2e-disc",             // X2E service discovery protocol
	Port{11877, 17}:  "x2e-disc",             // X2E service discovery protocol
	Port{11967, 6}:   "sysinfo-sp",           // SysInfo Service Protocol | SysInfo Sercice Protocol
	Port{11967, 17}:  "sysinfo-sp",           // SysInfo Sercice Protocol
	Port{11997, 132}: "wmereceiving",         // WorldMailExpress
	Port{11998, 132}: "wmedistribution",      // WorldMailExpress
	Port{11999, 132}: "wmereporting",         // WorldMailExpress
	Port{12000, 6}:   "cce4x",                // entextxid | ClearCommerce Engine 4.x (www.clearcommerce.com) | IBM Enterprise Extender SNA XID Exchange
	Port{12000, 17}:  "entextxid",            // IBM Enterprise Extender SNA XID Exchange
	Port{12001, 6}:   "entextnetwk",          // IBM Enterprise Extender SNA COS Network Priority
	Port{12001, 17}:  "entextnetwk",          // IBM Enterprise Extender SNA COS Network Priority
	Port{12002, 6}:   "entexthigh",           // IBM Enterprise Extender SNA COS High Priority
	Port{12002, 17}:  "entexthigh",           // IBM Enterprise Extender SNA COS High Priority
	Port{12003, 6}:   "entextmed",            // IBM Enterprise Extender SNA COS Medium Priority
	Port{12003, 17}:  "entextmed",            // IBM Enterprise Extender SNA COS Medium Priority
	Port{12004, 6}:   "entextlow",            // IBM Enterprise Extender SNA COS Low Priority
	Port{12004, 17}:  "entextlow",            // IBM Enterprise Extender SNA COS Low Priority
	Port{12005, 6}:   "dbisamserver1",        // DBISAM Database Server - Regular
	Port{12005, 17}:  "dbisamserver1",        // DBISAM Database Server - Regular
	Port{12006, 6}:   "dbisamserver2",        // DBISAM Database Server - Admin
	Port{12006, 17}:  "dbisamserver2",        // DBISAM Database Server - Admin
	Port{12007, 6}:   "accuracer",            // Accuracer Database System ñ Server | Accuracer Database System Server
	Port{12007, 17}:  "accuracer",            // Accuracer Database System ñ Server
	Port{12008, 6}:   "accuracer-dbms",       // Accuracer Database System ñ Admin | Accuracer Database System Admin
	Port{12008, 17}:  "accuracer-dbms",       // Accuracer Database System ñ Admin
	Port{12009, 6}:   "ghvpn",                // Green Hills VPN
	Port{12010, 6}:   "edbsrvr",              // ElevateDB Server
	Port{12012, 6}:   "vipera",               // Vipera Messaging Service
	Port{12012, 17}:  "vipera",               // Vipera Messaging Service
	Port{12013, 6}:   "vipera-ssl",           // Vipera Messaging Service over SSL Communication
	Port{12013, 17}:  "vipera-ssl",           // Vipera Messaging Service over SSL Communication
	Port{12109, 6}:   "rets-ssl",             // RETS over SSL
	Port{12109, 17}:  "rets-ssl",             // RETS over SSL
	Port{12121, 6}:   "nupaper-ss",           // NuPaper Session Service
	Port{12121, 17}:  "nupaper-ss",           // NuPaper Session Service
	Port{12168, 6}:   "cawas",                // CA Web Access Service
	Port{12168, 17}:  "cawas",                // CA Web Access Service
	Port{12172, 6}:   "hivep",                // Missing description for hivep
	Port{12172, 17}:  "hivep",                // HiveP
	Port{12300, 6}:   "linogridengine",       // LinoGrid Engine
	Port{12300, 17}:  "linogridengine",       // LinoGrid Engine
	Port{12302, 6}:   "rads",                 // Remote Administration Daemon (RAD) is a system service that offers secure, remote, programmatic access to Solaris system configuration and run-time state
	Port{12321, 6}:   "warehouse-sss",        // Warehouse Monitoring Syst SSS
	Port{12321, 17}:  "warehouse-sss",        // Warehouse Monitoring Syst SSS
	Port{12322, 6}:   "warehouse",            // Warehouse Monitoring Syst
	Port{12322, 17}:  "warehouse",            // Warehouse Monitoring Syst
	Port{12345, 6}:   "netbus",               // italk | NetBus backdoor trojan or Trend Micro Office Scan | Italk Chat System
	Port{12345, 17}:  "italk",                // Italk Chat System
	Port{12346, 6}:   "netbus",               // NetBus backdoor trojan
	Port{12753, 6}:   "tsaf",                 // tsaf port
	Port{12753, 17}:  "tsaf",                 // tsaf port
	Port{12865, 6}:   "netperf",              // control port for the netperf benchmark
	Port{13160, 6}:   "i-zipqd",              // Missing description for i-zipqd
	Port{13160, 17}:  "i-zipqd",              // I-ZIPQD
	Port{13216, 6}:   "bcslogc",              // Black Crow Software application logging
	Port{13216, 17}:  "bcslogc",              // Black Crow Software application logging
	Port{13217, 6}:   "rs-pias",              // R&S Proxy Installation Assistant Service
	Port{13217, 17}:  "rs-pias",              // R&S Proxy Installation Assistant Service
	Port{13218, 6}:   "emc-vcas-tcp",         // emc-vcas-udp | EMC Virtual CAS Service | EMV Virtual CAS Service Discovery
	Port{13218, 17}:  "emc-vcas-udp",         // EMV Virtual CAS Service Discovery
	Port{13223, 6}:   "powwow-client",        // PowWow Client
	Port{13223, 17}:  "powwow-client",        // PowWow Client
	Port{13224, 6}:   "powwow-server",        // PowWow Server
	Port{13224, 17}:  "powwow-server",        // PowWow Server
	Port{13400, 6}:   "doip-data",            // doip-disc | DoIP Data | DoIP Discovery
	Port{13701, 6}:   "netbackup",            // vmd           server
	Port{13702, 6}:   "netbackup",            // ascd          server
	Port{13705, 6}:   "netbackup",            // tl8cd         server
	Port{13706, 6}:   "netbackup",            // odld          server
	Port{13708, 6}:   "netbackup",            // vtlcd         server
	Port{13709, 6}:   "netbackup",            // ts8d          server
	Port{13710, 6}:   "netbackup",            // tc8d          server
	Port{13711, 6}:   "netbackup",            // server
	Port{13712, 6}:   "netbackup",            // tc4d          server
	Port{13713, 6}:   "netbackup",            // tl4d          server
	Port{13714, 6}:   "netbackup",            // tsdd          server
	Port{13715, 6}:   "netbackup",            // tshd          server
	Port{13716, 6}:   "netbackup",            // tlmd          server
	Port{13717, 6}:   "netbackup",            // tlhcd         server
	Port{13718, 6}:   "netbackup",            // lmfcd         server
	Port{13720, 6}:   "netbackup",            // bprd | bprd          server | BPRD Protocol (VERITAS NetBackup)
	Port{13720, 17}:  "bprd",                 // BPRD Protocol (VERITAS NetBackup)
	Port{13721, 6}:   "netbackup",            // bpdbm | bpdbm         server | BPDBM Protocol (VERITAS NetBackup)
	Port{13721, 17}:  "bpdbm",                // BPDBM Protocol (VERITAS NetBackup)
	Port{13722, 6}:   "netbackup",            // bpjava-msvc | bpjava-msvc   client | BP Java MSVC Protocol
	Port{13722, 17}:  "bpjava-msvc",          // BP Java MSVC Protocol
	Port{13724, 6}:   "vnetd",                // Veritas Network Utility
	Port{13724, 17}:  "vnetd",                // Veritas Network Utility
	Port{13782, 6}:   "netbackup",            // bpcd | bpcd          client | VERITAS NetBackup
	Port{13782, 17}:  "bpcd",                 // VERITAS NetBackup
	Port{13783, 6}:   "netbackup",            // vopied | vopied        client | VOPIED Protocol
	Port{13783, 17}:  "vopied",               // VOPIED Protocol
	Port{13785, 6}:   "nbdb",                 // NetBackup Database
	Port{13785, 17}:  "nbdb",                 // NetBackup Database
	Port{13786, 6}:   "nomdb",                // Veritas-nomdb
	Port{13786, 17}:  "nomdb",                // Veritas-nomdb
	Port{13818, 6}:   "dsmcc-config",         // DSMCC Config
	Port{13818, 17}:  "dsmcc-config",         // DSMCC Config
	Port{13819, 6}:   "dsmcc-session",        // DSMCC Session Messages
	Port{13819, 17}:  "dsmcc-session",        // DSMCC Session Messages
	Port{13820, 6}:   "dsmcc-passthru",       // DSMCC Pass-Thru Messages
	Port{13820, 17}:  "dsmcc-passthru",       // DSMCC Pass-Thru Messages
	Port{13821, 6}:   "dsmcc-download",       // DSMCC Download Protocol
	Port{13821, 17}:  "dsmcc-download",       // DSMCC Download Protocol
	Port{13822, 6}:   "dsmcc-ccp",            // DSMCC Channel Change Protocol
	Port{13822, 17}:  "dsmcc-ccp",            // DSMCC Channel Change Protocol
	Port{13823, 6}:   "bmdss",                // Blackmagic Design Streaming Server
	Port{13882, 6}:   "vunknown",             // Missing description for vunknown
	Port{13894, 6}:   "ucontrol",             // Ultimate Control communication protocol
	Port{13929, 6}:   "dta-systems",          // D-TA SYSTEMS
	Port{13929, 17}:  "dta-systems",          // D-TA SYSTEMS
	Port{13930, 6}:   "medevolve",            // MedEvolve Port Requester
	Port{14000, 6}:   "scotty-ft",            // SCOTTY High-Speed Filetransfer
	Port{14000, 17}:  "scotty-ft",            // SCOTTY High-Speed Filetransfer
	Port{14001, 132}: "sua",                  // De-Registered
	Port{14001, 6}:   "sua",                  // SUA
	Port{14001, 17}:  "sua",                  // De-Registered (2001 June 06)
	Port{14002, 6}:   "scotty-disc",          // Discovery of a SCOTTY hardware codec board
	Port{14033, 6}:   "sage-best-com1",       // sage Best! Config Server 1
	Port{14033, 17}:  "sage-best-com1",       // sage Best! Config Server 1
	Port{14034, 6}:   "sage-best-com2",       // sage Best! Config Server 2
	Port{14034, 17}:  "sage-best-com2",       // sage Best! Config Server 2
	Port{14141, 6}:   "bo2k",                 // vcs-app | Back Orifice 2K BoPeep mouse keyboard input | VCS Application
	Port{14141, 17}:  "vcs-app",              // VCS Application
	Port{14142, 6}:   "icpp",                 // IceWall Cert Protocol
	Port{14142, 17}:  "icpp",                 // IceWall Cert Protocol
	Port{14143, 6}:   "icpps",                // IceWall Cert Protocol over TLS
	Port{14145, 6}:   "gcm-app",              // GCM Application
	Port{14145, 17}:  "gcm-app",              // GCM Application
	Port{14149, 6}:   "vrts-tdd",             // Veritas Traffic Director
	Port{14149, 17}:  "vrts-tdd",             // Veritas Traffic Director
	Port{14150, 6}:   "vcscmd",               // Veritas Cluster Server Command Server
	Port{14154, 6}:   "vad",                  // Veritas Application Director
	Port{14154, 17}:  "vad",                  // Veritas Application Director
	Port{14250, 6}:   "cps",                  // Fencing Server
	Port{14250, 17}:  "cps",                  // Fencing Server
	Port{14414, 6}:   "ca-web-update",        // CA eTrust Web Update Service
	Port{14414, 17}:  "ca-web-update",        // CA eTrust Web Update Service
	Port{14500, 6}:   "xpra",                 // xpra network protocol
	Port{14936, 6}:   "hde-lcesrvr-1",        // Missing description for hde-lcesrvr-1
	Port{14936, 17}:  "hde-lcesrvr-1",        // Missing description for hde-lcesrvr-1
	Port{14937, 6}:   "hde-lcesrvr-2",        // Missing description for hde-lcesrvr-2
	Port{14937, 17}:  "hde-lcesrvr-2",        // Missing description for hde-lcesrvr-2
	Port{15000, 6}:   "hydap",                // Hypack Hydrographic Software Packages Data Acquisition | Hypack Data Aquisition
	Port{15000, 17}:  "hydap",                // Hypack Hydrographic Software Packages Data Acquisition
	Port{15002, 6}:   "onep-tls",             // Open Network Environment TLS
	Port{15118, 6}:   "v2g-secc",             // v2g Supply Equipment Communication Controller Discovery Protocol
	Port{15126, 6}:   "swgps",                // Nortel Java S WGPS Global Payment Solutions for US credit card authorizations
	Port{15151, 6}:   "bo2k",                 // Back Orifice 2K BoPeep video output
	Port{15345, 6}:   "xpilot",               // XPilot Contact Port
	Port{15345, 17}:  "xpilot",               // XPilot Contact Port
	Port{15363, 6}:   "3link",                // 3Link Negotiation
	Port{15363, 17}:  "3link",                // 3Link Negotiation
	Port{15555, 6}:   "cisco-snat",           // Cisco Stateful NAT
	Port{15555, 17}:  "cisco-snat",           // Cisco Stateful NAT
	Port{15660, 6}:   "bex-xr",               // Backup Express Restore Server
	Port{15660, 17}:  "bex-xr",               // Backup Express Restore Server
	Port{15740, 6}:   "ptp",                  // Picture Transfer Protocol
	Port{15740, 17}:  "ptp",                  // Picture Transfer Protocol
	Port{15998, 6}:   "2ping",                // 2ping Bi-Directional Ping Service
	Port{15998, 17}:  "2ping",                // 2ping Bi-Directional Ping Service
	Port{15999, 6}:   "programmar",           // ProGrammar Enterprise
	Port{16000, 6}:   "fmsas",                // Administration Server Access
	Port{16001, 6}:   "fmsascon",             // Administration Server Connector
	Port{16002, 6}:   "gsms",                 // GoodSync Mediation Service
	Port{16003, 6}:   "alfin",                // Automation and Control by REGULACE.ORG
	Port{16003, 17}:  "alfin",                // Automation and Control by REGULACE.ORG
	Port{16020, 6}:   "jwpc",                 // Filemaker Java Web Publishing Core
	Port{16021, 6}:   "jwpc-bin",             // Filemaker Java Web Publishing Core Binary
	Port{16080, 6}:   "osxwebadmin",          // Apple OS X WebAdmin
	Port{16161, 6}:   "sun-sea-port",         // Solaris SEA Port
	Port{16161, 17}:  "sun-sea-port",         // Solaris SEA Port
	Port{16162, 6}:   "solaris-audit",        // Solaris Audit - secure remote audit log
	Port{16309, 6}:   "etb4j",                // Missing description for etb4j
	Port{16309, 17}:  "etb4j",                // Missing description for etb4j
	Port{16310, 6}:   "pduncs",               // Policy Distribute, Update Notification
	Port{16310, 17}:  "pduncs",               // Policy Distribute, Update Notification
	Port{16311, 6}:   "pdefmns",              // Policy definition and update management
	Port{16311, 17}:  "pdefmns",              // Policy definition and update management
	Port{16360, 6}:   "netserialext1",        // Network Serial Extension Ports One
	Port{16360, 17}:  "netserialext1",        // Network Serial Extension Ports One
	Port{16361, 6}:   "netserialext2",        // Network Serial Extension Ports Two
	Port{16361, 17}:  "netserialext2",        // Network Serial Extension Ports Two
	Port{16367, 6}:   "netserialext3",        // Network Serial Extension Ports Three
	Port{16367, 17}:  "netserialext3",        // Network Serial Extension Ports Three
	Port{16368, 6}:   "netserialext4",        // Network Serial Extension Ports Four
	Port{16368, 17}:  "netserialext4",        // Network Serial Extension Ports Four
	Port{16384, 6}:   "connected",            // Connected Corp
	Port{16384, 17}:  "connected",            // Connected Corp
	Port{16385, 6}:   "rdgs",                 // Reliable Datagram Sockets
	Port{16444, 6}:   "overnet",              // Overnet file sharing
	Port{16444, 17}:  "overnet",              // Overnet file sharing
	Port{16619, 6}:   "xoms",                 // X509 Objects Management Service
	Port{16665, 6}:   "axon-tunnel",          // Reliable multipath data transport for high latencies
	Port{16666, 6}:   "vtp",                  // Vidder Tunnel Protocol
	Port{16789, 6}:   "cadsisvr",             // This server provides callable services to mainframe External Security Managers from any TCP IP platform
	Port{16900, 6}:   "newbay-snc-mc",        // Newbay Mobile Client Update Service
	Port{16900, 17}:  "newbay-snc-mc",        // Newbay Mobile Client Update Service
	Port{16950, 6}:   "sgcip",                // Simple Generic Client Interface Protocol
	Port{16950, 17}:  "sgcip",                // Simple Generic Client Interface Protocol
	Port{16959, 6}:   "subseven",             // Subseven trojan
	Port{16991, 6}:   "intel-rci-mp",         // Missing description for intel-rci-mp
	Port{16991, 17}:  "intel-rci-mp",         // INTEL-RCI-MP
	Port{16992, 6}:   "amt-soap-http",        // Intel(R) AMT SOAP HTTP
	Port{16992, 17}:  "amt-soap-http",        // Intel(R) AMT SOAP HTTP
	Port{16993, 6}:   "amt-soap-https",       // Intel(R) AMT SOAP HTTPS
	Port{16993, 17}:  "amt-soap-https",       // Intel(R) AMT SOAP HTTPS
	Port{16994, 6}:   "amt-redir-tcp",        // Intel(R) AMT Redirection TCP
	Port{16994, 17}:  "amt-redir-tcp",        // Intel(R) AMT Redirection TCP
	Port{16995, 6}:   "amt-redir-tls",        // Intel(R) AMT Redirection TLS
	Port{16995, 17}:  "amt-redir-tls",        // Intel(R) AMT Redirection TLS
	Port{17007, 6}:   "isode-dua",            // Missing description for isode-dua
	Port{17007, 17}:  "isode-dua",            // Missing description for isode-dua
	Port{17184, 6}:   "vestasdlp",            // Vestas Data Layer Protocol
	Port{17185, 6}:   "soundsvirtual",        // Sounds Virtual
	Port{17185, 17}:  "wdbrpc",               // vxWorks WDB remote debugging ONCRPC
	Port{17219, 6}:   "chipper",              // Missing description for chipper
	Port{17219, 17}:  "chipper",              // Chipper
	Port{17220, 6}:   "avtp",                 // IEEE 1722 Transport Protocol for Time Sensitive Applications
	Port{17221, 6}:   "avdecc",               // IEEE 1722.1 AVB Discovery, Enumeration, Connection management, and Control
	Port{17222, 6}:   "cpsp",                 // Control Plane Synchronization Protocol (SPSP)
	Port{17223, 6}:   "isa100-gci",           // ISA100 GCI is a service utilizing a common interface between an ISA100 Wireless gateway and a client application
	Port{17224, 6}:   "trdp-pd",              // Train Realtime Data Protocol (TRDP) Process Data
	Port{17225, 6}:   "trdp-md",              // Train Realtime Data Protocol (TRDP) Message Data
	Port{17234, 6}:   "integrius-stp",        // Integrius Secure Tunnel Protocol
	Port{17234, 17}:  "integrius-stp",        // Integrius Secure Tunnel Protocol
	Port{17235, 6}:   "ssh-mgmt",             // SSH Tectia Manager
	Port{17235, 17}:  "ssh-mgmt",             // SSH Tectia Manager
	Port{17300, 6}:   "kuang2",               // Kuang2 backdoor
	Port{17500, 6}:   "db-lsp",               // db-lsp-disc | Dropbox LanSync Protocol | Dropbox LanSync Discovery
	Port{17500, 17}:  "db-lsp-disc",          // Dropbox LanSync Discovery
	Port{17555, 6}:   "ailith",               // Ailith management of routers
	Port{17729, 6}:   "ea",                   // Eclipse Aviation
	Port{17729, 17}:  "ea",                   // Eclipse Aviation
	Port{17754, 6}:   "zep",                  // Encap. ZigBee Packets
	Port{17754, 17}:  "zep",                  // Encap. ZigBee Packets
	Port{17755, 6}:   "zigbee-ip",            // ZigBee IP Transport Service
	Port{17755, 17}:  "zigbee-ip",            // ZigBee IP Transport Service
	Port{17756, 6}:   "zigbee-ips",           // ZigBee IP Transport Secure Service
	Port{17756, 17}:  "zigbee-ips",           // ZigBee IP Transport Secure Service
	Port{17777, 6}:   "sw-orion",             // SolarWinds Orion
	Port{18000, 6}:   "biimenu",              // Beckman Instruments, Inc.
	Port{18000, 17}:  "biimenu",              // Beckman Instruments, Inc.
	Port{18104, 6}:   "radpdf",               // RAD PDF Service
	Port{18136, 6}:   "racf",                 // z OS Resource Access Control Facility
	Port{18181, 6}:   "opsec-cvp",            // Check Point OPSEC | OPSEC CVP
	Port{18181, 17}:  "opsec-cvp",            // OPSEC CVP
	Port{18182, 6}:   "opsec-ufp",            // Check Point OPSEC | OPSEC UFP
	Port{18182, 17}:  "opsec-ufp",            // OPSEC UFP
	Port{18183, 6}:   "opsec-sam",            // Check Point OPSEC | OPSEC SAM
	Port{18183, 17}:  "opsec-sam",            // OPSEC SAM
	Port{18184, 6}:   "opsec-lea",            // Check Point OPSEC | OPSEC LEA
	Port{18184, 17}:  "opsec-lea",            // OPSEC LEA
	Port{18185, 6}:   "opsec-omi",            // Check Point OPSEC | OPSEC OMI
	Port{18185, 17}:  "opsec-omi",            // OPSEC OMI
	Port{18186, 6}:   "ohsc",                 // Occupational Health SC | Occupational Health Sc
	Port{18186, 17}:  "ohsc",                 // Occupational Health Sc
	Port{18187, 6}:   "opsec-ela",            // Check Point OPSEC | OPSEC ELA
	Port{18187, 17}:  "opsec-ela",            // OPSEC ELA
	Port{18241, 6}:   "checkpoint-rtm",       // Check Point RTM
	Port{18241, 17}:  "checkpoint-rtm",       // Check Point RTM
	Port{18242, 6}:   "iclid",                // Checkpoint router monitoring
	Port{18243, 6}:   "clusterxl",            // Checkpoint router state backup
	Port{18262, 6}:   "gv-pf",                // GV NetConfig Service
	Port{18262, 17}:  "gv-pf",                // GV NetConfig Service
	Port{18333, 6}:   "bitcoin",              // Bitcoin crypto currency - https:  en.bitcoin.it wiki Running_Bitcoin
	Port{18463, 6}:   "ac-cluster",           // AC Cluster
	Port{18463, 17}:  "ac-cluster",           // AC Cluster
	Port{18634, 6}:   "rds-ib",               // Reliable Datagram Service
	Port{18634, 17}:  "rds-ib",               // Reliable Datagram Service
	Port{18635, 6}:   "rds-ip",               // Reliable Datagram Service over IP
	Port{18635, 17}:  "rds-ip",               // Reliable Datagram Service over IP
	Port{18668, 6}:   "vdmmesh",              // vdmmesh-disc | Manufacturing Execution Systems Mesh Communication
	Port{18769, 6}:   "ique",                 // IQue Protocol
	Port{18769, 17}:  "ique",                 // IQue Protocol
	Port{18881, 6}:   "infotos",              // Missing description for infotos
	Port{18881, 17}:  "infotos",              // Infotos
	Port{18888, 6}:   "apc-necmp",            // APCNECMP
	Port{18888, 17}:  "apc-necmp",            // APCNECMP
	Port{19000, 6}:   "igrid",                // iGrid Server
	Port{19000, 17}:  "igrid",                // iGrid Server
	Port{19007, 6}:   "scintilla",            // Scintilla protocol for device services
	Port{19020, 6}:   "j-link",               // J-Link TCP IP Protocol
	Port{19150, 6}:   "gkrellm",              // GKrellM remote system activity meter daemon
	Port{19191, 6}:   "opsec-uaa",            // OPSEC UAA
	Port{19191, 17}:  "opsec-uaa",            // OPSEC UAA
	Port{19194, 6}:   "ua-secureagent",       // UserAuthority SecureAgent
	Port{19194, 17}:  "ua-secureagent",       // UserAuthority SecureAgent
	Port{19220, 6}:   "cora",                 // cora-disc | Client Connection Management and Data Exchange Service | Discovery for Client Connection Management and Data Exchange Service
	Port{19283, 6}:   "keysrvr",              // Key Server for SASSAFRAS
	Port{19283, 17}:  "keysrvr",              // Key Server for SASSAFRAS
	Port{19315, 6}:   "keyshadow",            // Key Shadow for SASSAFRAS
	Port{19315, 17}:  "keyshadow",            // Key Shadow for SASSAFRAS
	Port{19333, 6}:   "litecoin",             // Litecoin crypto currency testnet - https:  litecoin.info Litecoin.conf
	Port{19398, 6}:   "mtrgtrans",            // Missing description for mtrgtrans
	Port{19398, 17}:  "mtrgtrans",            // Missing description for mtrgtrans
	Port{19410, 6}:   "hp-sco",               // Missing description for hp-sco
	Port{19410, 17}:  "hp-sco",               // Missing description for hp-sco
	Port{19411, 6}:   "hp-sca",               // Missing description for hp-sca
	Port{19411, 17}:  "hp-sca",               // Missing description for hp-sca
	Port{19412, 6}:   "hp-sessmon",           // Missing description for hp-sessmon
	Port{19412, 17}:  "hp-sessmon",           // HP-SESSMON
	Port{19539, 6}:   "fxuptp",               // Missing description for fxuptp
	Port{19539, 17}:  "fxuptp",               // FXUPTP
	Port{19540, 6}:   "sxuptp",               // Missing description for sxuptp
	Port{19540, 17}:  "sxuptp",               // SXUPTP
	Port{19541, 6}:   "jcp",                  // JCP Client
	Port{19541, 17}:  "jcp",                  // JCP Client
	Port{19788, 6}:   "mle",                  // Mesh Link Establishment
	Port{19998, 6}:   "iec-104-sec",          // IEC 60870-5-104 process control - secure
	Port{19999, 6}:   "dnp-sec",              // Distributed Network Protocol - Secure
	Port{19999, 17}:  "dnp-sec",              // Distributed Network Protocol - Secure
	Port{20000, 6}:   "dnp",                  // Missing description for dnp
	Port{20000, 17}:  "dnp",                  // DNP
	Port{20001, 6}:   "microsan",             // Missing description for microsan
	Port{20001, 17}:  "microsan",             // MicroSAN
	Port{20002, 6}:   "commtact-http",        // Commtact HTTP
	Port{20002, 17}:  "commtact-http",        // Commtact HTTP
	Port{20003, 6}:   "commtact-https",       // Commtact HTTPS
	Port{20003, 17}:  "commtact-https",       // Commtact HTTPS
	Port{20005, 6}:   "btx",                  // openwebnet | xcept4 (Interacts with German Telekom's CEPT videotext service) | OpenWebNet protocol for electric network
	Port{20005, 17}:  "openwebnet",           // OpenWebNet protocol for electric network
	Port{20012, 6}:   "ss-idi-disc",          // Samsung Interdevice Interaction discovery
	Port{20012, 17}:  "ss-idi-disc",          // Samsung Interdevice Interaction discovery
	Port{20013, 6}:   "ss-idi",               // Samsung Interdevice Interaction
	Port{20014, 6}:   "opendeploy",           // OpenDeploy Listener
	Port{20014, 17}:  "opendeploy",           // OpenDeploy Listener
	Port{20031, 17}:  "bakbonenetvault",      // BakBone NetVault primary communications port
	Port{20034, 6}:   "nburn_id",             // nburn-id | NetBurner ID Port
	Port{20034, 17}:  "nburn_id",             // NetBurner ID Port
	Port{20046, 6}:   "tmophl7mts",           // TMOP HL7 Message Transfer Service
	Port{20046, 17}:  "tmophl7mts",           // TMOP HL7 Message Transfer Service
	Port{20048, 6}:   "mountd",               // NFS mount protocol
	Port{20048, 17}:  "mountd",               // NFS mount protocol
	Port{20049, 132}: "nfsrdma",              // Network File System (NFS) over RDMA
	Port{20049, 6}:   "nfsrdma",              // Network File System (NFS) over RDMA
	Port{20049, 17}:  "nfsrdma",              // Network File System (NFS) over RDMA
	Port{20057, 6}:   "avesterra",            // AvesTerra Hypergraph Transfer Protocol (HGTP)
	Port{20167, 6}:   "tolfab",               // TOLfab Data Change
	Port{20167, 17}:  "tolfab",               // TOLfab Data Change
	Port{20202, 6}:   "ipdtp-port",           // IPD Tunneling Port
	Port{20202, 17}:  "ipdtp-port",           // IPD Tunneling Port
	Port{20222, 6}:   "ipulse-ics",           // Missing description for ipulse-ics
	Port{20222, 17}:  "ipulse-ics",           // iPulse-ICS
	Port{20480, 6}:   "emwavemsg",            // emWave Message Service
	Port{20480, 17}:  "emwavemsg",            // emWave Message Service
	Port{20670, 6}:   "track",                // Missing description for track
	Port{20670, 17}:  "track",                // Track
	Port{20999, 6}:   "athand-mmp",           // At Hand MMP | AT Hand MMP
	Port{20999, 17}:  "athand-mmp",           // AT Hand MMP
	Port{21000, 6}:   "irtrans",              // IRTrans Control
	Port{21000, 17}:  "irtrans",              // IRTrans Control
	Port{21010, 6}:   "notezilla-lan",        // Notezilla.Lan Server
	Port{21201, 6}:   "memcachedb",           // Missing description for memcachedb
	Port{21212, 6}:   "trinket-agent",        // Distributed artificial intelligence
	Port{21221, 6}:   "aigairserver",         // Services for Air Server
	Port{21553, 6}:   "rdm-tfs",              // Raima RDM TFS
	Port{21554, 6}:   "dfserver",             // MineScape Design File Server
	Port{21554, 17}:  "dfserver",             // MineScape Design File Server
	Port{21590, 6}:   "vofr-gateway",         // VoFR Gateway
	Port{21590, 17}:  "vofr-gateway",         // VoFR Gateway
	Port{21800, 6}:   "tvpm",                 // TVNC Pro Multiplexing
	Port{21800, 17}:  "tvpm",                 // TVNC Pro Multiplexing
	Port{21845, 6}:   "webphone",             // Missing description for webphone
	Port{21845, 17}:  "webphone",             // Missing description for webphone
	Port{21846, 6}:   "netspeak-is",          // NetSpeak Corp. Directory Services
	Port{21846, 17}:  "netspeak-is",          // NetSpeak Corp. Directory Services
	Port{21847, 6}:   "netspeak-cs",          // NetSpeak Corp. Connection Services
	Port{21847, 17}:  "netspeak-cs",          // NetSpeak Corp. Connection Services
	Port{21848, 6}:   "netspeak-acd",         // NetSpeak Corp. Automatic Call Distribution
	Port{21848, 17}:  "netspeak-acd",         // NetSpeak Corp. Automatic Call Distribution
	Port{21849, 6}:   "netspeak-cps",         // NetSpeak Corp. Credit Processing System
	Port{21849, 17}:  "netspeak-cps",         // NetSpeak Corp. Credit Processing System
	Port{22000, 6}:   "snapenetio",           // Missing description for snapenetio
	Port{22000, 17}:  "snapenetio",           // SNAPenetIO
	Port{22001, 6}:   "optocontrol",          // Missing description for optocontrol
	Port{22001, 17}:  "optocontrol",          // OptoControl
	Port{22002, 6}:   "optohost002",          // Opto Host Port 2
	Port{22002, 17}:  "optohost002",          // Opto Host Port 2
	Port{22003, 6}:   "optohost003",          // Opto Host Port 3
	Port{22003, 17}:  "optohost003",          // Opto Host Port 3
	Port{22004, 6}:   "optohost004",          // Opto Host Port 4
	Port{22004, 17}:  "optohost004",          // Opto Host Port 4
	Port{22005, 6}:   "optohost004",          // Opto Host Port 5
	Port{22005, 17}:  "optohost004",          // Opto Host Port 5
	Port{22125, 6}:   "dcap",                 // dCache Access Protocol
	Port{22128, 6}:   "gsidcap",              // GSI dCache Access Protocol
	Port{22222, 6}:   "easyengine",           // EasyEngine is CLI tool to manage WordPress Sites on Nginx server
	Port{22273, 6}:   "wnn6",                 // Wnn6 (Japanese input)
	Port{22273, 17}:  "wnn6",                 // Missing description for wnn6
	Port{22289, 6}:   "wnn6_Cn",              // Wnn6 (Chinese input)
	Port{22305, 6}:   "wnn6_Kr",              // cis | Wnn6 (Korean input) | CompactIS Tunnel
	Port{22305, 17}:  "cis",                  // CompactIS Tunnel
	Port{22321, 6}:   "wnn6_Tw",              // Wnn6 (Taiwanse input)
	Port{22335, 6}:   "shrewd-control",       // shrewd-stream | Initium Labs Security and Automation Control | Initium Labs Security and Automation Streaming
	Port{22343, 6}:   "cis-secure",           // CompactIS Secure Tunnel
	Port{22343, 17}:  "cis-secure",           // CompactIS Secure Tunnel
	Port{22347, 6}:   "WibuKey",              // WibuKey Standard WkLan
	Port{22347, 17}:  "WibuKey",              // WibuKey Standard WkLan
	Port{22350, 6}:   "CodeMeter",            // CodeMeter Standard
	Port{22350, 17}:  "CodeMeter",            // CodeMeter Standard
	Port{22351, 6}:   "codemeter-cmwan",      // TPC IP requests of copy protection software to a server
	Port{22370, 6}:   "hpnpd",                // Hewlett-Packard Network Printer daemon
	Port{22370, 17}:  "hpnpd",                // Hewlett-Packard Network Printer daemon
	Port{22537, 6}:   "caldsoft-backup",      // CaldSoft Backup server file transfer
	Port{22555, 6}:   "vocaltec-wconf",       // vocaltec-phone | Vocaltec Web Conference | Vocaltec Internet Phone
	Port{22555, 17}:  "vocaltec-phone",       // Vocaltec Internet Phone
	Port{22763, 6}:   "talikaserver",         // Talika Main Server
	Port{22763, 17}:  "talikaserver",         // Talika Main Server
	Port{22800, 6}:   "aws-brf",              // Telerate Information Platform LAN
	Port{22800, 17}:  "aws-brf",              // Telerate Information Platform LAN
	Port{22951, 6}:   "brf-gw",               // Telerate Information Platform WAN
	Port{22951, 17}:  "brf-gw",               // Telerate Information Platform WAN
	Port{23000, 6}:   "inovaport1",           // Inova LightLink Server Type 1
	Port{23000, 17}:  "inovaport1",           // Inova LightLink Server Type 1
	Port{23001, 6}:   "inovaport2",           // Inova LightLink Server Type 2
	Port{23001, 17}:  "inovaport2",           // Inova LightLink Server Type 2
	Port{23002, 6}:   "inovaport3",           // Inova LightLink Server Type 3
	Port{23002, 17}:  "inovaport3",           // Inova LightLink Server Type 3
	Port{23003, 6}:   "inovaport4",           // Inova LightLink Server Type 4
	Port{23003, 17}:  "inovaport4",           // Inova LightLink Server Type 4
	Port{23004, 6}:   "inovaport5",           // Inova LightLink Server Type 5
	Port{23004, 17}:  "inovaport5",           // Inova LightLink Server Type 5
	Port{23005, 6}:   "inovaport6",           // Inova LightLink Server Type 6
	Port{23005, 17}:  "inovaport6",           // Inova LightLink Server Type 6
	Port{23053, 6}:   "gntp",                 // Generic Notification Transport Protocol
	Port{23272, 6}:   "s102",                 // S102 application
	Port{23272, 17}:  "s102",                 // S102 application
	Port{23294, 6}:   "5afe-dir",             // 5afe-disc | 5AFE SDN Directory | 5AFE SDN Directory discovery
	Port{23333, 6}:   "elxmgmt",              // Emulex HBAnyware Remote Management
	Port{23333, 17}:  "elxmgmt",              // Emulex HBAnyware Remote Management
	Port{23400, 6}:   "novar-dbase",          // Novar Data
	Port{23400, 17}:  "novar-dbase",          // Novar Data
	Port{23401, 6}:   "novar-alarm",          // Novar Alarm
	Port{23401, 17}:  "novar-alarm",          // Novar Alarm
	Port{23402, 6}:   "novar-global",         // Novar Global
	Port{23402, 17}:  "novar-global",         // Novar Global
	Port{23456, 6}:   "aequus",               // Aequus Service
	Port{23457, 6}:   "aequus-alt",           // Aequus Service Mgmt
	Port{23546, 6}:   "areaguard-neo",        // AreaGuard Neo - WebServer
	Port{24000, 6}:   "med-ltp",              // Missing description for med-ltp
	Port{24000, 17}:  "med-ltp",              // Missing description for med-ltp
	Port{24001, 6}:   "med-fsp-rx",           // Missing description for med-fsp-rx
	Port{24001, 17}:  "med-fsp-rx",           // Missing description for med-fsp-rx
	Port{24002, 6}:   "med-fsp-tx",           // Missing description for med-fsp-tx
	Port{24002, 17}:  "med-fsp-tx",           // Missing description for med-fsp-tx
	Port{24003, 6}:   "med-supp",             // Missing description for med-supp
	Port{24003, 17}:  "med-supp",             // Missing description for med-supp
	Port{24004, 6}:   "med-ovw",              // Missing description for med-ovw
	Port{24004, 17}:  "med-ovw",              // Missing description for med-ovw
	Port{24005, 6}:   "med-ci",               // Missing description for med-ci
	Port{24005, 17}:  "med-ci",               // Missing description for med-ci
	Port{24006, 6}:   "med-net-svc",          // Missing description for med-net-svc
	Port{24006, 17}:  "med-net-svc",          // Missing description for med-net-svc
	Port{24242, 6}:   "filesphere",           // Missing description for filesphere
	Port{24242, 17}:  "filesphere",           // fileSphere
	Port{24249, 6}:   "vista-4gl",            // Vista 4GL
	Port{24249, 17}:  "vista-4gl",            // Vista 4GL
	Port{24321, 6}:   "ild",                  // Isolv Local Directory
	Port{24321, 17}:  "ild",                  // Isolv Local Directory
	Port{24322, 6}:   "hid",                  // Transport of Human Interface Device data streams
	Port{24386, 6}:   "intel_rci",            // intel-rci | Intel RCI
	Port{24386, 17}:  "intel_rci",            // Intel RCI
	Port{24465, 6}:   "tonidods",             // Tonido Domain Server
	Port{24465, 17}:  "tonidods",             // Tonido Domain Server
	Port{24554, 6}:   "binkp",                // Missing description for binkp
	Port{24554, 17}:  "binkp",                // BINKP
	Port{24577, 6}:   "bilobit",              // bilobit-update | bilobit Service | bilobit Service Update
	Port{24666, 6}:   "sdtvwcam",             // Service used by SmarDTV to communicate between a CAM and a second screen application
	Port{24676, 6}:   "canditv",              // Canditv Message Service
	Port{24676, 17}:  "canditv",              // Canditv Message Service
	Port{24677, 6}:   "flashfiler",           // Missing description for flashfiler
	Port{24677, 17}:  "flashfiler",           // FlashFiler
	Port{24678, 6}:   "proactivate",          // Turbopower Proactivate
	Port{24678, 17}:  "proactivate",          // Turbopower Proactivate
	Port{24680, 6}:   "tcc-http",             // TCC User HTTP Service
	Port{24680, 17}:  "tcc-http",             // TCC User HTTP Service
	Port{24754, 6}:   "cslg",                 // Citrix StorageLink Gateway
	Port{24850, 6}:   "assoc-disc",           // Device Association Discovery
	Port{24922, 6}:   "find",                 // Find Identification of Network Devices
	Port{24922, 17}:  "find",                 // Find Identification of Network Devices
	Port{25000, 6}:   "icl-twobase1",         // Missing description for icl-twobase1
	Port{25000, 17}:  "icl-twobase1",         // Missing description for icl-twobase1
	Port{25001, 6}:   "icl-twobase2",         // Missing description for icl-twobase2
	Port{25001, 17}:  "icl-twobase2",         // Missing description for icl-twobase2
	Port{25002, 6}:   "icl-twobase3",         // Missing description for icl-twobase3
	Port{25002, 17}:  "icl-twobase3",         // Missing description for icl-twobase3
	Port{25003, 6}:   "icl-twobase4",         // Missing description for icl-twobase4
	Port{25003, 17}:  "icl-twobase4",         // Missing description for icl-twobase4
	Port{25004, 6}:   "icl-twobase5",         // Missing description for icl-twobase5
	Port{25004, 17}:  "icl-twobase5",         // Missing description for icl-twobase5
	Port{25005, 6}:   "icl-twobase6",         // Missing description for icl-twobase6
	Port{25005, 17}:  "icl-twobase6",         // Missing description for icl-twobase6
	Port{25006, 6}:   "icl-twobase7",         // Missing description for icl-twobase7
	Port{25006, 17}:  "icl-twobase7",         // Missing description for icl-twobase7
	Port{25007, 6}:   "icl-twobase8",         // Missing description for icl-twobase8
	Port{25007, 17}:  "icl-twobase8",         // Missing description for icl-twobase8
	Port{25008, 6}:   "icl-twobase9",         // Missing description for icl-twobase9
	Port{25008, 17}:  "icl-twobase9",         // Missing description for icl-twobase9
	Port{25009, 6}:   "icl-twobase10",        // Missing description for icl-twobase10
	Port{25009, 17}:  "icl-twobase10",        // Missing description for icl-twobase10
	Port{25471, 6}:   "rna",                  // RNSAP User Adaptation for Iurh
	Port{25565, 6}:   "minecraft",            // A video game - http:  en.wikipedia.org wiki Minecraft
	Port{25576, 6}:   "sauterdongle",         // Sauter Dongle
	Port{25604, 6}:   "idtp",                 // Identifier Tracing Protocol
	Port{25793, 6}:   "vocaltec-hos",         // Vocaltec Address Server
	Port{25793, 17}:  "vocaltec-hos",         // Vocaltec Address Server
	Port{25900, 6}:   "tasp-net",             // TASP Network Comm
	Port{25900, 17}:  "tasp-net",             // TASP Network Comm
	Port{25901, 6}:   "niobserver",           // Missing description for niobserver
	Port{25901, 17}:  "niobserver",           // NIObserver
	Port{25902, 6}:   "nilinkanalyst",        // Missing description for nilinkanalyst
	Port{25902, 17}:  "nilinkanalyst",        // NILinkAnalyst
	Port{25903, 6}:   "niprobe",              // Missing description for niprobe
	Port{25903, 17}:  "niprobe",              // NIProbe
	Port{25954, 6}:   "bf-game",              // Bitfighter game server
	Port{25955, 6}:   "bf-master",            // Bitfighter master server
	Port{26000, 6}:   "quake",                // Missing description for quake
	Port{26000, 17}:  "quake",                // Quake game server
	Port{26133, 6}:   "scscp",                // Symbolic Computation Software Composability Protocol
	Port{26133, 17}:  "scscp",                // Symbolic Computation Software Composability Protocol
	Port{26208, 6}:   "wnn6_DS",              // Wnn6 (Dserver) | wnn6-ds
	Port{26208, 17}:  "wnn6-ds",              // Missing description for wnn6-ds
	Port{26257, 6}:   "cockroach",            // CockroachDB
	Port{26260, 6}:   "ezproxy",              // Missing description for ezproxy
	Port{26260, 17}:  "ezproxy",              // eZproxy
	Port{26261, 6}:   "ezmeeting",            // Missing description for ezmeeting
	Port{26261, 17}:  "ezmeeting",            // eZmeeting
	Port{26262, 6}:   "k3software-svr",       // K3 Software-Server
	Port{26262, 17}:  "k3software-svr",       // K3 Software-Server
	Port{26263, 6}:   "k3software-cli",       // K3 Software-Client
	Port{26263, 17}:  "k3software-cli",       // K3 Software-Client
	Port{26486, 6}:   "exoline-tcp",          // exoline-udp | EXOline-UDP
	Port{26486, 17}:  "exoline-udp",          // EXOline-UDP
	Port{26487, 6}:   "exoconfig",            // Missing description for exoconfig
	Port{26487, 17}:  "exoconfig",            // EXOconfig
	Port{26489, 6}:   "exonet",               // Missing description for exonet
	Port{26489, 17}:  "exonet",               // EXOnet
	Port{26900, 17}:  "hexen2",               // Hexen 2 game server
	Port{27000, 6}:   "flexlm0",              // FlexLM license manager additional ports
	Port{27000, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27001, 6}:   "flexlm1",              // FlexLM license manager additional ports
	Port{27001, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27002, 6}:   "flexlm2",              // FlexLM license manager additional ports
	Port{27002, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27003, 6}:   "flexlm3",              // FlexLM license manager additional ports
	Port{27003, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27004, 6}:   "flexlm4",              // FlexLM license manager additional ports
	Port{27004, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27005, 6}:   "flexlm5",              // FlexLM license manager additional ports
	Port{27005, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27006, 6}:   "flexlm6",              // FlexLM license manager additional ports
	Port{27006, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27007, 6}:   "flexlm7",              // FlexLM license manager additional ports
	Port{27007, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27008, 6}:   "flexlm8",              // FlexLM license manager additional ports
	Port{27008, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27009, 6}:   "flexlm9",              // FlexLM license manager additional ports
	Port{27009, 17}:  "flex-lm",              // FLEX LM (1-10)
	Port{27010, 6}:   "flexlm10",             // FlexLM license manager additional ports
	Port{27015, 17}:  "halflife",             // Half-life game server
	Port{27017, 6}:   "mongod",               // http:  docs.mongodb.org manual reference default-mongodb-port
	Port{27018, 6}:   "mongod",               // http:  docs.mongodb.org manual reference default-mongodb-port
	Port{27019, 6}:   "mongod",               // http:  docs.mongodb.org manual reference default-mongodb-port
	Port{27345, 6}:   "imagepump",            // Missing description for imagepump
	Port{27345, 17}:  "imagepump",            // ImagePump
	Port{27374, 6}:   "subseven",             // Subseven Windows trojan
	Port{27442, 6}:   "jesmsjc",              // Job controller service
	Port{27442, 17}:  "jesmsjc",              // Job controller service
	Port{27444, 17}:  "Trinoo_Bcast",         // Trinoo distributed attack tool Master
	Port{27500, 17}:  "quakeworld",           // Quake world
	Port{27504, 6}:   "kopek-httphead",       // Kopek HTTP Head Port
	Port{27504, 17}:  "kopek-httphead",       // Kopek HTTP Head Port
	Port{27665, 6}:   "Trinoo_Master",        // Trinoo distributed attack tool Master server control port
	Port{27782, 6}:   "ars-vista",            // ARS VISTA Application
	Port{27782, 17}:  "ars-vista",            // ARS VISTA Application
	Port{27876, 6}:   "astrolink",            // Astrolink Protocol
	Port{27910, 17}:  "quake2",               // Quake 2 game server
	Port{27960, 17}:  "quake3",               // Quake 3 Arena Server
	Port{27999, 6}:   "tw-auth-key",          // TW Authentication Key Distribution and | Attribute Certificate Services
	Port{27999, 17}:  "tw-auth-key",          // Attribute Certificate Services
	Port{28000, 6}:   "nxlmd",                // NX License Manager
	Port{28000, 17}:  "nxlmd",                // NX License Manager
	Port{28001, 6}:   "pqsp",                 // PQ Service
	Port{28017, 6}:   "mongod",               // http:  docs.mongodb.org manual reference default-mongodb-port
	Port{28119, 6}:   "a27-ran-ran",          // A27 cdma2000 RAN Management
	Port{28200, 6}:   "voxelstorm",           // VoxelStorm game server
	Port{28240, 6}:   "siemensgsm",           // Siemens GSM
	Port{28240, 17}:  "siemensgsm",           // Siemens GSM
	Port{28589, 6}:   "bosswave",             // Building operating system services wide area verified exchange
	Port{28910, 17}:  "heretic2",             // Heretic 2 game server
	Port{29118, 132}: "sgsap",                // SGsAP in 3GPP
	Port{29167, 6}:   "otmp",                 // ObTools Message Protocol
	Port{29167, 17}:  "otmp",                 // ObTools Message Protocol
	Port{29168, 132}: "sbcap",                // SBcAP in 3GPP
	Port{29169, 132}: "iuhsctpassoc",         // HNBAP and RUA Common Association
	Port{29999, 6}:   "bingbang",             // data exchange protocol for IEC61850 in wind power plants
	Port{30000, 6}:   "ndmps",                // Secure Network Data Management Protocol
	Port{30001, 6}:   "pago-services1",       // Pago Services 1
	Port{30001, 17}:  "pago-services1",       // Pago Services 1
	Port{30002, 6}:   "pago-services2",       // Pago Services 2
	Port{30002, 17}:  "pago-services2",       // Pago Services 2
	Port{30003, 6}:   "amicon-fpsu-ra",       // Amicon FPSU-IP Remote Administration
	Port{30004, 6}:   "amicon-fpsu-s",        // Amicon FPSU-IP VPN
	Port{30100, 6}:   "rwp",                  // Remote Window Protocol
	Port{30260, 6}:   "kingdomsonline",       // Kingdoms Online (CraigAvenue)
	Port{30260, 17}:  "kingdomsonline",       // Kingdoms Online (CraigAvenue)
	Port{30400, 6}:   "gs-realtime",          // GroundStar RealTime System
	Port{30832, 6}:   "samsung-disc",         // Samsung Convergence Discovery Protocol
	Port{30999, 6}:   "ovobs",                // OpenView Service Desk Client
	Port{30999, 17}:  "ovobs",                // OpenView Service Desk Client
	Port{31016, 6}:   "ka-sddp",              // ka-kdp | Kollective Agent Secure Distributed Delivery Protocol | Kollective Agent Kollective Delivery Protocol
	Port{31020, 6}:   "autotrac-acp",         // Autotrac ACP 245
	Port{31029, 6}:   "yawn",                 // YaWN - Yet Another Windows Notifier
	Port{31029, 17}:  "yawn",                 // YaWN - Yet Another Windows Notifier
	Port{31335, 17}:  "Trinoo_Register",      // Trinoo distributed attack tool Bcast Daemon registration port
	Port{31337, 6}:   "Elite",                // Sometimes interesting stuff can be found here
	Port{31337, 17}:  "BackOrifice",          // cDc Back Orifice remote admin tool
	Port{31400, 6}:   "pace-licensed",        // PACE license server
	Port{31416, 6}:   "boinc",                // xqosd | BOINC Client Control | XQoS network monitor
	Port{31416, 17}:  "xqosd",                // XQoS network monitor
	Port{31457, 6}:   "tetrinet",             // TetriNET Protocol
	Port{31457, 17}:  "tetrinet",             // TetriNET Protocol
	Port{31620, 6}:   "lm-mon",               // lm mon
	Port{31620, 17}:  "lm-mon",               // lm mon
	Port{31685, 6}:   "dsx_monitor",          // dsx-monitor | DS Expert Monitor
	Port{31727, 6}:   "diagd",                // Missing description for diagd
	Port{31765, 6}:   "gamesmith-port",       // GameSmith Port
	Port{31765, 17}:  "gamesmith-port",       // GameSmith Port
	Port{31948, 6}:   "iceedcp_tx",           // iceedcp-tx | Embedded Device Configuration Protocol TX
	Port{31948, 17}:  "iceedcp_tx",           // Embedded Device Configuration Protocol TX
	Port{31949, 6}:   "iceedcp_rx",           // iceedcp-rx | Embedded Device Configuration Protocol RX
	Port{31949, 17}:  "iceedcp_rx",           // Embedded Device Configuration Protocol RX
	Port{32034, 6}:   "iracinghelper",        // iRacing helper service
	Port{32034, 17}:  "iracinghelper",        // iRacing helper service
	Port{32249, 6}:   "t1distproc60",         // T1 Distributed Processor
	Port{32249, 17}:  "t1distproc60",         // T1 Distributed Processor
	Port{32400, 6}:   "plex",                 // Plex multimedia
	Port{32483, 6}:   "apm-link",             // Access Point Manager Link
	Port{32483, 17}:  "apm-link",             // Access Point Manager Link
	Port{32635, 6}:   "sec-ntb-clnt",         // SecureNotebook-CLNT
	Port{32635, 17}:  "sec-ntb-clnt",         // SecureNotebook-CLNT
	Port{32636, 6}:   "DMExpress",            // Missing description for DMExpress
	Port{32636, 17}:  "DMExpress",            // Missing description for DMExpress
	Port{32767, 6}:   "filenet-powsrm",       // FileNet BPM WS-ReliableMessaging Client
	Port{32767, 17}:  "filenet-powsrm",       // FileNet BPM WS-ReliableMessaging Client
	Port{32768, 6}:   "filenet-tms",          // Filenet TMS
	Port{32768, 17}:  "omad",                 // OpenMosix Autodiscovery Daemon
	Port{32769, 6}:   "filenet-rpc",          // Filenet RPC
	Port{32769, 17}:  "filenet-rpc",          // Filenet RPC
	Port{32770, 6}:   "sometimes-rpc3",       // filenet-nch | Sometimes an RPC port on my Solaris box | Filenet NCH
	Port{32770, 17}:  "sometimes-rpc4",       // Sometimes an RPC port on my Solaris box
	Port{32771, 6}:   "sometimes-rpc5",       // filenet-rmi | Sometimes an RPC port on my Solaris box (rusersd) | FileNET RMI | FileNet RMI
	Port{32771, 17}:  "sometimes-rpc6",       // Sometimes an RPC port on my Solaris box (rusersd)
	Port{32772, 6}:   "sometimes-rpc7",       // filenet-pa | Sometimes an RPC port on my Solaris box (status) | FileNET Process Analyzer
	Port{32772, 17}:  "sometimes-rpc8",       // Sometimes an RPC port on my Solaris box (status)
	Port{32773, 6}:   "sometimes-rpc9",       // filenet-cm | Sometimes an RPC port on my Solaris box (rquotad) | FileNET Component Manager
	Port{32773, 17}:  "sometimes-rpc10",      // Sometimes an RPC port on my Solaris box (rquotad)
	Port{32774, 6}:   "sometimes-rpc11",      // filenet-re | Sometimes an RPC port on my Solaris box (rusersd) | FileNET Rules Engine
	Port{32774, 17}:  "sometimes-rpc12",      // Sometimes an RPC port on my Solaris box (rusersd)
	Port{32775, 6}:   "sometimes-rpc13",      // filenet-pch | Sometimes an RPC port on my Solaris box (status) | Performance Clearinghouse
	Port{32775, 17}:  "sometimes-rpc14",      // Sometimes an RPC port on my Solaris box (status)
	Port{32776, 6}:   "sometimes-rpc15",      // filenet-peior | Sometimes an RPC port on my Solaris box (sprayd) | FileNET BPM IOR
	Port{32776, 17}:  "sometimes-rpc16",      // Sometimes an RPC port on my Solaris box (sprayd)
	Port{32777, 6}:   "sometimes-rpc17",      // filenet-obrok | Sometimes an RPC port on my Solaris box (walld) | FileNet BPM CORBA
	Port{32777, 17}:  "sometimes-rpc18",      // Sometimes an RPC port on my Solaris box (walld)
	Port{32778, 6}:   "sometimes-rpc19",      // Sometimes an RPC port on my Solaris box (rstatd)
	Port{32778, 17}:  "sometimes-rpc20",      // Sometimes an RPC port on my Solaris box (rstatd)
	Port{32779, 6}:   "sometimes-rpc21",      // Sometimes an RPC port on my Solaris box
	Port{32779, 17}:  "sometimes-rpc22",      // Sometimes an RPC port on my Solaris box
	Port{32780, 6}:   "sometimes-rpc23",      // Sometimes an RPC port on my Solaris box
	Port{32780, 17}:  "sometimes-rpc24",      // Sometimes an RPC port on my Solaris box
	Port{32786, 6}:   "sometimes-rpc25",      // Sometimes an RPC port (mountd)
	Port{32786, 17}:  "sometimes-rpc26",      // Sometimes an RPC port
	Port{32787, 6}:   "sometimes-rpc27",      // Sometimes an RPC port dmispd (DMI Service Provider)
	Port{32787, 17}:  "sometimes-rpc28",      // Sometimes an RPC port
	Port{32801, 6}:   "mlsn",                 // Multiple Listing Service Network
	Port{32801, 17}:  "mlsn",                 // Multiple Listing Service Network
	Port{32811, 6}:   "retp",                 // Real Estate Transport Protocol
	Port{32896, 6}:   "idmgratm",             // Attachmate ID Manager
	Port{32896, 17}:  "idmgratm",             // Attachmate ID Manager
	Port{33060, 6}:   "mysqlx",               // MySQL Database Extended Interface
	Port{33123, 6}:   "aurora-balaena",       // Aurora (Balaena Ltd)
	Port{33123, 17}:  "aurora-balaena",       // Aurora (Balaena Ltd)
	Port{33331, 6}:   "diamondport",          // DiamondCentral Interface
	Port{33331, 17}:  "diamondport",          // DiamondCentral Interface
	Port{33333, 6}:   "dgi-serv",             // Digital Gaslight Service
	Port{33334, 6}:   "speedtrace",           // speedtrace-disc | SpeedTrace TraceAgent | SpeedTrace TraceAgent Discovery
	Port{33434, 6}:   "traceroute",           // traceroute use
	Port{33434, 17}:  "traceroute",           // traceroute use
	Port{33435, 6}:   "mtrace",               // IP Multicast Traceroute
	Port{33656, 6}:   "snip-slave",           // SNIP Slave
	Port{33656, 17}:  "snip-slave",           // SNIP Slave
	Port{34249, 6}:   "turbonote-2",          // TurboNote Relay Server Default Port
	Port{34249, 17}:  "turbonote-2",          // TurboNote Relay Server Default Port
	Port{34378, 6}:   "p-net-local",          // P-Net on IP local
	Port{34378, 17}:  "p-net-local",          // P-Net on IP local
	Port{34379, 6}:   "p-net-remote",         // P-Net on IP remote
	Port{34379, 17}:  "p-net-remote",         // P-Net on IP remote
	Port{34567, 6}:   "dhanalakshmi",         // edi_service | dhanalakshmi.org EDI Service
	Port{34962, 6}:   "profinet-rt",          // PROFInet RT Unicast
	Port{34962, 17}:  "profinet-rt",          // PROFInet RT Unicast
	Port{34963, 6}:   "profinet-rtm",         // PROFInet RT Multicast
	Port{34963, 17}:  "profinet-rtm",         // PROFInet RT Multicast
	Port{34964, 6}:   "profinet-cm",          // PROFInet Context Manager
	Port{34964, 17}:  "profinet-cm",          // PROFInet Context Manager
	Port{34980, 6}:   "ethercat",             // EtherCAT Port | EhterCAT Port
	Port{34980, 17}:  "ethercat",             // EhterCAT Port
	Port{35000, 6}:   "heathview",            // Missing description for heathview
	Port{35001, 6}:   "rt-viewer",            // ReadyTech Viewer
	Port{35002, 6}:   "rt-sound",             // ReadyTech Sound Server
	Port{35003, 6}:   "rt-devicemapper",      // ReadyTech DeviceMapper Server
	Port{35004, 6}:   "rt-classmanager",      // ReadyTech ClassManager
	Port{35005, 6}:   "rt-labtracker",        // ReadyTech LabTracker
	Port{35006, 6}:   "rt-helper",            // ReadyTech Helper Service
	Port{35100, 6}:   "axio-disc",            // Axiomatic discovery protocol
	Port{35354, 6}:   "kitim",                // KIT Messenger
	Port{35355, 6}:   "altova-lm",            // altova-lm-disc | Altova License Management | Altova License Management Discovery
	Port{35356, 6}:   "guttersnex",           // Gutters Note Exchange
	Port{35357, 6}:   "openstack-id",         // OpenStack ID Service
	Port{36001, 6}:   "allpeers",             // AllPeers Network
	Port{36001, 17}:  "allpeers",             // AllPeers Network
	Port{36411, 6}:   "wlcp",                 // Wireless LAN Control plane Protocol (WLCP)
	Port{36412, 132}: "s1-control",           // S1-Control Plane (3GPP)
	Port{36422, 132}: "x2-control",           // X2-Control Plane (3GPP)
	Port{36423, 6}:   "slmap",                // SLm Interface Application Protocol
	Port{36424, 6}:   "nq-ap",                // Nq and Nq' Application Protocol
	Port{36443, 6}:   "m2ap",                 // M2 Application Part
	Port{36444, 6}:   "m3ap",                 // M3 Application Part
	Port{36462, 6}:   "xw-control",           // Xw-Control Plane (3GPP)
	Port{36524, 6}:   "febooti-aw",           // Febooti Automation Workshop
	Port{36602, 6}:   "observium-agent",      // Observium statistics collection agent
	Port{36700, 6}:   "mapx",                 // MapX communication
	Port{36865, 6}:   "kastenxpipe",          // KastenX Pipe
	Port{36865, 17}:  "kastenxpipe",          // KastenX Pipe
	Port{37475, 6}:   "neckar",               // science + computing's Venus Administration Port
	Port{37475, 17}:  "neckar",               // science + computing's Venus Administration Port
	Port{37483, 6}:   "gdrive-sync",          // Google Drive Sync
	Port{37601, 6}:   "eftp",                 // Epipole File Transfer Protocol
	Port{37654, 6}:   "unisys-eportal",       // Unisys ClearPath ePortal
	Port{37654, 17}:  "unisys-eportal",       // Unisys ClearPath ePortal
	Port{38000, 6}:   "ivs-database",         // InfoVista Server Database
	Port{38001, 6}:   "ivs-insertion",        // InfoVista Server Insertion
	Port{38002, 6}:   "cresco-control",       // crescoctrl-disc | Cresco Controller | Cresco Controller Discovery
	Port{38037, 6}:   "landesk-cba",          // Missing description for landesk-cba
	Port{38037, 17}:  "landesk-cba",          // Missing description for landesk-cba
	Port{38201, 6}:   "galaxy7-data",         // Galaxy7 Data Tunnel
	Port{38201, 17}:  "galaxy7-data",         // Galaxy7 Data Tunnel
	Port{38202, 6}:   "fairview",             // Fairview Message Service
	Port{38202, 17}:  "fairview",             // Fairview Message Service
	Port{38203, 6}:   "agpolicy",             // AppGate Policy Server
	Port{38203, 17}:  "agpolicy",             // AppGate Policy Server
	Port{38292, 6}:   "landesk-cba",          // Missing description for landesk-cba
	Port{38293, 17}:  "landesk-cba",          // Missing description for landesk-cba
	Port{38412, 6}:   "ng-control",           // NG Control Plane (3GPP)
	Port{38422, 6}:   "xn-control",           // Xn Control Plane (3GPP)
	Port{38472, 6}:   "f1-control",           // F1 Control Plane (3GPP)
	Port{38800, 6}:   "sruth",                // Sruth is a service for the distribution of routinely- generated but arbitrary files based on a publish subscribe distribution model and implemented using a peer-to-peer transport mechanism
	Port{38865, 6}:   "secrmmsafecopya",      // Security approval process for use of the secRMM SafeCopy program
	Port{39213, 17}:  "sygatefw",             // Sygate Firewall management port version 3.0 build 521 and above
	Port{39681, 6}:   "turbonote-1",          // TurboNote Default Port
	Port{39681, 17}:  "turbonote-1",          // TurboNote Default Port
	Port{40000, 6}:   "safetynetp",           // SafetyNET p
	Port{40000, 17}:  "safetynetp",           // SafetyNET p
	Port{40023, 6}:   "k-patentssensor",      // K-PatentsSensorInformation
	Port{40404, 6}:   "sptx",                 // Simplify Printing TX
	Port{40841, 6}:   "cscp",                 // Missing description for cscp
	Port{40841, 17}:  "cscp",                 // CSCP
	Port{40842, 6}:   "csccredir",            // Missing description for csccredir
	Port{40842, 17}:  "csccredir",            // CSCCREDIR
	Port{40843, 6}:   "csccfirewall",         // Missing description for csccfirewall
	Port{40843, 17}:  "csccfirewall",         // CSCCFIREWALL
	Port{40853, 6}:   "ortec-disc",           // ORTEC Service Discovery
	Port{40853, 17}:  "ortec-disc",           // ORTEC Service Discovery
	Port{41111, 6}:   "fs-qos",               // Foursticks QoS Protocol
	Port{41111, 17}:  "fs-qos",               // Foursticks QoS Protocol
	Port{41121, 6}:   "tentacle",             // Tentacle Server
	Port{41230, 6}:   "z-wave-s",             // Z-Wave Protocol over SSL TLS | Z-Wave Protocol over DTLS
	Port{41794, 6}:   "crestron-cip",         // Crestron Control Port
	Port{41794, 17}:  "crestron-cip",         // Crestron Control Port
	Port{41795, 6}:   "crestron-ctp",         // Crestron Terminal Port
	Port{41795, 17}:  "crestron-ctp",         // Crestron Terminal Port
	Port{41796, 6}:   "crestron-cips",        // Crestron Secure Control Port
	Port{41797, 6}:   "crestron-ctps",        // Crestron Secure Terminal Port
	Port{42508, 6}:   "candp",                // Computer Associates network discovery protocol
	Port{42508, 17}:  "candp",                // Computer Associates network discovery protocol
	Port{42509, 6}:   "candrp",               // CA discovery response
	Port{42509, 17}:  "candrp",               // CA discovery response
	Port{42510, 6}:   "caerpc",               // CA eTrust RPC
	Port{42510, 17}:  "caerpc",               // CA eTrust RPC
	Port{43000, 6}:   "recvr-rc",             // recvr-rc-disc | Receiver Remote Control | Receiver Remote Control Discovery
	Port{43188, 6}:   "reachout",             // Missing description for reachout
	Port{43188, 17}:  "reachout",             // REACHOUT
	Port{43189, 6}:   "ndm-agent-port",       // Missing description for ndm-agent-port
	Port{43189, 17}:  "ndm-agent-port",       // NDM-AGENT-PORT
	Port{43190, 6}:   "ip-provision",         // Missing description for ip-provision
	Port{43190, 17}:  "ip-provision",         // IP-PROVISION
	Port{43191, 6}:   "noit-transport",       // Reconnoiter Agent Data Transport
	Port{43210, 6}:   "shaperai",             // shaperai-disc | Shaper Automation Server Management | Shaper Automation Server Management Discovery
	Port{43439, 6}:   "eq3-update",           // eq3-config | EQ3 firmware update | EQ3 discovery and configuration
	Port{43440, 6}:   "ew-mgmt",              // ew-disc-cmd | Cisco EnergyWise Management | Cisco EnergyWise Discovery and Command Flooding
	Port{43440, 17}:  "ew-disc-cmd",          // Cisco EnergyWise Discovery and Command Flooding
	Port{43441, 6}:   "ciscocsdb",            // Cisco NetMgmt DB Ports
	Port{43441, 17}:  "ciscocsdb",            // Cisco NetMgmt DB Ports
	Port{44123, 6}:   "z-wave-tunnel",        // Z-Wave Secure Tunnel
	Port{44321, 6}:   "pmcd",                 // PCP server (pmcd)
	Port{44321, 17}:  "pmcd",                 // PCP server (pmcd)
	Port{44322, 6}:   "pmcdproxy",            // PCP server (pmcd) proxy
	Port{44322, 17}:  "pmcdproxy",            // PCP server (pmcd) proxy
	Port{44323, 6}:   "pmwebapi",             // HTTP binding for Performance Co-Pilot client API
	Port{44334, 6}:   "tinyfw",               // tiny personal firewall admin port
	Port{44442, 6}:   "coldfusion-auth",      // ColdFusion Advanced Security Siteminder Authentication Port (by Allaire Netegrity)
	Port{44443, 6}:   "coldfusion-auth",      // ColdFusion Advanced Security Siteminder Authentication Port (by Allaire Netegrity)
	Port{44444, 6}:   "cognex-dataman",       // Cognex DataMan Management Protocol
	Port{44544, 6}:   "domiq",                // DOMIQ Building Automation
	Port{44553, 6}:   "rbr-debug",            // REALbasic Remote Debug
	Port{44553, 17}:  "rbr-debug",            // REALbasic Remote Debug
	Port{44600, 6}:   "asihpi",               // AudioScience HPI
	Port{44818, 6}:   "EtherNetIP-2",         // EtherNet IP-2 | EtherNet-IP-2 | EtherNet IP messaging
	Port{44818, 17}:  "EtherNetIP-2",         // EtherNet IP messaging
	Port{44900, 6}:   "m3da",                 // m3da-disc | M3DA is used for efficient machine-to-machine communications | M3DA Discovery is used for efficient machine-to-machine communications
	Port{45000, 6}:   "asmp",                 // asmp-mon | Nuance AutoStore Status Monitoring Protocol (data transfer) | Nuance AutoStore Status Monitoring Protocol (device monitoring)
	Port{45000, 17}:  "ciscopop",             // Cisco Postoffice Protocol for Cisco Secure IDS
	Port{45001, 6}:   "asmps",                // Nuance AutoStore Status Monitoring Protocol (secure data transfer)
	Port{45002, 6}:   "rs-status",            // Redspeed Status Monitor
	Port{45045, 6}:   "synctest",             // Remote application control protocol
	Port{45054, 6}:   "invision-ag",          // InVision AG
	Port{45054, 17}:  "invision-ag",          // InVision AG
	Port{45514, 6}:   "cloudcheck",           // cloudcheck-ping | ASSIA CloudCheck WiFi Management System | ASSIA CloudCheck WiFi Management keepalive
	Port{45678, 6}:   "eba",                  // EBA PRISE
	Port{45678, 17}:  "eba",                  // EBA PRISE
	Port{45824, 6}:   "dai-shell",            // Server for the DAI family of client-server products
	Port{45825, 6}:   "qdb2service",          // Qpuncture Data Access Service
	Port{45825, 17}:  "qdb2service",          // Qpuncture Data Access Service
	Port{45966, 6}:   "ssr-servermgr",        // SSRServerMgr
	Port{45966, 17}:  "ssr-servermgr",        // SSRServerMgr
	Port{46336, 6}:   "inedo",                // Listen port used for Inedo agent communication
	Port{46998, 6}:   "spremotetablet",       // Connection between a desktop computer or server and a signature tablet to capture handwritten signatures
	Port{46999, 6}:   "mediabox",             // MediaBox Server
	Port{46999, 17}:  "mediabox",             // MediaBox Server
	Port{47000, 6}:   "mbus",                 // Message Bus
	Port{47000, 17}:  "mbus",                 // Message Bus
	Port{47001, 6}:   "winrm",                // Windows Remote Management Service
	Port{47100, 6}:   "jvl-mactalk",          // Configuration of motors connected to Industrial Ethernet
	Port{47557, 6}:   "dbbrowse",             // Databeam Corporation
	Port{47557, 17}:  "dbbrowse",             // Databeam Corporation
	Port{47624, 6}:   "directplaysrvr",       // Direct Play Server
	Port{47624, 17}:  "directplaysrvr",       // Direct Play Server
	Port{47806, 6}:   "ap",                   // ALC Protocol
	Port{47806, 17}:  "ap",                   // ALC Protocol
	Port{47808, 6}:   "bacnet",               // Building Automation and Control Networks
	Port{47808, 17}:  "bacnet",               // Building Automation and Control Networks
	Port{47809, 6}:   "presonus-ucnet",       // PreSonus Universal Control Network Protocol
	Port{48000, 6}:   "nimcontroller",        // Nimbus Controller
	Port{48000, 17}:  "nimcontroller",        // Nimbus Controller
	Port{48001, 6}:   "nimspooler",           // Nimbus Spooler
	Port{48001, 17}:  "nimspooler",           // Nimbus Spooler
	Port{48002, 6}:   "nimhub",               // Nimbus Hub
	Port{48002, 17}:  "nimhub",               // Nimbus Hub
	Port{48003, 6}:   "nimgtw",               // Nimbus Gateway
	Port{48003, 17}:  "nimgtw",               // Nimbus Gateway
	Port{48004, 6}:   "nimbusdb",             // NimbusDB Connector
	Port{48005, 6}:   "nimbusdbctrl",         // NimbusDB Control
	Port{48049, 6}:   "3gpp-cbsp",            // 3GPP Cell Broadcast Service Protocol
	Port{48050, 6}:   "weandsf",              // WeFi Access Network Discovery and Selection Function
	Port{48128, 6}:   "isnetserv",            // Image Systems Network Services
	Port{48128, 17}:  "isnetserv",            // Image Systems Network Services
	Port{48129, 6}:   "blp5",                 // Bloomberg locator
	Port{48129, 17}:  "blp5",                 // Bloomberg locator
	Port{48556, 6}:   "com-bardac-dw",        // Missing description for com-bardac-dw
	Port{48556, 17}:  "com-bardac-dw",        // Missing description for com-bardac-dw
	Port{48619, 6}:   "iqobject",             // Missing description for iqobject
	Port{48619, 17}:  "iqobject",             // Missing description for iqobject
	Port{48653, 6}:   "robotraconteur",       // Robot Raconteur transport
	Port{49000, 6}:   "matahari",             // Matahari Broker
	Port{49001, 6}:   "nusrp",                // nusdp-disc | Nuance Unity Service Request Protocol | Nuance Unity Service Discovery Protocol
	Port{49150, 6}:   "inspider",             // InSpider System
	Port{49400, 6}:   "compaqdiag",           // Compaq Web-based management
	Port{50000, 6}:   "ibm-db2",              // (also Internet Intranet Input Method Server Framework?)
	Port{50002, 6}:   "iiimsf",               // Internet Intranet Input Method Server Framework
	Port{54320, 6}:   "bo2k",                 // Back Orifice 2K Default Port
	Port{54321, 17}:  "bo2k",                 // Back Orifice 2K Default Port
	Port{61439, 6}:   "netprowler-manager",   // Missing description for netprowler-manager
	Port{61440, 6}:   "netprowler-manager2",  // Missing description for netprowler-manager2
	Port{61441, 6}:   "netprowler-sensor",    // Missing description for netprowler-sensor
	Port{62078, 6}:   "iphone-sync",          // Apparently used by iPhone while syncing - http:  code.google.com p iphone-elite source browse wiki Port_62078.wiki
	Port{64738, 17}:  "murmur",               // Murmur is the server-side software for Mumble open source voice chat software
	Port{65301, 6}:   "pcanywhere",           // Missing description for pcanywhere
	Port{0, 0}:       "HOPOPT",               // Adding in some default IP protocol names below.
	Port{0, 1}:       "ICMP",
	Port{0, 2}:       "IGMP",
	Port{0, 4}:       "IPc4",
	Port{0, 8}:       "EGP",
	Port{0, 9}:       "IGP",
	Port{0, 20}:      "HMP",
	Port{0, 31}:      "MFE-NSP",
	Port{0, 40}:      "IL",
	Port{0, 41}:      "IPv6",
	Port{0, 46}:      "RSVP",
	Port{0, 47}:      "GRE",
	Port{0, 50}:      "ESP",
	Port{0, 51}:      "AH",
	Port{0, 58}:      "IPv6-ICMP",
	Port{0, 97}:      "ETHERIP",
	Port{0, 103}:     "PIM",
	Port{0, 104}:     "ARIS",
	Port{0, 105}:     "SCPS",
	Port{0, 112}:     "VRRP",
	Port{0, 115}:     "L2TP",
	Port{0, 132}:     "SCTP",
	Port{0, 240}:     "240",
}
