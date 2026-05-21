
import csv
import matplotlib.pyplot as plt

systems = []
sizes = []

with open(
    "benchmark/advanced/results/memory.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        systems.append(
            row["system"],
        )

        sizes.append(
            float(row["size_mb"]),
        )

plt.figure(
    figsize=(8,6),
)

plt.bar(
    systems,
    sizes,
)

plt.ylabel(
    "Memory (MB)",
)

plt.title(
    "Memory Compression Efficiency",
)

plt.tight_layout()

plt.savefig(
    "benchmark/advanced/charts/memory_chart.png",
)
