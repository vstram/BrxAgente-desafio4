# Guia de Estilo - BrxAgente VR/VA

## 1. Paleta de Cores

### 1.1 Cores Primárias
| Nome | HEX | RGB | Uso |
|------|-----|-----|-----|
| Azul Principal | `#2196f3` | `rgb(33, 150, 243)` | Cores primárias da marca, botões primários, links |
| Azul Escuro | `#1976d2` | `rgb(25, 118, 210)` | Hover de botões primários, variações escuras |
| Azul Claro | `#bbdefb` | `rgb(187, 222, 251)` | Fundos, estados hover suaves |

### 1.2 Cores Secundárias
| Nome | HEX | RGB | Uso |
|------|-----|-----|-----|
| Verde Principal | `#4caf50` | `rgb(76, 175, 80)` | Cores de sucesso, botões secundários |
| Verde Escuro | `#388e3c` | `rgb(56, 142, 60)` | Hover de botões secundários, variações escuras |
| Verde Claro | `#c8e6c9` | `rgb(200, 230, 201)` | Fundos de sucesso, estados hover suaves |

### 1.3 Cores de Destaque
| Nome | HEX | RGB | Uso |
|------|-----|-----|-----|
| Laranja | `#ff9800` | `rgb(255, 152, 0)` | Avisos, botões de ação importante |
| Vermelho | `#f44336` | `rgb(244, 67, 54)` | Erros, ações destrutivas |

### 1.4 Cores Neutras
| Nome | HEX | RGB | Uso |
|------|-----|-----|-----|
| Fundo Principal | `#f5f7fa` | `rgb(245, 247, 250)` | Cor de fundo da aplicação |
| Superfícies | `#ffffff` | `rgb(255, 255, 255)` | Cartões, modais, elementos flutuantes |
| Texto Principal | `#212121` | `rgb(33, 33, 33)` | Texto principal e títulos |
| Texto Secundário | `#757575` | `rgb(117, 117, 117)` | Texto secundário, descrições |
| Divisores | `#e0e0e0` | `rgb(224, 224, 224)` | Bordas, linhas divisórias |

### 1.5 Sombras
| Nome | Valor | Uso |
|------|-------|-----|
| Sombra 1 | `0 2px 4px rgba(0,0,0,0.1)` | Elementos elevados levemente |
| Sombra 2 | `0 4px 8px rgba(0,0,0,0.1)` | Hover em elementos interativos |
| Sombra 3 | `0 8px 16px rgba(0,0,0,0.1)` | Modais, elementos flutuantes importantes |

## 2. Tipografia

### 2.1 Fonte Principal
A aplicação utiliza a fonte **Inter**, uma tipografia moderna e altamente legível em telas. A fonte está incluída localmente no projeto para garantir carregamento rápido e disponibilidade offline.

### 2.2 Hierarquia Tipográfica
| Elemento | Tamanho | Peso | Cor | Uso |
|----------|---------|------|-----|-----|
| Títulos Principais (H1) | `1.8rem` | `700` | `var(--text-primary)` | Títulos de seção principais |
| Títulos Secundários (H2) | `1.2rem` | `600` | `var(--surface-color)` | Títulos de cartões e seções |
| Subtítulos | `1rem` | `400` | `var(--text-secondary)` | Subtítulos e descrições |
| Corpo do Texto | `1rem` | `400` | `var(--text-primary)` | Texto principal do conteúdo |
| Texto Secundário | `0.9rem` | `400` | `var(--text-secondary)` | Textos auxiliares e notas |

### 2.3 Espaçamento entre Linhas
- Padrão: `1.6` (160% do tamanho da fonte)

## 3. Componentes

### 3.1 Botões
#### Estrutura Visual
- Altura mínima: `44px` (para acessibilidade)
- Padding: `0.75rem 1.5rem`
- Borda arredondada: `4px`
- Sombra: `var(--shadow-1)`
- Transição suave: `all 0.3s ease`

#### Estados
- **Normal**: Cor de fundo apropriada, texto branco
- **Hover**: Elevação com `transform: translateY(-2px)` e `box-shadow: var(--shadow-2)`
- **Ativo**: Retorno à posição original com `transform: translateY(0)`
- **Desabilitado**: Opacidade reduzida (`opacity: 0.6`) e cursor `not-allowed`

#### Variações
- **Primário**: Fundo `var(--primary-color)`, hover `var(--primary-dark)`
- **Secundário**: Fundo `var(--secondary-color)`, hover `var(--secondary-dark)`
- **Rodapé**: Fundo `var(--accent-color)`, hover `#f57c00`

### 3.2 Campos de Entrada
#### Estrutura Visual
- Altura: `44px` (para acessibilidade)
- Padding: `0.75rem`
- Borda: `2px solid var(--divider-color)`
- Borda arredondada: `8px`
- Fundo: `var(--surface-color)`
- Sombra: `var(--shadow-1)`

