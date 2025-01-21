package com.example.paligemma.presentation

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableFloatStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.example.paligemma.data.CoordinatesModelRepo
import com.example.paligemma.data.RequestModel
import com.example.paligemma.data.Result
import kotlinx.coroutines.launch
import java.util.Locale

class CoordinatesModelViewModel(
    private val coordinatesModelRepo: CoordinatesModelRepo
) : ViewModel() {

    var uiState by mutableStateOf<UiState>(UiState.Idle)

    var imageLeftDistance by mutableFloatStateOf(0.0f)

    fun getCoordinatesModel(requestModel: RequestModel) {
        uiState = UiState.Loading
        viewModelScope.launch {
            try {
                val coordinatesModel = coordinatesModelRepo
                    .getCoordinatesModel(requestModel)
                    .body()

                uiState = when {
                    coordinatesModel?.result != null -> {
                        UiState.ObjectDetectionResponse(
                            getObjectDetectionUiData(coordinatesModel.result)
                        )
                    }

                    coordinatesModel?.response != null -> {
                        UiState.CaptionResponse(coordinatesModel.response)
                    }

                    coordinatesModel?.polygons != null && coordinatesModel.labels != null -> {
                        UiState.SegmentationResponse(
                            getSegmentationUiData(
                                coordinatesModel.polygons,
                                coordinatesModel.labels
                            )
                        )
                    }

                    coordinatesModel?.error != null -> {
                        UiState.Error(
                            coordinatesModel.error
                        )
                    }

                    else -> {
                        UiState.Error("No result found.")
                    }
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

    private fun getSegmentationUiData(
        polygons: List<List<List<Double>>>,
        labels: List<String>
    ): SegmentationUiData {
        val map = HashMap<String, Color>()
        val segmentUiList = mutableListOf<SegmentationUiData.SegmentUiData>()
        polygons.getOrNull(0)?.forEachIndexed { idx, it ->
            map.putIfAbsent(labels[idx].filterLabel(), getRandomColor())
            val points = it.chunked(2).map { (x, y) ->
                Offset(x.toFloat() + imageLeftDistance, y.toFloat())
            }

            val path = Path()
            path.moveTo(points.first().x, points.first().y)
            points.drop(1).forEach { point ->
                path.lineTo(point.x, point.y)
            }
            path.close()

            segmentUiList.add(
                SegmentationUiData.SegmentUiData(
                    path = path,
                    color = map.getOrDefault(labels[idx].filterLabel(), Color.Transparent)
                )
            )
        }
        return SegmentationUiData(
            data = segmentUiList,
            colorsMap = map
        )
    }

    private fun getObjectDetectionUiData(results: List<Result>): List<ObjectDetectionUiData> {
        val map = HashMap<String, Color>()
        return results.map { result ->
            val (y1, x1, y2, x2) = result.coordinates
            map.putIfAbsent(result.label.filterLabel(), getRandomColor())

            ObjectDetectionUiData(
                topLeft = Offset(x1.toFloat() + imageLeftDistance, y1.toFloat()),
                color = map.getOrDefault(result.label.filterLabel(), Color.Transparent),
                size = Size(
                    width = (x2 - x1).toFloat(),
                    height = (y2 - y1).toFloat()
                ),
                text = result.label.filterLabel(),
                textTopLeft = Offset(x1.toFloat() + imageLeftDistance, y1.toFloat() - 40)
            )
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

private fun String.filterLabel(): String {
    return this.replace("'", "").trim().replaceFirstChar { if (it.isLowerCase()) it.titlecase(Locale.US) else it.toString() }
}

private fun getRandomColor(): Color {
    val r = (1..254).random()
    val g = (1..254).random()
    val b = (1..254).random()
    return Color(
        red = r,
        green = g,
        blue = b
    )
}