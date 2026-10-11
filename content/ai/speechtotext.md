---
date: 2026-09-22
---
Azure と Google の Speech-to-Text について触ってみたので、現在の仕様と感想について記載しておきます。

[サンプルで使った音源](https://note.com/note_0530/n/neb6a205744bd)

## Google Vertex AI

### Speech Studio

基本は Speech Studio 使って音声変換行う。  
日本リージョンでサポートしているモデルは

- Long
- short
- Chirp2 \*一番賢いモデル
- Chirp
- Telephony

詳細な言語および機能のサポートページは[こちら](https://cloud.google.com/speech-to-text/v2/docs/speech-to-text-supported-languages)

上記音源の WER 率は 7%  
WER については[こちら](https://learn.microsoft.com/en-us/azure/ai-services/speech-service/how-to-custom-speech-evaluate-data?pivots=speech-studio#evaluate-word-error-rate-wer)

チューニング方法：

1. **Model Adaptations**
   > 予め語彙セットをアップロードして、音声変換中利用できる。  
   > ただ現在は上記 Chirp2 モデルサポートされていなく、Long モデルの場合は元々精度が悪すぎて検証してみましたが、Long モデル + Model Adaptations < Chirp2 の感触でした。
2. **Custom Models**
   > 事前にモデルとデータセットをトレーニングすることで精度を向上させる。  
   > ドキュメントでは日本語サポートされていますが、実際コンソール上日本語選択できない（問い合わせ中）  
   > ただ、Long モデルしか使えないのであんまり期待できない。

## Azure Speech-to-Text

### Whisper について

Whisper モデルは現在 AzureOpenAI(japaneast は未対応) と AzureAISpeech 両方から起動できる  
選択基準は Azure の[ドキュメント](https://learn.microsoft.com/ja-jp/azure/ai-services/speech-service/whisper-overview#whisper-model-via-azure-ai-speech-or-via-azure-openai-service)に記載されている

### 選ぶポイント:

| 項目                       | AzureOpenAI Whisper | Azure AI Speech Whisper Mode | Azure AI Speech General Mode | Azure AI Speech Batch Model |
| -------------------------- | ------------------- | ---------------------------- | ---------------------------- | --------------------------- |
| リアルタイム文字起こし     | ×                   | ×                            | ◯                            | ×                           |
| 単一ファイルが 25MB 超える | ×                   | ◯(1GB まで)                  | ◯(1GB まで)                  | ◯(1GB まで)                 |
| ダイアライゼーション       | ×                   | ×                            | ◯                            | ◯                           |

Azure AI Speech で確認した上記音源のの WER 率は 6.13%

チューニング方法：

1. **Phrase List**
   > 事前に単語を設定することで、認識の可能性を高める。
2. **Custom Model**
   > モデルカスタム可能だが、かなりハードルが高く、時間とお金もかかる。  
   > トレーニングとテストデータ作成の詳細は[こちら](https://learn.microsoft.com/ja-jp/azure/ai-services/speech-service/how-to-custom-speech-test-and-train)  
   > 音声+ラベル付けの場合は一つの音声の長さは最大 60s、合計 1-20 時間のデータが推奨されている。  
   > また、東日本リージョンに専用トレーニングハードウェアがないため、さらにトレーニングに時間を要してしまう。（数日〜数週間、保証はない）  
   > トレーニング開始後、完了まで削除できない。

感想:

金額については Azure と Google のバッチ処理はほぼ同じ、28 円/時間  
データ読み込みについてはどちらも Storage から音声文字起こしサポートしている。  
正確率は Google より Azure の方がやや上、AzureAI Speech の Whisper モードと General モードと Batch モードはほぼ同じ。  
カスタムモデルのトレーニングは Azure も Google も現状サポートがあんまり良くない。  
実際案件の場合は Phrase List 使って少し正解率向上目指せますが、それ以上は難しい。他の方法考える必要がある。
