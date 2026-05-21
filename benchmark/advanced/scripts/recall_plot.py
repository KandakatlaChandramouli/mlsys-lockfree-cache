
import csv
import matplotlib.pyplot as plt

systems = []
recall = []

with open(
    "benchmark/advanced/results/recall.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        systems.append(
            row["system"],
        )

        recall.append(
            float(row["recall"]),
        )

plt.figure(
    figsize=(8,6),
)

plt.bar(
    systems,
    recall,
)

plt.ylabel(
    "Recall@10",
)

plt.title(
    "ANN Recall Comparison",
)

plt.tight_layout()

plt.savefig(
    "benchmark/advanced/charts/recall_chart.png",
)
