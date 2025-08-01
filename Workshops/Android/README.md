# Getting Started with Gemma on Android 🚀

This guide provides the resources and instructions to integrate the Gemma model into your Android applications. We offer some simple boilerplate projects to get you started quickly.

The `GemmaAppText` directory contains a minimal-implementation example that uses text-only input.

-----

## Option 1: Using the Boilerplate (Quickstart)

The easiest way to start is by using the `GemmaAppText` project. It's already configured with all the necessary dependencies, UI components, and model integration logic.

1.  **Open** the `GemmaAppText` project in Android Studio.
2.  **Download** the `Gemma 3n` model from [Kaggle](https://www.kaggle.com/models/google/gemma-3n/tfLite).
3.  **Push the model** to your test device using the Android Debug Bridge (adb):
    ```bash
    $ adb shell mkdir -p /data/local/tmp/llm/
    $ adb push gemma-3n-E2B-it-int4.task /data/local/tmp/llm/
    ```
4.  **Run** the app.

-----

## Option 2: Manual Setup from an "Empty Activity" Project

If you prefer to integrate Gemma into a new or existing project, follow these steps.

### 1. Add Dependencies

First, you need to add the MediaPipe LLM Inference API library to your project.

#### Add the Library to `app/build.gradle.kts`

In your app-level `build.gradle.kts` file, add the following dependency:

```kotlin
dependencies {
    implementation(libs.mediapipe.tasks.genai)
}
```

#### Add the Version to `gradle/libs.versions.toml`

To manage the dependency version, add the following entries to your `libs.versions.toml` file:

```toml
# In the [versions] block
mediapipeTasksGenai = "0.10.25"

# In the [libraries] block
mediapipe-tasks-genai = { group = "com.google.mediapipe", name = "tasks-genai", version.ref = "mediapipeTasksGenai" }
```

### 2. Download and Deploy the Model

Download the **Gemma 3n** model from [Kaggle](https://www.kaggle.com/models/google/gemma-3n/tfLite).

Next, push the downloaded `.task` file to your Android device.

```bash
$ adb shell mkdir -p /data/local/tmp/llm/
$ adb push gemma-3n-E2B-it-int4.task /data/local/tmp/llm/
```

> **Note**: For development, using `adb` to push the model to your device is fine. For a production app, you should host the model on a server and download it at runtime, as the model file is too large to be bundled directly within your APK.

### 3. Copy Sample Code

To get the UI and model logic working quickly, copy the following files from the `GemmaAppText` boilerplate into your project:

  * **UI and State Management:**
      * `ChatScreen.kt`
      * `ChatViewModel.kt`
      * `UiState.kt`
  * **Model Wrapper:**
      * `GemmaModel.kt` (This class wraps the model and sets default configuration.)

### 4. Update MainActivity

Finally, modify your `MainActivity.kt` to use the new `ChatScreen` Composable. Replace the default `setContent` block with the following:

```kotlin
    setContent {
        MyApplicationTheme { // Or your app's theme
            ChatScreen()
        }
    }
```

After completing these steps, you can run your app to see the Gemma-powered chat interface in action.

![image](GemmaAppText.gif)
