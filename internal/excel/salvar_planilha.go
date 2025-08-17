// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"BrxAgente-desafio4/internal/modelo"
)

// SalvarPlanilhaEmDownloads salva uma planilha Excel na pasta de Downloads do usuário
// com o nome de arquivo especificado
func SalvarPlanilhaEmDownloads(colaboradores map[string]*modelo.Colaborador, nomeArquivo string) error {
	// Obter o diretório de Downloads do usuário
	diretorioDownloads, err := obterDiretorioDownloads()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório de Downloads: %w", err)
	}
	
	// Criar o caminho completo do arquivo
	caminhoArquivo := filepath.Join(diretorioDownloads, nomeArquivo)
	
	// Gerar a planilha no diretório de Downloads
	if err := GerarPlanilhaResultado(colaboradores, caminhoArquivo); err != nil {
		return fmt.Errorf("erro ao gerar planilha: %w", err)
	}
	
	return nil
}

// SalvarPlanilhaEmDownloadsComTemplate salva uma planilha Excel na pasta de Downloads do usuário
// usando um template existente como base
func SalvarPlanilhaEmDownloadsComTemplate(colaboradores map[string]*modelo.Colaborador, caminhoTemplate, nomeArquivo string) error {
	// Obter o diretório de Downloads do usuário
	diretorioDownloads, err := obterDiretorioDownloads()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório de Downloads: %w", err)
	}
	
	// Criar o caminho completo do arquivo
	caminhoArquivo := filepath.Join(diretorioDownloads, nomeArquivo)
	
	// Gerar a planilha no diretório de Downloads usando template
	if err := GerarPlanilhaResultadoComTemplate(colaboradores, caminhoTemplate, caminhoArquivo); err != nil {
		return fmt.Errorf("erro ao gerar planilha com template: %w", err)
	}
	
	return nil
}

// obterDiretorioDownloads obtém o diretório de Downloads do usuário de acordo com o sistema operacional
func obterDiretorioDownloads() (string, error) {
	var diretorio string
	
	switch runtime.GOOS {
	case "windows":
		// No Windows, a pasta Downloads geralmente está em %USERPROFILE%\Downloads
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			return "", fmt.Errorf("variável de ambiente USERPROFILE não definida")
		}
		diretorio = filepath.Join(userProfile, "Downloads")
	case "darwin":
		// No macOS, a pasta Downloads geralmente está em /Users/{username}/Downloads
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("erro ao obter diretório home: %w", err)
		}
		diretorio = filepath.Join(homeDir, "Downloads")
	case "linux":
		// No Linux, a pasta Downloads geralmente está em /home/{username}/Downloads
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("erro ao obter diretório home: %w", err)
		}
		diretorio = filepath.Join(homeDir, "Downloads")
	default:
		// Para outros sistemas, tentar usar o diretório home e assumir a pasta Downloads
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("erro ao obter diretório home: %w", err)
		}
		diretorio = filepath.Join(homeDir, "Downloads")
	}
	
	// Verificar se o diretório existe
	if _, err := os.Stat(diretorio); os.IsNotExist(err) {
		// Tentar criar o diretório se não existir
		if err := os.MkdirAll(diretorio, 0755); err != nil {
			return "", fmt.Errorf("erro ao criar diretório de Downloads: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("erro ao verificar diretório de Downloads: %w", err)
	}
	
	return diretorio, nil
}

// VerificarPermissoesDiretorio verifica se o usuário tem permissões para escrever no diretório especificado
func VerificarPermissoesDiretorio(diretorio string) error {
	// Tentar criar um arquivo temporário no diretório
	arquivoTemp := filepath.Join(diretorio, ".temp_write_test")
	
	// Criar o arquivo
	file, err := os.Create(arquivoTemp)
	if err != nil {
		return fmt.Errorf("sem permissão para escrever no diretório %s: %w", diretorio, err)
	}
	
	// Fechar o arquivo
	if err := file.Close(); err != nil {
		return fmt.Errorf("erro ao fechar arquivo temporário: %w", err)
	}
	
	// Remover o arquivo temporário
	if err := os.Remove(arquivoTemp); err != nil {
		return fmt.Errorf("erro ao remover arquivo temporário: %w", err)
	}
	
	return nil
}