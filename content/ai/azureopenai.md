---
date: 2026-09-22
---
## Preparation

```python
import os
import requests
import azure.cognitiveservices.speech as speechsdk

from openai import AzureOpenAI
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
```

## Azure OpenAI

### GPT-35-turbo

```python
client = AzureOpenAI(
  azure_endpoint = os.getenv("AZURE_OPENAI_ENDPOINT"),
  api_key=os.getenv("AZURE_OPENAI_API_KEY"),
  api_version="2024-02-01"
)

response = client.chat.completions.create(
  model="gpt-35-turbo", # model = "deployment_name".
  messages=[
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Does Azure OpenAI support customer managed keys?"},
    {"role": "assistant", "content": "Yes, customer managed keys are supported by Azure OpenAI."},
    {"role": "user", "content": "Do other Azure AI services support this too?"}
  ]
)

print(response.choices[0].message.content)
```

### Whisper

```python
client = AzureOpenAI(
    api_key=os.getenv("AZURE_OPENAI_API_KEY"),
    api_version="2024-02-01",
    azure_endpoint = os.getenv("AZURE_OPENAI_ENDPOINT")
)

deployment_id = "whisper" #This will correspond to the custom name you chose for your deployment when you deployed a model."
audio_test_file = "./audio/audio.wav"

result = client.audio.transcriptions.create(
    file=open(audio_test_file, "rb"),
    model=deployment_id
)

print(result)
```

## Azure AI Speech

### Speech-to-Text

