
import csv
import matplotlib.pyplot as plt

x = []
y = []

with open(
    "benchmark/advanced/results/scaling.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        x.append(
            int(row["vectors"]),
        )

        y.append(
            float(row["latency_ms"]),
        )

plt.figure(
    figsize=(8,6),
)

plt.plot(
    x,
    y,
    marker="o",
)

plt.xscale("log")

plt.xlabel(
    "Vector Count",
)

plt.ylabel(
    "Latency (ms)",
)

plt.title(
    "Scaling Curve",
)

plt.tight_layout()

plt.savefig(
    "benchmark/advanced/charts/scaling_curve.png",
)
