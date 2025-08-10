/**
 * Copyright (c) 2022, Xerra Earth Observation Institute.
 * Copyright (c) 2025, Simeon Miteff.
 *
 * See LICENSE.TXT in the root directory of this source tree.
 */

#include "mock_tcpinfo.h"

void zero( void* ptr ) {
    memset(ptr, 0, sizeof(struct tcp_info));
}

void set_snd_wscale(
    void* ptr,
    uint8_t snd_wscale
    ) {
    struct tcp_info *t = (struct tcp_info *)ptr;
    t->tcpi_snd_wscale = snd_wscale;
}

void set_rcv_wscale(
    void* ptr,
    uint8_t rcv_wscale
    ) {
    struct tcp_info *t = (struct tcp_info *)ptr;
    t->tcpi_rcv_wscale = rcv_wscale;
}

void set_delivery_rate_app_limited(
    void* ptr,
    bool delivery_rate_app_limited
    ) {
    struct tcp_info *t = (struct tcp_info *)ptr;
    t->tcpi_delivery_rate_app_limited = delivery_rate_app_limited;
}

void set_fastopen_client_fail(
    void* ptr,
    uint8_t fastopen_client_fail
    ) {
    struct tcp_info *t = (struct tcp_info *)ptr;
    t->tcpi_fastopen_client_fail = fastopen_client_fail;
}

//void set_fields( void* ptr,
//                 uint8_t snd_wscale,
//                 uint8_t rcv_wscale,
//                 bool *delivery_rate_app_limited,
//                 uint8_t *fastopen_client_fail
//                 ) {
//
//    struct tcp_info t;
//    memset(&t, 0, sizeof(struct tcp_info));
//
//    t.tcpi_snd_wscale = snd_wscale;
//    t.tcpi_rcv_wscale = rcv_wscale;
//    if delivery_rate_app_limited != NULL {}
//        t.tcpi_delivery_rate_app_limited = *delivery_rate_app_limited ? 1 : 0;
//    }
//    if fastopen_client_fail != NULL {
//        t.tcpi_fastopen_client_fail = *fastopen_client_fail;
//    }
//
//    memcpy(ptr, &t, sizeof(struct tcp_info));
//}