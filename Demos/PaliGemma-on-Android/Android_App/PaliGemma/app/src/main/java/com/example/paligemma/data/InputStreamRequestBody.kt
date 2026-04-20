package com.example.paligemma.data
import android.content.ContentResolver
import android.net.Uri
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody
import okio.BufferedSink
import okio.source
import java.io.IOException

class InputStreamRequestBody(
    private val contentResolver: ContentResolver,
    private val uri: Uri
) : RequestBody() {

    override fun contentType() = contentResolver.getType(uri)?.toMediaTypeOrNull()

    override fun contentLength(): Long = -1

    @Throws(IOException::class)
    override fun writeTo(sink: BufferedSink) {

        val input = contentResolver.openInputStream(uri)
        input?.use { sink.writeAll(it.source()) }
            ?: throw IOException("Could not open $uri")
    }
}