[圧縮ファイルの場合](https://learn.microsoft.com/ja-jp/azure/ai-services/speech-service/how-to-use-codec-compressed-audio-input-streams?tabs=linux%2Cdebian%2Cjava-android%2Cterminal&pivots=programming-language-python)

```python
def recognize_from_microphone():
  # This example requires environment variables named "SPEECH_KEY" and "SPEECH_REGION"
  speech_config = speechsdk.SpeechConfig(subscription=os.environ.get('AZURE_AI_SPEECH_KEY'), region=os.environ.get('AZURE_AI_SPEECH_REGION'))
  speech_config.speech_recognition_language="ja-JP"

  # use microphone
  # audio_config = speechsdk.audio.AudioConfig(use_default_microphone=True)

  # use audio file
  audio_config = speechsdk.audio.AudioConfig(filename="./audio/audio.wav")
  speech_recognizer = speechsdk.SpeechRecognizer(speech_config=speech_config, audio_config=audio_config)

  print("Speak into your microphone.")
  speech_recognition_result = speech_recognizer.recognize_once_async().get()

  if speech_recognition_result.reason == speechsdk.ResultReason.RecognizedSpeech:
    print("Recognized: {}".format(speech_recognition_result.text))
  elif speech_recognition_result.reason == speechsdk.ResultReason.NoMatch:
    print("No speech could be recognized: {}".format(speech_recognition_result.no_match_details))
  elif speech_recognition_result.reason == speechsdk.ResultReason.Canceled:
    cancellation_details = speech_recognition_result.cancellation_details
    print("Speech Recognition canceled: {}".format(cancellation_details.reason))
    if cancellation_details.reason == speechsdk.CancellationReason.Error:
      print("Error details: {}".format(cancellation_details.error_details))
      print("Did you set the speech resource key and region values?")

recognize_from_microphone()
```

### Continuous Speech-to-Text(Real-time)

```python
def recognize_from_microphone():
  # This example requires environment variables named "SPEECH_KEY" and "SPEECH_REGION"
  speech_config = speechsdk.SpeechConfig(subscription=os.environ.get('AZURE_AI_SPEECH_KEY'), region=os.environ.get('AZURE_AI_SPEECH_REGION'))
  speech_config.speech_recognition_language="ja-JP"

  # use microphone
  audio_config = speechsdk.audio.AudioConfig(use_default_microphone=True)
  speech_recognizer = speechsdk.SpeechRecognizer(speech_config=speech_config, audio_config=audio_config)

  print("Speak into your microphone.")

  done = False
  last_time_speech_detected = time.time()

  def stop_cb(evt):
    print('CLOSING on {}'.format(evt))
    speech_recognizer.stop_continuous_recognition()
    nonlocal done
    done = True

  def update_last_time_speech_detected(evt):
    nonlocal last_time_speech_detected
    last_time_speech_detected = time.time()
    print('RECOGNIZING: {}'.format(evt))

  speech_recognizer.recognizing.connect(update_last_time_speech_detected)
  speech_recognizer.recognized.connect(update_last_time_speech_detected)
  speech_recognizer.session_started.connect(lambda evt: print('SESSION STARTED: {}'.format(evt)))
  speech_recognizer.session_stopped.connect(lambda evt: print('SESSION STOPPED {}'.format(evt)))
  speech_recognizer.canceled.connect(lambda evt: print('CANCELED {}'.format(evt)))

  speech_recognizer.session_stopped.connect(stop_cb)
  speech_recognizer.canceled.connect(stop_cb)

  speech_recognizer.start_continuous_recognition()
  while not done:
    if time.time() - last_time_speech_detected > 5:
      done = True
    time.sleep(0.5)

  speech_recognizer.stop_continuous_recognition()

recognize_from_microphone()
```

### Batch Request

#### Get All Speech Model

```python
url = "https://japaneast.api.cognitive.microsoft.com/speechtotext/v3.2-preview.2/models/base"
headers = {
    "Ocp-Apim-Subscription-Key": os.environ.get('AZURE_AI_SPEECH_KEY')
}

response = requests.get(url, headers=headers)

print(response.text)
```

#### Batch Request - Create

[Whisper Model](https://learn.microsoft.com/en-us/azure/ai-services/speech-service/batch-transcription-create?pivots=rest-api#use-a-whisper-model)

```python
url = "https://japaneast.api.cognitive.microsoft.com/speechtotext/v3.2-preview.2/transcriptions"
headers = {
    "Ocp-Apim-Subscription-Key": os.environ.get('AZURE_AI_SPEECH_KEY'),
    "Content-Type": "application/json"
}

data = {
  "displayName": "My Transcription",
  "description": "Speech Studio Batch speech to text",
  "locale": "ja-jp",
  "contentUrls": [
    "https://{AUDIO_FILE_SAS_URL}",
    "https://{AUDIO_FILE_SAS_URL}"
  ],
  "model": {
    "self": "https://japaneast.api.cognitive.microsoft.com/speechtotext/v3.2-preview.2/models/base/466fe29a-fc93-45c6-8ec3-ff3c7c36d1b9"
  },
  "properties": {
    "wordLevelTimestampsEnabled": False,
    "displayFormWordLevelTimestampsEnabled": True,
    "diarizationEnabled": True,
    "diarization": {
      "speakers": {
        "minCount": 1,
        "maxCount": 10
      }
    },
    "punctuationMode": "DictatedAndAutomatic",
    "profanityFilterMode": "Masked",
    "destinationContainerUrl": "https://{CONTAINER_SAS_URL}"
  },
  "customProperties": {}
}

response = requests.post(url, json=data, headers=headers)

print(response.text)
```

#### Batch Request - Check Status

```python
url = "https://japaneast.api.cognitive.microsoft.com/speechtotext/v3.2-preview.2/transcriptions/{TranscriptionsId}"
headers = {
    "Ocp-Apim-Subscription-Key": os.environ.get('AZURE_AI_SPEECH_KEY')
}

response = requests.get(url, headers=headers)

print(response.text)
```

#### Batch Request - Get Result

```python
url = "https://japaneast.api.cognitive.microsoft.com/speechtotext/v3.2-preview.2/transcriptions/{TranscriptionsId}/files"
headers = {
    "Ocp-Apim-Subscription-Key": os.environ.get('AZURE_AI_SPEECH_KEY')
}

response = requests.get(url, headers=headers)

print(response.text)
```
