
import csv
import matplotlib.pyplot as plt

systems = []
values = []

with open(
    "benchmark/results/research_ann.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        if row["metric"] != "recall_at_10":
            continue

        systems.append(
            row["system"],
        )

        values.append(
            float(row["value"]),
        )

plt.figure(
    figsize=(8,6),
)

plt.bar(
    systems,
    values,
)

plt.ylabel(
    "Recall@10",
)

plt.title(
    "Research HNSW Recall",
)

plt.tight_layout()

plt.savefig(
    "benchmark/charts/research_recall.png",
)
