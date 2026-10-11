---
date: 2026-09-22
---
## Install Jupyter Notebook

```shell
python -m venv example
cd example
source bin/activate

pip install jupyter
jupyter notebook
```

## Get Gemini Key

Open Google AI Studio https://aistudio.google.com/

## Test API Key

```curl
API_KEY="YOUR_API_KEY"
curl -H 'Content-Type: application/json' \
  -d '{"contents":[
    {"role": "user",
      "parts":[{"text": "Give me five subcategories of jazz?"}]}]}' \
  "https://generativelanguage.googleapis.com/v1/models/gemini-pro:generateContent?key=${API_KEY}"
```

## Install Google AI Python SDK

```shell
pip install google-generativeai
pip show google-generativeai
```

## List Available Models

```python
import os
import google.generativeai as genai
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
genai.configure(api_key=os.environ.get('GEMINI_API_KEY'))

# Listing available models
for m in genai.list_models():
  if 'generateContent' in m.supported_generation_methods:
    print(m.name)
```

## Generating Text From Text Inputs: Gemini Pro

```python
model = genai.GenerativeModel('gemini-pro')
response = model.generate_content('Tell me a joke that only an AI would understand')
print(response.text)
# response.prompt_feedback
```

## Streaming Response

```python
prompt = 'Write an imaginary scientific paper that proposes a new theory of quantum gravity'
response = model.generate_content(prompt, stream=True)

for chunk in response:
  print(chunk.text, end='')
```

## Generating Text from Image and Text Inputs: Gemini Pro Vision

```python
from IPython.display import Image
from IPython.core.display import HTML

img = Image('img/dog.jpg')
model = genai.GenerativeModel('gemini-pro-vision')
response = model.generate_content(img)
print(response.text)
```

With prompt

```python
prompt = 'describe this image for me.'
img = Image('img/breakfast.jpg')
response = model.generate_content([prompt, img])
print(response.text)
```

## Gemini API Generation Parameters

candidate_count: 1, 変更不可，回答候補の数を設定する  
stop_sequences: 回答が止まる記号を設定する。最大 5 個  
max_output_tokens: 最大の生成されるトークン数  
temperature: [0.0, 1.0], 創造性を設定する. 1.0 はもっと創造的  
top_p: 確率値を上から順に合計したとき、X 値を超えない範囲で候補が選ばれる  
top_k (int): 確率値を上位 X から選ぶ

```python
import os
import google.generativeai as genai
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
genai.configure(api_key=os.environ.get('GEMINI_API_KEY'))

generation_config = genai.types.GenerationConfig()

# change for all query
model = genai.GenerativeModel('gemini-pro', generation_config=generation_config)

# change for single query
response = model.generate_content(prompt, generation_config=generation_config)

# list all parameters
help(genai.types.GenerationConfig)

# check default value for models
genai.get_model('models/gemini-pro')

# Model(name='models/gemini-pro',
#       base_model_id='',
#       version='001',
#       display_name='Gemini 1.0 Pro',
#       description='The best model for scaling across a wide range of tasks',
#       input_token_limit=30720,
#       output_token_limit=2048,
#       supported_generation_methods=['generateContent', 'countTokens'],
#       temperature=0.9,
#       top_p=1.0,
#       top_k=1)

genai.get_model('models/gemini-pro-vision')

# Model(name='models/gemini-pro-vision',
#       base_model_id='',
#       version='001',
#       display_name='Gemini 1.0 Pro Vision',
#       description='The best image understanding model to handle a broad range of applications',
#       input_token_limit=12288,
#       output_token_limit=4096,
#       supported_generation_methods=['generateContent', 'countTokens'],
#       temperature=0.4,
#       top_p=1.0,
#       top_k=32)
```

Example

```python
import os
import google.generativeai as genai
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
genai.configure(api_key=os.environ.get('GEMINI_API_KEY'))

generation_config = genai.types.GenerationConfig(
  candidate_count=1,
  #stop_sequences=[','],
  max_output_tokens=32000,
  temperature=0.9,
  top_p=1,
  top_k=3
)

model = genai.GenerativeModel('gemini-pro')
prompt = 'Give me an example of an SQL SELECT statement.'
response = model.generate_content(prompt, generation_config=generation_config)
print(response.text)
```

## Building Chat Conversations

```python
model = genai.GenerativeModel('gemini-pro')
chat = model.start_chat(history=[])

response = chat.send_message('In one paragraph, explain how the Internet works to a beginner')

# all response
chat.history

# user message
chat.history[0]

# model response
chat.history[1]

# get only model response text
chat.history[1].parts[0].text

# get last response text
chat.history[-1].parts[0].text
```

Example:

```python
for item in chat.history:
  print(f'{item.role.capitalize()}: {item.parts[0].text}')
  print('-' * 100)
```

## Building a Conversational Agent Using Gemini Pro

