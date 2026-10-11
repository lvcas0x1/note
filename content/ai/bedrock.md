---
date: 2026-09-22
---
### CLI

```bash
# list all foundation models
aws bedrock list-foundation-models --region ap-northeast-1

```

### AI Agent

- Planning
  - Planning without Feedback
    - Chain-of-Throught (COT)
    - Tree of Through Reasoning (ToT)
  - Planning with Feedback
    - ReAct or Reasoning + Action
    - Reflexion
- Memory
  - Short-term Memory
  - Long-term Memory
- Tools & Action
  - Enables Agents to connect and retrieve data from external and internal sources
    - Google Search, API, RAG, Code Interpreter, Database
- Guardrails
  - Build responsible AI applications
- Communication (Single/Multi-Agent)
