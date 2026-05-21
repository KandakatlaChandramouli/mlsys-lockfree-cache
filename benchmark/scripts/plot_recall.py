
import csv
import matplotlib.pyplot as plt

systems = []
recall = []

with open(
    "benchmark/results/research_runtime.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        if row["metric"] != "recall_at_10":
            continue

        systems.append(
            row["system"],
        )

        recall.append(
            float(row["value"]),
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
    "ANN Recall Evaluation",
)

plt.tight_layout()

plt.savefig(
    "benchmark/charts/recall_eval.png",
)
