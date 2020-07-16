from copy import deepcopy

db=dict()

def iterSearch(orig:str, dest:str, hist:dict, bestSched:dict, bugget:float, bestOffer:float) -> (dict, float):
    import pdb; pdb.set_trace()
    # print(f"\nstarting iterSearch from {orig} to {dest}")
    # print(f"history {hist}")
    mids = db[orig]
    # print(f"mids from {orig}: {mids}")
    for mid, price in mids.items():
        if mid == dest:
            # print(f"Destination found: {orig} to {mid}")
            if bugget+price < bestOffer:
                hist[len(hist)] = dict(
                    origin=orig,
                    destination=dest,
                    price=price
                )
                bestSched = hist
                bestOffer = bugget + price
                # print(f"Found best: {bestSched}\n\tPrice: {bestOffer}")
            continue
        else:
            if db.get(mid):
                _hist = deepcopy(hist)
                _hist[len(_hist)] = dict(
                    origin=orig,
                    destination=mid,
                    price=price
                )
                bestSched, bestOffer = iterSearch(mid, dest, _hist, bestSched, bugget+price, bestOffer)
            else:
                pass
                # print(f"{mid} hasn't origin")
    # print(f"finish: {bestSched}, {bestOffer}")
    return bestSched, bestOffer


def FindBestOffer(orig:str, dest:str) -> (dict, float):
    bestOffer:float = 2**64
    schedule = dict()

    if not db.get(orig):
        return schedule, bestOffer
    # print(f"{orig} founded")

    return iterSearch(orig, dest, dict(), dict(), 0.0, bestOffer)


if __name__ == "__main__":
    db = dict(
        BRC=dict(SCL=50),
        GRD=dict(DFD=10.3),
        GRU=dict(BRC=10,
                 CDG=75,
                 ORL=56,
                 SCL=20),
        ORL=dict(CDG=5),
        SDU=dict(GRU=10.3)
    )
    print(FindBestOffer("SDU", "CDG"))
