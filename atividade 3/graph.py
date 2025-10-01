import matplotlib.pyplot as plt

# Dados dos tempos médios (em ms)
gomaxprocs = [1, 2, 6]
tempos_mutex = [206, 162, 181]       # do código com Mutex
tempos_canal = [2.706, 2.0177, 2.17289]  # do código com channels

# Criar gráfico
plt.figure(figsize=(10,6))
plt.plot(gomaxprocs, tempos_mutex, marker='o', linestyle='-', color='red', label='Mutex')
plt.plot(gomaxprocs, tempos_canal, marker='s', linestyle='--', color='blue', label='Channels')

# Títulos e labels
plt.title("Comparação de Tempo Médio por GOMAXPROCS", fontsize=16)
plt.xlabel("GOMAXPROCS", fontsize=14)
plt.ylabel("Tempo médio (ms)", fontsize=14)
plt.xticks(gomaxprocs)
plt.grid(True, linestyle='--', alpha=0.5)
plt.legend(fontsize=12)

# Mostrar valores nos pontos
for i, txt in enumerate(tempos_mutex):
    plt.annotate(f"{txt:.1f}", (gomaxprocs[i], tempos_mutex[i]), textcoords="offset points", xytext=(0,10), ha='center')
for i, txt in enumerate(tempos_canal):
    plt.annotate(f"{txt:.2f}", (gomaxprocs[i], tempos_canal[i]), textcoords="offset points", xytext=(0,-15), ha='center')

plt.tight_layout()
plt.show()
