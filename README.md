
# Welcome to the Gemma Cookbook
This is a collection of guides and examples for [Google Gemma](https://ai.google.dev/gemma/). Gemma is a family of lightweight, state-of-the art open models built from the same research and technology used to create the Gemini models.

## Get started with the Gemma models
Gemma is a family of lightweight, state-of-the art open models built from the same research and technology used to create the Gemini models. The Gemma model family includes:
* base Gemma
  * [Gemma](https://ai.google.dev/gemma/docs/model_card)
  * [Gemma 2](https://ai.google.dev/gemma/docs/model_card_2)
* Gemma variants
  * [CodeGemma](https://ai.google.dev/gemma/docs/codegemma)
  * [PaliGemma](https://ai.google.dev/gemma/docs/paligemma)
  * [RecurrentGemma](https://ai.google.dev/gemma/docs/recurrentgemma)
  * [ShieldGemma](https://ai.google.dev/gemma/docs/shieldgemma)
  * [DataGemma](https://ai.google.dev/gemma/docs/datagemma)

You can find the Gemma models on GitHub, Hugging Face models, Kaggle, Google Cloud Vertex AI Model Garden, and [ai.nvidia.com](https://ai.nvidia.com).

## Partner quickstart guides
| Company      | Description                                                                                                                                                                                  |
| ------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Hugging Face | [Utilize Hugging Face Transformers and TRL for fine-tuning and inference tasks with Gemma models.](partner-quickstarts/gemma-huggingface.md)                                                 |
| NVIDIA       | [Fine-tune Gemma models with NVIDIA NeMo Framework and export to TensorRT-LLM for production.](partner-quickstarts/Gemma-NVidia/)                                                            |
| LangChain    | This [tutorial](partner-quickstarts/gemma-langchain.ipynb) shows you how to get started with Gemma and LangChain, running in Google Cloud or in your Colab environment.                      |
| MongoDB      | This [article](partner-quickstarts/rag_with_hugging_face_gemma_mongodb.ipynb) presents how to leverage Gemma as the foundation model in a retrieval-augmented generation pipeline or system. |

## Workshops and technical talks
| Notebook                                                                                                                    | Description                                                                                                                                        |
| --------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Workshop_How_to_Fine_tuning_Gemma.ipynb](Workshops/Workshop_How_to_Fine_tuning_Gemma.ipynb) | Recommended finetuning notebook for getting started                                                                                                |
| [Self_extend_Gemma.ipynb](Gemma/Self_extend_Gemma.ipynb)                                                                    | Self-extend context window for Gemma in the I/O 2024 [Keras talk](https://www.youtube.com/watch?v=TV7qCk1dBWA)                                     |
| [Gemma_control_vectors.ipynb](Gemma/Gemma_control_vectors.ipynb)                                                            | Implement [control vectors](https://arxiv.org/abs/2310.01405) with Gemma in the I/O 2024 [Keras talk](https://www.youtube.com/watch?v=TV7qCk1dBWA) |

## Accompanying notebooks for the Build with AI video series
| Folder                                                      |
| ----------------------------------------------------------- |
| [Business email assistant](Gemma/business-email-assistant/) |
| [Personal code assistant](Gemma/personal-code-assistant/)   |
| [Spoken language tasks](Gemma/spoken-language-tasks/)       |

## Cookbook table of contents

| **Gemma model overview**                         |                                                                      |
| :----------------------------------------------- | -------------------------------------------------------------------- |
| [Common_use_cases.ipynb](Common_use_cases.ipynb) | Illustrate some common use cases for Gemma, CodeGemma and PaliGemma. |

#### Gemma

| **Inference and serving**                                                                                            |                                                                                                                                                                                         |
| :------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Keras_Gemma_2_Quickstart.ipynb](Gemma/Keras_Gemma_2_Quickstart.ipynb)                                               | Gemma 2 pre-trained 9B model quickstart tutorial with Keras.                                                                                                                            |
| [Keras_Gemma_2_Quickstart_Chat.ipynb](Gemma/Keras_Gemma_2_Quickstart_Chat.ipynb)                                     | Gemma 2 instruction-tuned 9B model quickstart tutorial with Keras. Referenced in this [blog](https://developers.googleblog.com/en/fine-tuning-gemma-2-with-keras-hugging-face-update/). |
| [Chat_and_distributed_pirate_tuning.ipynb](Gemma/Chat_and_distributed_pirate_tuning.ipynb)                           | Chat with Gemma 7B and finetune it so that it generates responses in pirates' tone.                                                                                                     |
| [gemma_inference_on_tpu.ipynb](Gemma/gemma_inference_on_tpu.ipynb)                                                   | Basic inference of Gemma with JAX/Flax on TPU.                                                                                                                                          |
| [gemma_data_parallel_inference_in_jax_tpu.ipynb](Gemma/gemma_data_parallel_inference_in_jax_tpu.ipynb)               | Parallel inference of Gemma with JAX/Flax on TPU.                                                                                                                                       |
| [Gemma_Basics_with_HF.ipynb](Gemma/Gemma_Basics_with_HF.ipynb)                                                       | Load, run, finetune and deploy Gemma using [Hugging Face](https://huggingface.co/).                                                                                                     |
| [Gemma_with_Langfun_and_LlamaCpp.ipynb](Gemma/Gemma_with_Langfun_and_LlamaCpp.ipynb)                                 | Leverage [Langfun](https://github.com/google/langfun) to seamlessly integrate natural language with programming using Gemma 2 and [LlamaCpp](https://github.com/ggerganov/llama.cpp).   |
| [Gemma_with_Langfun_and_LlamaCpp_Python_Bindings.ipynb](Gemma/Gemma_with_Langfun_and_LlamaCpp_Python_Bindings.ipynb) | Leverage [Langfun](https://github.com/google/langfun) for smooth language-program interaction with Gemma 2 and [llama-cpp-python](https://github.com/abetlen/llama-cpp-python).         |
| [Guess_the_word.ipynb](Gemma/Guess_the_word.ipynb)                                                                   | Play a word guessing game with Gemma using Keras.                                                                                                                                       |
| [Game_Design_Brainstorming.ipynb](Gemma/Game_Design_Brainstorming.ipynb)                                             | Use Gemma to brainstorm ideas during game design using Keras.                                                                                                                           |
| [Translator_of_Old_Korean_Literature.ipynb](Gemma/Translator_of_Old_Korean_Literature.ipynb)                         | Use Gemma to translate old Korean literature using Keras.                                                                                                                               |
| [Gemma2_on_Groq.ipynb](Gemma/Gemma2_on_Groq.ipynb)                                                                   | Leverage the free Gemma 2 9B IT model hosted on [Groq](https://groq.com/) (super fast speed).                                                                                           |
| [Run_with_Ollama.ipynb](Gemma/Run_with_Ollama.ipynb)                                                                 | Run Gemma models using [Ollama](https://www.ollama.com/).                                                                                                                               |
| [Using_Gemma_with_Llamafile.ipynb](Gemma/Using_Gemma_with_Llamafile.ipynb)                                           | Run Gemma models using [Llamafile](https://github.com/Mozilla-Ocho/llamafile/).                                                                                                         |
| [Using_Gemma_with_LlamaCpp.ipynb](Gemma/Using_Gemma_with_LlamaCpp.ipynb)                                             | Run Gemma models using [LlamaCpp](https://github.com/abetlen/llama-cpp-python/).                                                                                                        |
| [Using_Gemma_with_LocalGemma.ipynb](Gemma/Using_Gemma_with_LocalGemma.ipynb)                                         | Run Gemma models using [Local Gemma](https://github.com/huggingface/local-gemma/).                                                                                                      |
| [Using_Gemini_and_Gemma_with_RouteLLM.ipynb](Gemma/Using_Gemini_and_Gemma_with_RouteLLM.ipynb)                       | Route Gemma and Gemini models using [RouteLLM](https://github.com/lm-sys/RouteLLM/).                                                                                                    |
| [Using_Gemma_with_SGLang.ipynb](Gemma/Using_Gemma_with_SGLang.ipynb)                                                 | Run Gemma models using [SGLang](https://github.com/sgl-project/sglang/).                                                                                                                |
| [Constrained_generation_with_Gemma.ipynb](Gemma/Constrained_generation_with_Gemma.ipynb)                             | Constrained generation with Gemma models using [LlamaCpp](https://github.com/abetlen/llama-cpp-python/) and [Guidance](https://github.com/guidance-ai/guidance/tree/main/).             |
| [Integrate_with_Mesop.ipynb](Gemma/Integrate_with_Mesop.ipynb)                                                       | Integrate Gemma with [Google Mesop](https://google.github.io/mesop/).                                                                                                                   |
| [Integrate_with_OneTwo.ipynb](Gemma/Integrate_with_OneTwo.ipynb)                                                     | Integrate Gemma with [Google OneTwo](https://github.com/google-deepmind/onetwo).                                                                                                        |
| [Deploy_with_vLLM.ipynb](Gemma/Deploy_with_vLLM.ipynb)                                                               | Deploy a Gemma model using [vLLM](https://github.com/vllm-project/vllm).                                                                                                                |
| [Deploy_Gemma_in_Vertex_AI.ipynb](Gemma/Deploy_Gemma_in_Vertex_AI.ipynb)                                             | Deploy a Gemma model using [Vertex AI](https://cloud.google.com/vertex-ai).                                                                                                             |
| **Prompting**                                                                                                        |                                                                                                                                                                                         |
| [Prompt_chaining.ipynb](Gemma/Prompt_chaining.ipynb)                                                                 | Illustrate prompt chaining and iterative generation with Gemma.                                                                                                                         |
| [LangChain_chaining.ipynb](Gemma/LangChain_chaining.ipynb)                                                           | Illustrate LangChain chaining  with Gemma.                                                                                                                                              |
| [Advanced_Prompting_Techniques.ipynb](Gemma/Advanced_Prompting_Techniques.ipynb)                                     | Illustrate advanced prompting techniques with Gemma.                                                                                                                                    |
| **RAG**                                                                                                              |                                                                                                                                                                                         |
| [RAG_with_ChromaDB.ipynb](Gemma/RAG_with_ChromaDB.ipynb)                                                             | Build a Retrieval Augmented Generation (RAG) system with Gemma using [ChromaDB](https://www.trychroma.com/) and [Hugging Face](https://huggingface.co/).                                |
| [Minimal_RAG.ipynb](Gemma/Minimal_RAG.ipynb)                                                                         | Minimal example of building a RAG system with Gemma using [Google UniSim](https://github.com/google/unisim) and [Hugging Face](https://huggingface.co/).                                |
| [RAG_PDF_Search_in_multiple_documents_on_Colab.ipynb](Gemma/RAG_PDF_Search_in_multiple_documents_on_Colab.ipynb)     | RAG PDF Search in multiple documents using Gemma 2 2B on Google Colab.                                                                                                                  |
| [Using_Gemma_with_LangChain.ipynb](Gemma/Using_Gemma_with_LangChain.ipynb)                                           | Examples to demonstrate using Gemma with [LangChain](https://www.langchain.com/).                                                                                                       |
| [Using_Gemma_with_Elasticsearch_and_LangChain.ipynb](Gemma/Using_Gemma_with_Elasticsearch_and_LangChain.ipynb)       | Example to demonstrate using Gemma with [Elasticsearch](https://www.elastic.co/elasticsearch/), [Ollama](https://www.ollama.com/) and [LangChain](https://www.langchain.com/).          |
| [Gemma_with_Firebase_Genkit_and_Ollama.ipynb](Gemma/Gemma_with_Firebase_Genkit_and_Ollama.ipynb)                     | Example to demonstrate using Gemma with [Firebase Genkit](https://firebase.google.com/docs/genkit/) and [Ollama](https://www.ollama.com/)                                               |
| [Gemma_RAG_LlamaIndex.ipynb](Gemma/Gemma_RAG_LlamaIndex.ipynb)                                                       | RAG example with [LlamaIndex](https://www.llamaindex.ai/) using Gemma.                                                                                                                  |
| **Finetuning**                                                                                                       |                                                                                                                                                                                         |
| [Finetune_with_Axolotl.ipynb](Gemma/Finetune_with_Axolotl.ipynb)                                                     | Finetune Gemma using [Axolotl](https://github.com/OpenAccess-AI-Collective/axolotl).                                                                                                    |
| [Finetune_with_XTuner.ipynb](Gemma/Finetune_with_XTuner.ipynb)                                                       | Finetune Gemma using [XTuner](https://github.com/InternLM/xtuner).                                                                                                                      |
| [Finetune_with_LLaMA_Factory.ipynb](Gemma/Finetune_with_LLaMA_Factory.ipynb)                                         | Finetune Gemma using [LLaMA-Factory](https://github.com/hiyouga/LLaMA-Factory).                                                                                                         |
| [Finetune_with_Torch_XLA.ipynb](Gemma/Finetune_with_Torch_XLA.ipynb)                                                 | Finetune Gemma using [PyTorch/XLA](https://github.com/pytorch/xla).                                                                                                                     |
| [Custom_Vocabulary.ipynb](Gemma/Custom_Vocabulary.ipynb)                                                             | Demonstrate how to use a custom vocabulary "&lt;unused[0-98]&gt;" tokens in Gemma.                                                                                                      |
| **Alignment**                                                                                                        |                                                                                                                                                                                         |
| [Aligning_DPO_Gemma_2b_it.ipynb](Gemma/Aligning_DPO_Gemma_2b_it.ipynb)                                               | Demonstrate how to align a Gemma model using DPO (Direct Preference Optimization) with [Hugging Face TRL](https://huggingface.co/docs/trl/en/index).                                    |
| **Evaluation**                                                                                                       |                                                                                                                                                                                         |
| [Gemma_evaluation.ipynb](Gemma/Gemma_evaluation.ipynb)                                                               | Demonstrate how to use Eleuther AI's LM evaluation harness to perform model evaluation on Gemma.                                                                                        |
| **Mobile**                                                                                                           |                                                                                                                                                                                         |
| [Gemma on Android](Gemma/Gemma-on-Android)                                                                           | Android app to deploy fine-tuned Gemma-2B-it model using MediaPipe LLM Inference API.                                                                                                   |

#### PaliGemma

| **Inference**                                                                                                                                    |                                                                                                                                                                     |
| :----------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Image_captioning_using_PaliGemma.ipynb](PaliGemma/Image_captioning_using_PaliGemma.ipynb)                                                       | Use PaliGemma to generate image captions using Keras.                                                                                                               |
| [Image_captioning_using_finetuned_PaliGemma.ipynb](PaliGemma/Image_captioning_using_finetuned_PaliGemma.ipynb)                                   | Compare the image captioning results using different PaliGemma versions with [Hugging Face](https://huggingface.co/).                                               |
| [Finetune_PaliGemma_for_image_description.ipynb](PaliGemma/Finetune_PaliGemma_for_image_description.ipynb)                                       | Finetune PaliGemma for image description using [JAX](https://github.com/google/jax).                                                                                |
| [Integrate_PaliGemma_with_Mesop.ipynb](PaliGemma/Integrate_PaliGemma_with_Mesop.ipynb)                                                           | Integrate PaliGemma with [Google Mesop](https://google.github.io/mesop/).                                                                                           |
| [Zero_shot_object_detection_in_images_using_PaliGemma.ipynb](PaliGemma/Zero_shot_object_detection_in_images_using_PaliGemma.ipynb)               | Zero-shot Object Detection in images using PaliGemma.                                                                                                               |
| [Zero_shot_object_detection_in_videos_using_PaliGemma.ipynb](PaliGemma/Zero_shot_object_detection_in_videos_using_PaliGemma.ipynb)               | Zero-shot Object Detection in videos using PaliGemma.                                                                                                               |
| [Referring_expression_segmentation_in_images_using_PaliGemma.ipynb](PaliGemma/Referring_expression_segmentation_in_images_using_PaliGemma.ipynb) | Referring Expression Segmentation in images using PaliGemma.                                                                                                        |
| [Referring_expression_segmentation_in_videos_using_PaliGemma.ipynb](PaliGemma/Referring_expression_segmentation_in_videos_using_PaliGemma.ipynb) | Referring Expression Segmentation in videos using PaliGemma.                                                                                                        |
| **Finetuning**                                                                                                                                   |                                                                                                                                                                     |
| [Finetune_PaliGemma_with_Keras.ipynb](PaliGemma/Finetune_PaliGemma_with_Keras.ipynb)                                                             | Finetune PaliGemma with Keras.                                                                                                                                      |
| **Mobile**                                                                                                                                       |                                                                                                                                                                     |
| [PaliGemma on Android](PaliGemma/PaliGemma-on-Android)                                                                                           | Inference PaliGemma on Android using Hugging Face and Gradio Client API for tasks like zero-shot object detection, image captioning, and visual question-answering. |

#### CodeGemma
| **Finetuning**                                                                                 |                                                  |
| :--------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| [CodeGemma_finetuned_on_SQL_with_HF.ipynb](CodeGemma/CodeGemma_finetuned_on_SQL_with_HF.ipynb) | Fine-Tuning CodeGemma on the SQL Spider Dataset. |

## Get help
Ask a Gemma cookbook-related question on the [developer forum](https://discuss.ai.google.dev/c/gemma/10), or open an [issue](https://github.com/google-gemini/gemma-cookbook/issues) on GitHub.

## Wish list
If you want to see additional cookbooks implemented for specific features/integrations, please send us a pull request by adding your feature request(s) in the [wish list](https://github.com/google-gemini/gemma-cookbook/blob/main/WISHLIST.md).

If you want to make contributions to the Gemma Cookbook project, you are welcome to pick any idea in the [wish list](https://github.com/google-gemini/gemma-cookbook/blob/main/WISHLIST.md) and implement it.

## Contributing
Contributions are always welcome. Please read [contributing](https://github.com/google-gemini/gemma-cookbook/blob/main/CONTRIBUTING.md) before implementation.

Thank you for developing with Gemma! We’re excited to see what you create.

## Translation of this repository
* [Traditional Chinese](https://github.com/doggy8088/gemma-cookbook)
* [Simplified Chinese](https://github.com/xiaoxiong1006/gemma-cookbook)
