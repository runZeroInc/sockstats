/**
 * Copyright (c) 2022, Xerra Earth Observation Institute.
 * Copyright (c) 2025, Simeon Miteff.
 *
 * See LICENSE.TXT in the root directory of this source tree.
 */

package linux

/*
#include "mock_tcpinfo.h"
*/
import "C"
import (
	"unsafe"
)

func (packed *RawTCPInfo) MockSetFields(
	SndWScale uint8,
	RcvWScale uint8,
	DeliveryRateAppLimited NullableBool,
	FastOpenClientFail NullableUint8,
) {
	C.zero(unsafe.Pointer(packed))

	C.set_snd_wscale(
		unsafe.Pointer(packed),
		C.uchar(SndWScale),
	)

	C.set_rcv_wscale(
		unsafe.Pointer(packed),
		C.uchar(RcvWScale),
	)

	if DeliveryRateAppLimited.Valid {
		C.set_delivery_rate_app_limited(
			unsafe.Pointer(packed),
			C.bool(DeliveryRateAppLimited.Value),
		)
	}

	if FastOpenClientFail.Valid {
		C.set_fastopen_client_fail(
			unsafe.Pointer(packed),
			C.uchar(FastOpenClientFail.Value),
		)
	}
}
