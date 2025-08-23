# 📖 Guia Completo do Usuário - Agente de IA BrxAgente

Este guia detalha todas as funcionalidades da aplicação desktop BrxAgente e como utilizá-las efetivamente no seu trabalho diário.

## 📑 Índice

1. [Interface da Aplicação](#-interface-da-aplicação)
2. [Funcionalidades Principais](#-funcionalidades-principais)
3. [Workflows Disponíveis](#️-workflows-disponíveis)
4. [Sistema de Chat Inteligente](#-sistema-de-chat-inteligente)
5. [Dashboard e Monitoramento](#-dashboard-e-monitoramento)
6. [Configurações](#️-configurações)
7. [Interpretação de Resultados](#-interpretação-de-resultados)
8. [Melhores Práticas](#-melhores-práticas)

## 🖥️ Interface da Aplicação

### **Layout Principal**

A aplicação desktop BrxAgente possui uma interface moderna dividida em seções:

```
┌─────────────────────────────────────────────────────┐
│  🤖 BrxAgente - Agente de IA para VR               │
├─────────────────┬───────────────────────────────────┤
│                 │  📊 Dashboard Principal           │
│  📋 Menu        │  ┌─────────────────────────────┐   │
│  • Dashboard    │  │ Status: 🟢 Ativo           │   │
│  • Processar VR │  │ Último processo: 2h atrás   │   │
│  • Chat         │  │ Colaboradores: 2.847        │   │
│  • Configurações│  └─────────────────────────────┘   │
│  • Relatórios   │                                   │
│  • Ajuda        │  🔄 Workflow em Execução          │
│                 │  ┌─────────────────────────────┐   │
│                 │  │ ⚙️ Processando VR Outubro    │   │
│                 │  │ Progresso: ████████░░ 80%   │   │
│                 │  │ Etapa atual: Cálculos       │   │
│                 │  └─────────────────────────────┘   │
└─────────────────┴───────────────────────────────────┘
```

### **Seções da Interface**

#### **1. Menu Principal (Sidebar)**
- 📊 **Dashboard**: Visão geral e status
- 🔄 **Processar VR**: Workflows de processamento
- 💬 **Chat**: Consultor inteligente
- ⚙️ **Configurações**: API keys e parâmetros
- 📊 **Relatórios**: Resultados e análises
- ❓ **Ajuda**: Documentação e suporte

#### **2. Área Principal**
- **Status Cards**: Métricas em tempo real
- **Progress Bar**: Acompanhamento de workflows
- **Action Buttons**: Botões para ações principais
- **Data Tables**: Visualização de dados processados

#### **3. Barra de Status (Footer)**
- 🟢 **Conexão**: Status da API LLM
- 📡 **Cache**: Hit ratio e performance
- 💾 **Dados**: Último backup/processamento
- ⚠️ **Alertas**: Notificações importantes

## 🔍 Funcionalidades Principais

### **1. Auditor Inteligente**

**Como Acessar:**
1. Vá para **"Processar VR"** no menu
2. Clique em **"Auditoria Inteligente"**
3. Selecione o diretório das planilhas
4. Clique em **"Iniciar Auditoria"**

**Parâmetros Configuráveis:**
- 🎯 **Sensibilidade**: Baixa, Média, Alta
- 📊 **Limiar de Anomalia**: 1.5x, 2.0x, 2.5x da média
- 🔍 **Incluir Histórico**: Últimos 3, 6, 12 meses
- ⚠️ **Tipos de Alerta**: Crítico, Aviso, Informativo

**Exemplo de Resultado na Interface:**
```
🔍 AUDITORIA CONCLUÍDA

📋 Resumo:
• 2.847 colaboradores analisados
• 16 anomalias detectadas (0.56%)
• Tempo de execução: 3m 45s

🚨 Anomalias Críticas (3):
┌─────────────────────────────────────────────────┐
│ João Silva (MAT001234) - SINDPD                 │
│ • VR: R$ 1.840,00 (340% acima da média)        │
│ • Confiança: 95%                                │
│ • Ação: [Investigar] [Corrigir] [Ignorar]      │
└─────────────────────────────────────────────────┘
```

### **2. Orquestrador de Workflows**

**Interface de Workflows:**
1. **Seletor de Workflow**: Dropdown com opções
2. **Parâmetros**: Campos para configuração
3. **Preview**: Visualização dos passos
4. **Controles**: Iniciar, Pausar, Cancelar

**Workflows Predefinidos:**
- 📈 **Processamento Completo**: Fluxo completo A-Z
- 🔍 **Apenas Auditoria**: Detecção de anomalias
- ⚙️ **Cálculos Isolados**: Apenas cálculos de VR
- 📊 **Geração de Relatórios**: Outputs e insights

**Controles Durante Execução:**
- ⏸️ **Pausar**: Para revisão manual
- ⏩ **Continuar**: Retomar execução
- ❌ **Cancelar**: Parar completamente
- 🔄 **Reiniciar**: Começar do início

### **3. Consultor de Políticas**

**Como Usar:**
1. Clique em **"Chat"** no menu principal
2. Digite sua pergunta na caixa de texto
3. Aguarde a resposta do agente
4. Use **histórico de conversa** para contexto

**Tipos de Perguntas Suportadas:**

**Elegibilidade e Regras:**
```
👤 "Diretores têm direito a VR?"
🤖 "Não. Diretores não são elegíveis para VR conforme 
    Política VR-2025, Seção 2.3. Isso inclui:
    • Diretores Estatutários
    • Diretores Executivos  
    • Diretores de Área
    
    Exceção: Diretores com salário < R$ 15.000 podem 
    ser incluídos mediante aprovação especial."
```

**Cálculos Específicos:**
```
👤 "Como calcular VR para licença médica de 20 dias?"
🤖 "Para licença médica > 15 dias, usar cálculo proporcional:
    
    📊 Fórmula: (Dias trabalhados ÷ Dias úteis totais) × VR base
    
    📝 Exemplo:
    • Dias úteis outubro: 22 dias
    • Dias trabalhados: 2 dias (22-20)
    • Cálculo: (2 ÷ 22) × R$ 460,00 = R$ 41,82
    
    📚 Fonte: Manual VR-2025, Seção 5.4"
```

### **4. Análise Preditiva**

**Como Acessar:**
1. Vá para **"Relatórios"** → **"Análise Preditiva"**
2. Selecione o **sindicato** para análise
3. Escolha o **período** de previsão (3, 6, 12 meses)
4. Clique em **"Gerar Previsão"**

**Outputs Visuais:**
- 📈 **Gráficos de Tendência**: Evolução temporal
- 🥧 **Breakdown por Sindicato**: Distribuição percentual
- 📊 **Métricas Preditas**: Valores futuros estimados
- ⚠️ **Alertas Proativos**: Riscos identificados

## ⚙️ Workflows Disponíveis

### **1. Workflow: Processamento Completo de VR**

**Passos Visualizados na Interface:**
```
1. 📂 Validação do Diretório
   ├── Status: ✅ Concluído
   ├── Tempo: 15s
   └── Arquivos: 3 planilhas encontradas

2. 📊 Leitura das Planilhas
   ├── Status: ✅ Concluído  
   ├── Tempo: 45s
   └── Dados: 2.847 colaboradores

3. 🔍 Auditoria de Dados
   ├── Status: 🔄 Em execução...
   ├── Tempo: 2m 30s
   └── Progresso: 75%

4. ⚙️ Cálculos de VR
   ├── Status: ⏳ Aguardando
   └── Estimativa: 3m

5. 📋 Geração de Relatórios
   ├── Status: ⏳ Aguardando
   └── Estimativa: 1m

6. 📧 Notificações
   ├── Status: ⏳ Aguardando
   └── Estimativa: 15s
```

**Controles Interativos:**
- **Pausar após cada etapa**: Checkbox para revisão manual
- **Notificar por email**: Configurar destinatários
- **Backup automático**: Salvar estados intermediários

### **2. Workflow: Detecção de Anomalias**

**Interface Especializada:**
```
🔍 DETECTOR DE ANOMALIAS

📊 Configurações:
┌─────────────────────────────────────────┐
│ Sensibilidade:     [●●●○○] Alta         │
│ Limiar de Desvio:  [2.5x] da média     │
│ Incluir Histórico: [6 meses]           │
│ Tipos de Anomalia: [✓] Valor [✓] Data  │
└─────────────────────────────────────────┘

🎯 Resultados:
┌─────────────────────────────────────────┐
│ 🚨 3 Anomalias Críticas                │
│ 🟨 8 Alertas de Atenção                │
│ 📊 5 Informações Adicionais            │
│                                         │
│ [📥 Exportar] [📧 Enviar] [🔧 Corrigir] │
└─────────────────────────────────────────┘
```

## 💬 Sistema de Chat Inteligente

### **Interface de Chat**

**Layout do Chat:**
```
┌─────────────────────────────────────────────────────┐
│ 💬 Consultor Inteligente - Agente VR              │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 👤 Quantos colaboradores do SINDPD temos?          │
│                                                     │
│ 🤖 Encontrei 247 colaboradores ativos do SINDPD.   │
│    Destes, 234 são elegíveis para VR (94.7%).     │
│                                                     │
│    📊 Detalhamento:                                 │
│    • Ativos e elegíveis: 234                       │
│    • Em licença: 8                                 │
│    • Estagiários: 3 (não elegíveis)               │
│    • Diretores: 2 (não elegíveis)                 │
│                                                     │
│    💡 Deseja ver distribuição por setor?           │
│                                                     │
├─────────────────────────────────────────────────────┤
│ Digite sua pergunta...                    [Enviar] │
└─────────────────────────────────────────────────────┘
```

**Funcionalidades do Chat:**
- 💾 **Histórico Persistente**: Conversas salvas
- 📋 **Sugestões**: Perguntas comuns predefinidas
- 📊 **Visualizações**: Tabelas e gráficos inline
- 🔗 **Links**: Acesso direto a dados relacionados
- 📤 **Exportação**: Salvar conversas em PDF/Word

### **Comandos Especiais do Chat**

**Comandos de Dados:**
```
@colaboradores [filtro] - Lista colaboradores
@anomalias [período] - Exibe anomalias detectadas
@relatório [tipo] - Gera relatório específico
@previsão [sindicato] - Análise preditiva
@status - Status atual do sistema
```

**Exemplos de Uso:**
```
👤 @colaboradores SINDPD férias
🤖 Colaboradores do SINDPD em férias:
    • João Silva: 15/10 a 29/10
    • Maria Santos: 22/10 a 05/11
    • Pedro Costa: 01/11 a 15/11
    Total: 3 colaboradores

👤 @previsão SINDAC
🤖 📈 Previsão SINDAC (próximos 3 meses):
    • Nov/2025: R$ 245.670 (534 colaboradores)
    • Dez/2025: R$ 251.230 (+2.3% tendência)
    • Jan/2026: R$ 267.890 (+6.6% pico sazonal)
```

## 📊 Dashboard e Monitoramento

### **Dashboard Principal**

**Seções do Dashboard:**
```
┌─────────────┬─────────────┬─────────────┬─────────────┐
│ 👥 Total    │ 💰 VR Mensal│ 🔍 Anomalias│ ⚡ Performance │
│ 2.847       │ R$ 1.3M     │ 3 críticas │ 95% cache   │
│ Colaboradores│ Estimado    │ 8 alertas  │ hit ratio   │
└─────────────┴─────────────┴─────────────┴─────────────┘

┌─────────────────────────────────────────────────────┐
│ 📈 Evolução Mensal                                 │
│     R$ 1.4M ┤                                      │
│     R$ 1.3M ┤     ●───●───●                       │
│     R$ 1.2M ┤   ●               ●                 │
│     R$ 1.1M ┤ ●                                   │
│             └──────────────────────────────────── │
│              Jun  Jul  Ago  Set  Out             │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ 🚨 Alertas Recentes                                │
│ • 15:30 - Anomalia detectada: João Silva (MAT001) │
│ • 14:45 - Processamento VR iniciado                │
│ • 14:20 - Cache otimizado automaticamente          │
│ • 13:55 - Backup automático concluído              │
└─────────────────────────────────────────────────────┘
```

### **Métricas em Tempo Real**

**Performance do Sistema:**
- 🚀 **Throughput**: Colaboradores/minuto processados
- 💾 **Cache Hit Ratio**: % de acertos no cache
- ⚡ **Tempo de Resposta**: Latência média das consultas
- 🧠 **Uso de Memória**: Consumo atual vs. disponível

**Métricas de Negócio:**
- 📊 **Taxa de Sucesso**: % de workflows concluídos
- 🔍 **Anomalias por Período**: Tendência de detecção
- 👥 **Colaboradores Processados**: Total acumulado
- 💰 **Economia Estimada**: Tempo/custo poupado

## ⚙️ Configurações

### **Tela de Configurações**

**Abas de Configuração:**
```
┌─────────────────────────────────────────────────────┐
│ [🔑 API] [📁 Arquivos] [🚨 Alertas] [⚡ Performance] │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 🔑 Configuração de API                             │
│                                                     │
│ Provedor LLM: [OpenAI ▼]                           │
│ API Key: [sk-...] [👁️] [✅ Testado]                │
│ Modelo: [gpt-4-turbo ▼]                            │
│ Timeout: [30] segundos                             │
│                                                     │
│ ──────────────────────────                         │
│                                                     │
│ 🔄 Provedor Backup: [Claude ▼]                     │
│ API Key: [sk-ant-...] [👁️] [⚠️ Não testado]        │
│                                                     │
│ [💾 Salvar] [🧪 Testar Conexão] [↻ Restaurar]      │
└─────────────────────────────────────────────────────┘
```

### **Configurações Avançadas**

**Performance:**
- 🔄 **Workers Paralelos**: 4, 8, 16, 32
- 💾 **Tamanho do Cache**: 500MB, 1GB, 2GB
- ⏱️ **TTL do Cache**: 12h, 24h, 48h
- 🎯 **Sensibilidade de Anomalias**: Baixa, Média, Alta

**Alertas e Notificações:**
- 📧 **Email**: SMTP configurável
- 💬 **Slack**: Webhook integration
- 🔔 **Desktop**: Notificações do sistema
- 📱 **WhatsApp**: Via API (opcional)

**Arquivos e Backup:**
- 📂 **Diretório de Trabalho**: Pasta das planilhas
- 💾 **Backup Automático**: Frequência e retenção
- 📊 **Formato de Saída**: Excel, CSV, JSON
- 🗄️ **Histórico**: Período de retenção

## 📊 Interpretação de Resultados

### **Planilha Final de VR**

**Visualização na Interface:**
```
┌─────────────────────────────────────────────────────┐
│ 📊 Resultados - VR Outubro 2025                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 📈 Resumo Executivo:                               │
│ • Total de colaboradores: 2.847                    │  
│ • VR total calculado: R$ 1.309.620,00             │
│ • Média por colaborador: R$ 460,00                 │
│ • Anomalias detectadas: 3 críticas, 8 alertas     │
│                                                     │
│ 📋 Breakdown por Sindicato:                        │
│ ┌─────────┬──────────┬──────────────┬──────────┐    │
│ │Sindicato│Colaborado│ VR Total     │ Média    │    │
│ │         │res       │              │          │    │
│ ├─────────┼──────────┼──────────────┼──────────┤    │
│ │SINDPD   │ 1.247    │ R$ 573.620   │ R$ 460   │    │
│ │SINDAC   │   892    │ R$ 410.320   │ R$ 460   │    │
│ │SINDMET  │   708    │ R$ 325.680   │ R$ 460   │    │
│ └─────────┴──────────┴──────────────┴──────────┘    │
│                                                     │
│ [📥 Exportar Excel] [📧 Enviar] [🔍 Ver Detalhes]  │
└─────────────────────────────────────────────────────┘
```

### **Relatório de Anomalias**

**Interface de Anomalias:**
```
┌─────────────────────────────────────────────────────┐
│ 🚨 Relatório de Anomalias Detectadas               │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 🚨 Críticas (Requer ação imediata)                │
│                                                     │
│ 1️⃣ João Silva Santos (MAT001234)                   │
│    ├─ Sindicato: SINDPD                           │
│    ├─ VR Calculado: R$ 1.840,00                   │
│    ├─ Valor Esperado: R$ 460,00                   │
│    ├─ Desvio: +300% (4x a média)                  │
│    ├─ Confiança: 95%                              │
│    ├─ Possível Causa: Duplicação de registro      │
│    └─ Ação: [🔍 Investigar] [✏️ Corrigir]         │
│                                                     │
│ 2️⃣ Maria Santos (MAT005678)                        │
│    ├─ Problema: Matrícula duplicada               │
│    ├─ Ocorrências: 2 registros ativos             │
│    ├─ Confiança: 100%                             │
│    └─ Ação: [🔗 Consolidar] [❌ Remover]          │
│                                                     │
│ [📊 Ver Todos] [📧 Enviar Relatório] [⚙️ Configurar] │
└─────────────────────────────────────────────────────┘
```

## 🎯 Melhores Práticas

### **Fluxo de Trabalho Recomendado**

1. **📋 Preparação Semanal**
   - ✅ Verifique atualizações das planilhas
   - ✅ Confirme configurações de API
   - ✅ Execute teste de conectividade
   - ✅ Revise alertas pendentes

2. **🔄 Processamento Mensal**
   - ✅ Backup das planilhas originais
   - ✅ Execute "Validação de Planilhas" primeiro
   - ✅ Rode "Detecção de Anomalias"
   - ✅ Processe apenas após resolver críticas
   - ✅ Gere relatórios finais

3. **📊 Monitoramento Contínuo**
   - ✅ Acompanhe dashboard diariamente
   - ✅ Responda alertas em até 24h
   - ✅ Mantenha histórico organizado
   - ✅ Documente exceções tratadas

### **Configurações Recomendadas por Cenário**

**Para Volume Alto (>5000 colaboradores):**
```yaml
performance:
  workers: 16
  cache_size: "2GB"
  batch_size: 500
  
validation:
  strictness: "medium"
  timeout: 60
```

**Para Precisão Máxima (Compliance):**
```yaml
validation:
  strictness: "high"
  manual_review: true
  backup_every_step: true
  
anomaly_detection:
  sensitivity: "high"
  threshold: 1.5
```

**Para Velocidade (Processamento Frequente):**
```yaml
performance:
  cache_ttl: "72h"
  parallel_validation: true
  skip_non_critical: true
  
speed:
  fast_mode: true
  reduce_logging: true
```

---

## 🆘 Suporte e Recursos Adicionais

- 🔧 **[API Reference](api-reference.md)** - Métodos técnicos disponíveis
- ⚙️ **[Workflows](workflows.md)** - Configurações avançadas
- 🆘 **[Troubleshooting](troubleshooting.md)** - Soluções para problemas
- 💡 **[Exemplos](examples/)** - Cases de uso práticos

**Precisa de ajuda?** Use o botão "❓ Ajuda" na aplicação para acessar a documentação contextual!

*A interface gráfica foi projetada para ser intuitiva. Explore e experimente! 🚀*