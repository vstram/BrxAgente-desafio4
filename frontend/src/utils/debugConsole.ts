import WailsProfiler from './wailsProfiler';

// Comandos de debug disponíveis no console do navegador
class DebugConsole {
  static init() {
    // Comandos globais disponíveis no window
    (window as any).debugBrx = {
      // Mostrar estatísticas de performance do Wails
      wailsStats: () => {
        WailsProfiler.printAllStats();
      },

      // Limpar estatísticas
      clearStats: () => {
        console.log('🧹 Limpando estatísticas de performance...');
        (WailsProfiler as any).calls = new Map();
        console.log('✅ Estatísticas limpas!');
      },

      // Forçar coleta de lixo (se disponível)
      gc: () => {
        if ((window as any).gc) {
          console.log('🗑️ Executando coleta de lixo...');
          (window as any).gc();
          console.log('✅ Coleta de lixo executada!');
        } else {
          console.log('❌ Coleta de lixo não disponível. Execute Chrome com --js-flags="--expose-gc"');
        }
      },

      // Mostrar uso de memória
      memory: () => {
        const memory = (performance as any).memory;
        if (memory) {
          console.log('💾 Uso de Memória JavaScript:');
          console.log(`   • Heap Usado: ${(memory.usedJSHeapSize / 1024 / 1024).toFixed(2)} MB`);
          console.log(`   • Heap Total: ${(memory.totalJSHeapSize / 1024 / 1024).toFixed(2)} MB`);
          console.log(`   • Limite: ${(memory.jsHeapSizeLimit / 1024 / 1024).toFixed(2)} MB`);
          
          const percentage = (memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100;
          console.log(`   • Utilização: ${percentage.toFixed(1)}%`);
          
          if (percentage > 80) {
            console.warn('⚠️ Uso de memória alto! Possível memory leak.');
          }
        } else {
          console.log('❌ Informações de memória não disponíveis');
        }
      },

      // Simular carga pesada para testar performance
      stressTest: async (iterations = 10) => {
        console.log(`🔥 Executando teste de stress com ${iterations} iterações...`);
        
        const startTime = performance.now();
        
        for (let i = 0; i < iterations; i++) {
          try {
            console.log(`   Iteração ${i + 1}/${iterations}`);
            await WailsProfiler.profile(`StressTest_${i}`, () => 
              import("../../wailsjs/go/main/App").then(m => m.GetAgentStatus())
            );
          } catch (error) {
            console.error(`   ❌ Erro na iteração ${i + 1}:`, error);
          }
        }
        
        const totalTime = performance.now() - startTime;
        console.log(`🏁 Teste de stress concluído em ${totalTime.toFixed(2)}ms`);
        
        // Mostrar estatísticas
        WailsProfiler.printAllStats();
      },

      // Ajuda
      help: () => {
        console.log('🎛️ Comandos de Debug Disponíveis:');
        console.log('================================');
        console.log('debugBrx.wailsStats()  - Mostrar estatísticas Wails');
        console.log('debugBrx.clearStats()  - Limpar estatísticas');
        console.log('debugBrx.memory()      - Mostrar uso de memória');
        console.log('debugBrx.gc()          - Forçar coleta de lixo');
        console.log('debugBrx.stressTest(N) - Teste de stress com N iterações');
        console.log('debugBrx.help()        - Mostrar esta ajuda');
        console.log('');
        console.log('💡 Exemplo: debugBrx.stressTest(20)');
      }
    };

    console.log('🎛️ Debug Console inicializado! Digite debugBrx.help() para ver comandos disponíveis.');
  }
}

export default DebugConsole;