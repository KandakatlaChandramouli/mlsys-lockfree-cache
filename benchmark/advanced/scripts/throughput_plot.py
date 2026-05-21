
import csv
import matplotlib.pyplot as plt

workers = []
qps = []

with open(
    "benchmark/advanced/results/throughput.csv",
) as f:

    reader = csv.DictReader(f)

    for row in reader:

        workers.append(
            int(row["workers"]),
        )

        qps.append(
            int(row["qps"]),
        )

plt.figure(
    figsize=(8,6),
)

plt.plot(
    workers,
    qps,
    marker="o",
)

plt.xlabel(
    "Workers",
)

plt.ylabel(
    "Queries/sec",
)

plt.title(
    "Concurrency Throughput",
)

plt.tight_layout()

plt.savefig(
    "benchmark/advanced/charts/throughput_curve.png",
)
