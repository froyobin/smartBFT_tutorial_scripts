import numpy as np
import matplotlib.mlab as mlab
import matplotlib.pyplot as plt

BOUNDRY = 300


def LoadData(path):
    fo = open(path, "r")
    timedata = []
    for line in fo:
        time = int(line.split(" ")[2].strip('\n'))
        if time < BOUNDRY:
            timedata.append(time)
    return timedata


if __name__ == '__main__':
    timedata = load_data = LoadData('./goodresult_4PBFT.txt')

        # example data
    #
    # mu = 100  # mean of distribution
    # sigma = 15  # standard deviation of distribution
    # x = mu + sigma * np.random.randn(10000)
    # print (np.random.randn(10000))


    num_bins = 15
    # the histogram of the data
    n, bins, patches = plt.hist(timedata, num_bins, density=True,
                                histtype='bar',facecolor='blue', rwidth=0.5)

    # add a 'best fit' line
    # y = mlab.normpdf(bins, mu, sigma)
    # plt.plot(bins, y, 'r--')
    # plt.plot(bins, y, 'r--')
    plt.xlabel('Time (ms) ')
    plt.ylabel('Density of each interval')
    plt.title(r'Time delay distribution')

    # Tweak spacing to prevent clipping of ylabel
    plt.subplots_adjust(left=0.15)
    plt.show()