#### Estados
- **Normal**: Borda `var(--divider-color)`
- **Focus**: Borda `var(--primary-color)` e `box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.2)`
- **Hover**: Borda `var(--primary-light)` (exceto quando focado)

### 3.3 Cartões (Sections)
#### Estrutura Visual
- Fundo: `var(--surface-color)`
- Borda arredondada: `8px`
- Sombra: `var(--shadow-1)`
- Transição suave: `all 0.3s ease`

#### Estados
- **Normal**: Sombra `var(--shadow-1)`
- **Hover**: Elevação com `transform: translateY(-5px)` e `box-shadow: var(--shadow-2)`

### 3.4 Ícones
#### Estrutura Visual
- Tamanho padrão: `1.2rem` (aproximadamente 19px)
- Tamanho em botões: `1.2rem`
- Tamanho em dispositivos móveis: `1rem`

#### Cores
- Herdam a cor do texto do elemento pai
- Em botões: Contraste apropriado com o fundo do botão

## 4. Espaçamento e Grid

### 4.1 Sistema de Espaçamento
Utilizamos um sistema baseado em múltiplos de 0.5rem:
- XS: `0.25rem` (4px)
- S: `0.5rem` (8px)
- M: `1rem` (16px)
- L: `1.5rem` (24px)
- XL: `2rem` (32px)
- XXL: `3rem` (48px)

### 4.2 Grid Principal
- **Desktop**: 12 colunas com gutters de `1.5rem`
- **Tablet**: 8 colunas com gutters de `1rem`
- **Mobile**: 4 colunas com gutters de `0.75rem`

### 4.3 Padding e Margens
- **Containers principais**: `1rem` em todos os lados
- **Elementos internos**: `1rem` entre elementos principais
- **Elementos agrupados**: `0.5rem` entre elementos relacionados

## 5. Animações e Transições

### 5.1 Princípios
- Todas as transições devem ser suaves e rápidas (entre 0.2s e 0.5s)
- Usar funções de tempo apropriadas (`ease`, `ease-in-out`)
- Priorizar transformações em vez de propriedades que disparam repaints

### 5.2 Animações Comuns
- **Fade In**: `fadeIn 0.5s ease-out`
- **Slide Down**: `slideDown 0.5s ease-out`
- **Scale In**: `scaleIn 0.8s ease-out`
- **Pulse**: `pulse 2s infinite` (para estados de loading)

## 6. Responsividade

### 6.1 Breakpoints
- **Mobile**: Até `768px`
- **Tablet**: `769px` até `1024px`
- **Desktop**: Acima de `1025px`

### 6.2 Adaptações
- Redução de tamanhos de fonte em telas menores
- Reorganização de layouts complexos para colunas simples
- Ajuste de padding e margens para telas menores
- Botões e elementos interativos mantêm tamanho mínimo de toque (44px)

## 7. Acessibilidade

### 7.1 Contraste de Cores
- Todos os elementos seguem a diretriz WCAG 2.1 de contraste mínimo de 4.5:1
- Cores foram testadas para daltonismo
- Estados de foco claramente visíveis

### 7.2 Navegação por Teclado
- Todos os elementos interativos são acessíveis via Tab
- Ordem de tabulação lógica e previsível
- Indicador de foco visível em todos os elementos interativos

### 7.3 Leitores de Tela
- Estrutura semântica adequada com cabeçalhos hierárquicos
- Labels apropriados para todos os elementos interativos
- Atributos ARIA onde necessário para contexto adicional

## 8. Ícones Personalizados

### 8.1 Conjunto de Ícones
A aplicação utiliza ícones personalizados em formato SVG para garantir qualidade em qualquer resolução:
- `folder.svg`: Seleção de diretório
- `play.svg`: Iniciar processamento
- `settings.svg`: Configurações
- `spinner.svg`: Estado de carregamento
- `chat.svg`: Abrir chat
- `send.svg`: Enviar mensagem
- `test.svg`: Testar conexão
- `save.svg`: Salvar configurações
- `cancel.svg`: Cancelar ação
- `x.svg`: Fechar modais
- `clear.svg`: Limpar chat
- `info.svg`: Informações
- `arrow-up.svg`: Setas para expandir
- `arrow-down.svg`: Setas para recolher
- `refresh.svg`: Atualizar/recarregar

### 8.2 Uso
- Tamanho padrão: `1.2rem`
- Herdam a cor do texto do elemento pai
- Devem sempre ser acompanhados de texto para acessibilidade
- Em botões, usados como elementos visuais complementares ao texto

---
*Este guia de estilo foi criado para padronizar o desenvolvimento da aplicação BrxAgente VR/VA e deve ser seguido por todos os desenvolvedores que contribuem para o projeto.*