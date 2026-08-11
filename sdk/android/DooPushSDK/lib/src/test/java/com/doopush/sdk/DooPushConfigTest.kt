package com.doopush.sdk

import org.junit.Assert.assertEquals
import org.junit.Assert.assertThrows
import org.junit.Test

class DooPushConfigTest {

    @Test
    fun acceptsServerGeneratedAppKey() {
        val appKey = "dp_ak_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

        val config = DooPushConfig.create("1", appKey)

        assertEquals(appKey, config.appKey)
    }

    @Test
    fun rejectsOldLengthAppKey() {
        assertInvalidAppKey("dp_ak_LmsRDkJp395R4_ik0m-TWDAIyWqEjmbt")
    }

    @Test
    fun rejectsEmptyAppKeySuffix() {
        assertInvalidAppKey("dp_ak_")
    }

    @Test
    fun rejectsNonBase64UrlCharacters() {
        assertInvalidAppKey("dp_ak_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA+")
    }

    private fun assertInvalidAppKey(appKey: String) {
        assertThrows(DooPushConfigException::class.java) {
            DooPushConfig.create("1", appKey)
        }
    }
}
