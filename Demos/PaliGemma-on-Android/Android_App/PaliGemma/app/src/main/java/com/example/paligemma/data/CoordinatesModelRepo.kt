package com.example.paligemma.data

import retrofit2.Response

interface CoordinatesModelRepo {
    suspend fun getCoordinatesModel(requestModel: RequestModel): Response<CoordinatesModel>
}