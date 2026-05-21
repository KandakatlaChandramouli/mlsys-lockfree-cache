
import csv
import matplotlib.pyplot as plt

systems = []
values = []

with open(
    "benchmark/results/results.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        if row["metric"] != "latency_s":
            continue

        systems.append(
            row["system"],
        )

        values.append(
            float(row["value"]),
        )

plt.figure(
    figsize=(10, 6),
)

plt.bar(
    systems,
    values,
)

plt.ylabel(
    "Latency (s)",
)

plt.title(
    "Runtime Pipeline Performance",
)

plt.xticks(
    rotation=15,
)

plt.tight_layout()

plt.savefig(
    "benchmark/charts/runtime_pipeline.png",
)
