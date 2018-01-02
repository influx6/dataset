function ParseRecord(recsJSON){
	var recs = JSON.parse(recsJSON);
	return JSON.stringify([
		{
			total: recs.length,
			records: recs,
		}
	]);
};