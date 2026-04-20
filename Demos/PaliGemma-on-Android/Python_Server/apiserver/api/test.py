from gradio_client import Client, file, handle_file


def detect():
    client = Client("big-vision/paligemma")
    result = client.predict(
    handle_file('../media/images/person_and_car.jpg'), # filepath in 'Image' Image component
    "Detect person and car",
    "paligemma-3b-mix-224", # str in 'Prompt' Textbox component # Literal[] in 'Model' Dropdown component
    "greedy", # Literal['greedy', 'nucleus(0.1)', 'nucleus(0.3)', 'temperature(0.5)'] in 'Decoding' Dropdown component
    api_name="/compute"
    )
    print(result)
    return {"result": result}

detect()
