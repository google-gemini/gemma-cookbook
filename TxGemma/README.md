# TxGemma

This folder is organized into several categories, each focusing on a specific aspect of working with TxGemma models:

* [Inference and serving](#inference-and-serving): How to load, run, and deploy TxGemma models for inference
* [Fine-tuning](#fine-tuning): How to fine-tune TxGemma models for specific tasks and domains
* [Agentic](#agentic): How to integrate TxGemma models into agentic workflows

For more information on TxGemma, please visit the [HAI-DEF developer site](https://developers.devsite.corp.google.com/health-ai-developer-foundations/txgemma).

**Citation**

```bibtex
@article{wang2025txgemma,
    title={TxGemma: Efficient and Agentic LLMs for Therapeutics},
    author={Wang, Eric and Schmidgall, Samuel and Jaeger, Paul F. and Zhang, Fan and Pilgrim, Rory and Matias, Yossi and Barral, Joelle and Fleet, David and Azizi, Shekoofeh},
    year={2025},
}
```

Find the paper [here](https://arxiv.org/abs/2504.06196).

## Inference and serving

| Notebook Name | Description |
| :------------------------------------------------------------------------------------------| --------------------------------------------------------------------------------------------------------------------------------- |
| [[TxGemma]Quickstart_with_Hugging_Face.ipynb]([TxGemma]Quickstart_with_Hugging_Face.ipynb) | Load and run TxGemma using [Hugging Face](https://huggingface.co/).                                                               |
| [[TxGemma]Quickstart_with_Model_Garden.ipynb]([TxGemma]Quickstart_with_Model_Garden.ipynb) | Deploy TxGemma using [Model Garden](https://console.cloud.google.com/vertex-ai/model-garden) and get online or batch predictions. |

## Fine-tuning
| Notebook Name | Description |
| :--------------------------------------------------------------------------------------| ------------------------------------------------------------------------------------------------------------- |
| [[TxGemma]Finetune_with_Hugging_Face.ipynb]([TxGemma]Finetune_with_Hugging_Face.ipynb) | Fine-tune TxGemma on the [TrialBench](https://arxiv.org/abs/2407.00631) dataset using Hugging Face libraries. |

## Agentic
| Notebook Name | Description |
| :--------------------------------------------------------------------------------------------- | ------------------------------------------------- |
| [[TxGemma]Agentic_Demo_with_Hugging_Face.ipynb]([TxGemma]Agentic_Demo_with_Hugging_Face.ipynb) | Use Agentic-Tx, a therapeutics-focused LLM agent. |
