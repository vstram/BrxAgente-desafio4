// Profiler para medir performance das chamadas Wails
class WailsProfiler {
  private static calls: Map<string, number[]> = new Map();
  
  static async profile<T>(
    functionName: string, 
    wailsFunction: () => Promise<T>
  ): Promise<T> {
    const startTime = performance.now();
    
    console.log(`🔗 [Wails] ${functionName} - Iniciando chamada...`);
    
    try {
      const result = await wailsFunction();
      
      const duration = performance.now() - startTime;
      
      // Armazenar tempo
      if (!this.calls.has(functionName)) {
        this.calls.set(functionName, []);
      }
      this.calls.get(functionName)!.push(duration);
      
      const resultSize = JSON.stringify(result).length;
      
      console.log(`✅ [Wails] ${functionName} - Concluído em ${duration.toFixed(2)}ms (${resultSize} chars)`);
      
      // Alertar se muito lento
      if (duration > 1000) {
        console.warn(`⚠️ [Wails] ${functionName} - MUITO LENTO: ${duration.toFixed(2)}ms`);
      }
      
      return result;
      
    } catch (error) {
      const duration = performance.now() - startTime;
      console.error(`❌ [Wails] ${functionName} - Erro após ${duration.toFixed(2)}ms:`, error);
      throw error;
    }
  }
  
  static getStats(functionName: string) {
    const times = this.calls.get(functionName);
    if (!times || times.length === 0) {
      return null;
    }
    
    const sum = times.reduce((a, b) => a + b, 0);
    const avg = sum / times.length;
    const min = Math.min(...times);
    const max = Math.max(...times);
    
    return {
      calls: times.length,
      average: avg,
      minimum: min,
      maximum: max,
      total: sum
    };
  }
  
  static printAllStats() {
    console.log('📊 [Wails] Estatísticas de Performance:');
    console.log('==========================================');
    
    for (const [functionName, times] of this.calls.entries()) {
      const stats = this.getStats(functionName);
      if (stats) {
        console.log(`\n🔹 ${functionName}:`);
        console.log(`   • Chamadas: ${stats.calls}`);
        console.log(`   • Média: ${stats.average.toFixed(2)}ms`);
        console.log(`   • Mín/Máx: ${stats.minimum.toFixed(2)}ms / ${stats.maximum.toFixed(2)}ms`);
        console.log(`   • Total: ${stats.total.toFixed(2)}ms`);
        
        if (stats.average > 1000) {
          console.log(`   ❌ PROBLEMA: Muito lento!`);
        } else if (stats.average > 500) {
          console.log(`   ⚠️ ATENÇÃO: Um pouco lento`);
        } else {
          console.log(`   ✅ Performance boa`);
        }
      }
    }
  }
}

export default WailsProfiler;