```python
import os
import google.generativeai as genai
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
genai.configure(api_key=os.environ.get('GEMINI_API_KEY'))

model = genai.GenerativeModel('gemini-pro')
chat = model.start_chat(history=[])
while True:
    prompt = input('User: ')
    if prompt.lower() not in ['exit', 'quit']:
        response = chat.send_message(prompt)
        print(f'{chat.history[-1].role.capitalize()}: {chat.history[-1].parts[0].text}')
        print('\n' + '*' * 100 + '\n')
    else:
        print('Quitting ...')
        break
```

## Prompt Engineering

```python
import os
import google.generativeai as genai
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)

def generate_response(prompt, temperature=1):

  genai.configure(api_key=os.environ.get('GEMINI_API_KEY'))

  model = genai.GenerativeModel('gemini-pro')

  generation_config = genai.types.GenerationConfig(
      temperature=temperature, top_p=1, top_k=1
  )

  response = model.generate_content(prompt , generation_config=generation_config)

  return response.text

print(generate_response('Create a python function that calculates the factorial.'))
```

```python
text = '''
 '''

prompt = f'''
このテキストをまとめてください: {text}
'''

output = generate_response(prompt)
print(output)
```

````python
prompt = f'''
200文字以内に、3つのバッククォートで区切られたテキストを要約してください。
要約を日本語から中国語に翻訳してください。
Text: ```{text}```
'''
output = generate_response(prompt)
print(output)
````

````python
prompt = f'''
3つのバッククォートで区切られたTextから4択問題の練習問題を生成してください。

以下の指示を必ず守ってください。
3つのバッククォートで区切られたTextを全部読み込んでください。
問題文は20文字以上、100文字以内で生成してください、問題文はあなたがTextを全部理解した上で生成してください、曖昧な質問をしないでください、程度副詞を使う問題文は避けてください。
問題文を理解してください、Textから正解は1つだけ生成してください、正解は4文字以上、40文字以内に生成してください、正解の選択肢はTextの言葉からそのまま引用して生成してください、曖昧な言葉や言葉の添削をしないでください、程度副詞を使う選択肢は避けてください。
問題文を理解してください、Textから不正解の選択肢3つ生成してください、上記正解の段落以外のText部分から生成してください、架空の不正解は生成しないでください。不正解は4文字以上、40文字以内に生成してください、選択肢はTextの言葉からそのまま引用して生成してください、曖昧な言葉や言葉の添削をしないでください、程度副詞を使う選択肢は避けてください。
正解を生成された後、3つのバッククォートで区切られたTextから正解が正しいかどうかを検証してください。
検証成功であれば、正解の理由を生成してください、正解の理由は問題文に対してどうして正解かを記載してください。正解の理由にTextの中の引用元を含めてください。
検証失敗であれば、再度正解を生成してください、検証成功まで繰り返してください。
不正解の3つの選択肢は本当に不正解なのかを検証してください、曖昧な不正解な選択肢を生成しないでください。
検証失敗であれば、再度正解を生成してください、検証成功まで繰り返してください。
検証成功であれば、不正解の選択肢の理由を生成してください、不正解の理由は問題文に対してなぜ、どのように不正解かをTextを理解した上で具体的に記載してください。不正解の理由にTextの中の引用元を含めてください。
アウトプットのフォーマットは3つのバッククォートで区切られたOutputのフォーマットに合わせて作成してください。正解の選択肢は4択の中でランダムで配置してください。選択肢Aを正解の選択肢として出力することはなるべく避けてください。

Text: ```{text}```
Output: ```
問題文:

選択肢:
A.
B.
C.
D.
--------------------------------------------------------------------------------------------------------------------------------------------
正解:
正解の選択肢をここに記載
正解理由を上記正解の選択肢の下に記載

不正解:
不正解の選択肢をここに記載
不正解理由を上記不正解の選択肢の下に記載

全ての不正解の選択肢と理由を繰り返して出力してください
```

'''
output = generate_response(prompt)
print(output)
````

Output

```
問題文:

教科書で学んだ言葉の生活との関係について、正しい記述はどれですか？

選択肢:
A. かけ離れているものが多い
B. 日常生活においても活用できる
C. 教師の都合によって異なる
D. 教育出版の教科書では考慮されていない

--------------------------------------------------------------------------------------------------------------------------------------------
正解:
B. 日常生活においても活用できる
正解理由: 「教育出版の新しい教科書は，話す・聞く，書く，読むという言語活動が相互に関連し合い，日常生活の中で使いこなすことのできる力がつくようになっています。」

不正解:
A. かけ離れているものが多い
不正解理由: 上記正解の理由にあるとおり、教科書は日常生活で使用できる言語能力の習得を目的としているため、かけ離れたものではない。

C. 教師の都合によって異なる
不正解理由: 教科書は教育出版によって統一的に作成されており、教師の都合によって内容が異なることはない。

D. 教育出版の教科書では考慮されていない
不正解理由: 上記問題文のとおり、教科書は「日常の生活とつながるような言葉の力を伸ばすことができるように構成されている」ことを目的としている。
```
