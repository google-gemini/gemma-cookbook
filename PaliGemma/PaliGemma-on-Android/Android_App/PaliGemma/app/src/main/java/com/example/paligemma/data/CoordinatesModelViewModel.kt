package com.example.paligemma.data

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.launch

class CoordinatesModelViewModel(
    private val coordinatesModelRepo: CoordinatesModelRepo
) : ViewModel() {

    var uiState by mutableStateOf<UiState>(UiState.Idle)

    fun getCoordinatesModel(requestModel: RequestModel) {
        uiState = UiState.Loading
        viewModelScope.launch {
            try {
                val coordinatesModel = coordinatesModelRepo
                    .getCoordinatesModel(requestModel)
                    .body()

                if (coordinatesModel?.result != null) {
                    uiState = UiState.Success(
                        coordinatesModel.result
                    )
                } else if (coordinatesModel?.error != null) {
                    uiState = UiState.Error(
                        coordinatesModel.error
                    )
                } else if (coordinatesModel?.response != null) {
                    uiState = UiState.CaptionResponse(coordinatesModel.response)
                } else {
                    uiState = UiState.Error("No result found.")
                }
            } catch (e: Exception) {
                e.printStackTrace()
                uiState = if (e.message != null) {
                    UiState.Error(e.message!!)
                } else {
                    UiState.Error("An unknown error occurred.")
                }
            }
        }
    }

    fun resetData() {
        uiState = UiState.Idle
    }
}

class CoordinatesModelViewModelFactory(private val coordinatesModelRepo: CoordinatesModelRepo) :
    ViewModelProvider.NewInstanceFactory() {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return CoordinatesModelViewModel(coordinatesModelRepo) as T
    }
}

sealed interface UiState {
    data object Idle : UiState
    data object Loading : UiState
    data class Success(val result: List<Result>) : UiState
    data class CaptionResponse(val msg: String) : UiState
    data class Error(val e: String) : UiState
}