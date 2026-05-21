
import csv
import matplotlib.pyplot as plt

systems = []
values = []

with open(
    "benchmark/results/results.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        if "snapshot" not in row["system"]:
            continue

        systems.append(
            row["system"],
        )

        values.append(
            float(row["value"]),
        )

plt.figure(
    figsize=(8, 6),
)

plt.bar(
    systems,
    values,
)

plt.ylabel(
    "Size",
)

plt.title(
    "Snapshot Compression",
)

plt.tight_layout()

plt.savefig(
    "benchmark/charts/snapshot_compression.png",
)
