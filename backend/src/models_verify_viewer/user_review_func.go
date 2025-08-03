package models_verify_viewer

import "log"

func (pr_old *PendingReview) MergePendingReviewItems(pr_new PendingReview) {
	log.Printf("Merging %d new items into existing %d items", len(pr_new.Items), len(pr_old.Items))

	newItemsMap := make(map[string]bool)
	for _, newitem := range pr_new.Items {
		key := newitem.JobName + "|" + newitem.DatasetName + "|" + newitem.ImageName
		newItemsMap[key] = true
	}

	filteredOldItems := NewPendingReviewItemSet()
	for _, oldItem := range pr_old.Items {
		key := oldItem.JobName + "|" + oldItem.DatasetName + "|" + oldItem.ImageName
		if newItemsMap[key] {
			filteredOldItems = append(filteredOldItems, oldItem)
		}
	}
	pr_old.Items = filteredOldItems

	for _, newitem := range pr_new.Items {
		found := false
		for _, oldItem := range pr_old.Items {
			if newitem.JobName == oldItem.JobName && newitem.DatasetName == oldItem.DatasetName && newitem.ImageName == oldItem.ImageName {
				found = true
				break
			}
		}
		if !found {
			pr_old.Items = append(pr_old.Items, newitem)
		}
	}
}

func (pr *PendingReview) ClearPendingReviewItems() {
	pr.Items = []PendingReviewItem{}
}
