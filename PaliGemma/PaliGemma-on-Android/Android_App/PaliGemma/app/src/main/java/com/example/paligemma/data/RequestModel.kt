package com.example.paligemma.data

import android.net.Uri
import java.io.File

data class RequestModel(
    val text: String,
    val width: String,
    val height: String,
    val uri: Uri,
)
