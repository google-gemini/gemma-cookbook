package com.example.paligemma.data

import android.content.Context
import android.net.Uri
import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import retrofit2.Response
import java.io.File


class CoordinatesModelRepoImpl(
    val applicationContext: Context
) : CoordinatesModelRepo {

    override suspend fun getCoordinatesModel(requestModel: RequestModel): Response<CoordinatesModel> {
        return withContext(Dispatchers.IO) {
            val file = getTempFile(applicationContext, requestModel.uri)
            Log.d("TAG", "getCoordinatesModel: File = $file")
            CoordinatesModelApi.instance.getCoordinatesModel(
                text = requestModel.text.toRequestBody(MultipartBody.FORM),
                width = requestModel.width.toRequestBody(MultipartBody.FORM),
                height = requestModel.height.toRequestBody(MultipartBody.FORM),
                image = MultipartBody.Part.createFormData(
                    name = "image",
                    filename = file?.name,
                    body = file!!.asRequestBody(MultipartBody.FORM)
                )
            )
        }
    }

}

fun getTempFile(context: Context, uri: Uri): File? {
    try {
        val resolver = context.contentResolver
        val tempFile = File(context.cacheDir, "${System.currentTimeMillis()}_.jpg")
        val inputStream = resolver.openInputStream(uri) ?: return null
        val outputStream = tempFile.outputStream()
        val buffer = ByteArray(4 * 1024) // Adjust buffer size as needed
        while (true) {
            val read = inputStream.read(buffer)
            if (read == -1) break
            outputStream.write(buffer, 0, read)
        }
        outputStream.flush()
        return tempFile
    } catch (e: Exception) {
        e.printStackTrace()
        return null
    }
}